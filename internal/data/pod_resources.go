package data

import (
	"context"
	goJSON "encoding/json"
	"sort"

	"go.uber.org/zap"

	"gitlab.unanet.io/devops/eve/pkg/errors"
	"gitlab.unanet.io/devops/eve/pkg/json"
	"gitlab.unanet.io/devops/eve/pkg/log"
	"gitlab.unanet.io/devops/eve/pkg/mergemap"
)

type PodResourcesMap struct {
	ArtifactID                 *int      `db:"artifact_id"`
	ServiceID                  *int      `db:"service_id"`
	EnvironmentID              *int      `db:"environment_id"`
	NamespaceID                *int      `db:"namespace_id"`
	Data                       json.Text `db:"data"`
	StackingOrder              int       `db:"stacking_order"`
	PodResourcesDescription    string    `db:"pr_description"`
	PodResourcesMapDescription string    `db:"prm_description"`
}

func (r *Repo) HydrateDeployServicePodResource(ctx context.Context, svc DeployService) (json.Text, error) {
	log.Logger.Debug("hydrate pod resource", zap.Any("service", svc))
	// Get all of the matching records from the map table
	rows, err := r.db.QueryxContext(ctx, `
		SELECT
		   prm.artifact_id, 		       
		   prm.service_id, 
		   prm.environment_id, 
		   prm.namespace_id, 
		   pr.data, 
		   prm.stacking_order, 
		   prm.description as prm_description,
		   pr.description as pr_description
		FROM
			pod_resources as pr
		LEFT JOIN
			pod_resources_map prm on pr.id = prm.pod_resources_id
		WHERE
			(prm.service_id = $1 OR prm.environment_id = $2 OR prm.namespace_id = $3)
		OR 
		      ((prm.artifact_id = $4)  
		        AND (prm.environment_id IS null or prm.environment_id = $2)
		        AND (prm.namespace_id IS null or prm.namespace_id = $3))
		    		
		ORDER BY
		    prm.stacking_order,prm.service_id,prm.namespace_id,prm.environment_id
	`, svc.ServiceID, svc.EnvironmentID, svc.NamespaceID, svc.ArtifactID)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	defer rows.Close()

	// Hydrate a slice of the records to the Data Structure (PodAutoscaleMap)
	var prms []PodResourcesMap
	for rows.Next() {
		var prm PodResourcesMap
		err = rows.StructScan(&prm)
		if err != nil {
			return nil, errors.Wrap(err)
		}
		prms = append(prms, prm)
	}

	// Guard against no values set in the DB
	if len(prms) == 0 {
		log.Logger.Debug("no pod resource values set", zap.Any("service", svc))
		return json.EmptyJSONText, nil
	}

	// Explicitly Sort the slice based on stacking Order (this is done on ORDER BY, but I want to be explicit about it)
	sort.SliceStable(prms, func(i, j int) bool {
		return prms[i].StackingOrder < prms[j].StackingOrder
	})

	// Declare the "final" autoscale setting struct
	// this is the destination when we merge data on top of it
	targetPodResourceSettings := make(map[string]interface{})

	// Iterate over the sorted slice
	// Unmarshal the JSON Bytes to a temp map structure
	// Merge the temp structure onto the target/dest map struct
	// rinse and repeat, until all have been merged on top (Highest StackOrder is the last to be merged...it "wins")
	for _, prm := range prms {
		var dataMap map[string]interface{}
		merr := goJSON.Unmarshal(prm.Data, &dataMap)
		if merr != nil {
			return nil, errors.Wrap(merr)
		}
		targetPodResourceSettings = mergemap.Merge(targetPodResourceSettings, dataMap)
	}
	// Serialize the final struct back to a Byte Slice (JSON.Text) and return to the caller
	log.Logger.Debug("hydrate pod resource value", zap.Any("autoscale", targetPodResourceSettings), zap.Any("service", svc))
	return json.StructToJson(targetPodResourceSettings)
}

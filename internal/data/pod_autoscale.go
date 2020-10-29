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

type PodAutoscaleMap struct {
	ServiceID                  *int      `db:"service_id"`
	EnvironmentID              *int      `db:"environment_id"`
	NamespaceID                *int      `db:"namespace_id"`
	Data                       json.Text `db:"data"`
	StackingOrder              int       `db:"stacking_order"`
	PodAutoscaleDescription    string    `db:"pa_description"`
	PodAutoscaleMapDescription string    `db:"pam_description"`
}

func (r *Repo) HydrateDeployServicePodAutoscale(ctx context.Context, svc DeployService) (json.Text, error) {
	log.Logger.Info("hydrate pod autoscale")
	// Get all of the matching records from the map table
	rows, err := r.db.QueryxContext(ctx, `
		SELECT
		   pam.service_id, 
		   pam.environment_id, 
		   pam.namespace_id, 
		   pa.data, 
		   pam.stacking_order, 
		   pam.description as pam_description,
		   pa.description as pa_description
		FROM
			pod_autoscale as pa
		LEFT JOIN
			pod_autoscale_map pam on pa.id = pam.pod_autoscale_id
		WHERE
			pam.service_id = $1 OR pam.environment_id = $2 OR pam.namespace_id = $3
		ORDER BY
		    pam.stacking_order,pam.service_id,pam.namespace_id,pam.environment_id
	`, svc.ServiceID, svc.EnvironmentID, svc.NamespaceID)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	defer rows.Close()

	// Hydrate a slice of the records to the Data Structure (PodAutoscaleMap)
	var pads []PodAutoscaleMap
	for rows.Next() {
		var pad PodAutoscaleMap
		err = rows.StructScan(&pad)
		if err != nil {
			return nil, errors.Wrap(err)
		}
		pads = append(pads, pad)
	}

	// Guard against no values set in the DB
	if len(pads) == 0 {
		log.Logger.Debug("no pod autoscale values set", zap.Any("service", svc))
		return json.EmptyJSONText, nil
	}

	// Explicitly Sort the slice based on stacking Order
	// this is done on sql ORDER BY, but I want to be explicit about it
	sort.SliceStable(pads, func(i, j int) bool {
		return pads[i].StackingOrder < pads[j].StackingOrder
	})

	// Declare the "final" autoscale setting struct
	// this is the destination when we merge data on top of it
	targetAutoScaleSettings := make(map[string]interface{})

	// Iterate over the sorted slice
	// Unmarshal the JSON Bytes to a temp map structure
	// Merge the temp structure onto the target/dest map structure
	// rinse and repeat, until all have been merged on top (Highest StackOrder is the last to be merged...it "wins")
	for _, pad := range pads {
		var dataMap map[string]interface{}
		merr := goJSON.Unmarshal(pad.Data, &dataMap)
		if merr != nil {
			return nil, errors.Wrap(merr)
		}
		targetAutoScaleSettings = mergemap.Merge(targetAutoScaleSettings, dataMap)
	}
	// Serialize the final struct back to a Byte Slice (JSON.Text) and return to the caller
	log.Logger.Debug("hydrate pod autoscale value", zap.Any("autoscale", targetAutoScaleSettings), zap.Any("service", svc))
	return json.StructToJson(targetAutoScaleSettings)

}

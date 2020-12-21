package data

import (
	"context"
	goJSON "encoding/json"
	"sort"

	"gitlab.unanet.io/devops/go/pkg/errors"
	"gitlab.unanet.io/devops/go/pkg/json"
	"gitlab.unanet.io/devops/go/pkg/mergemap"
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

type PodResourcesByStackingOrder []PodResourcesMap

func (a PodResourcesByStackingOrder) Len() int { return len(a) }
func (a PodResourcesByStackingOrder) Less(i, j int) bool {
	return a[i].StackingOrder < a[j].StackingOrder
}
func (a PodResourcesByStackingOrder) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

func (r *Repo) NamespacePodResourcesMap(ctx context.Context, namespaceID int) ([]PodResourcesMap, error) {
	return r.PodResourcesMap(ctx, 0, 0, namespaceID, 0)
}

func (r *Repo) EnvironmentPodResourcesMap(ctx context.Context, environmentID int) ([]PodResourcesMap, error) {
	return r.PodResourcesMap(ctx, 0, environmentID, 0, 0)
}

func (r *Repo) ArtifactPodResourcesMap(ctx context.Context, artifactID, environmentID, namespaceID int) ([]PodResourcesMap, error) {
	return r.PodResourcesMap(ctx, 0, environmentID, namespaceID, artifactID)
}

func (r *Repo) PodResourcesMap(ctx context.Context, serviceID, environmentID, namespaceID, artifactID int) ([]PodResourcesMap, error) {
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
		JOIN
			pod_resources_map prm on pr.id = prm.pod_resources_id
		WHERE
		    prm.service_id = $1 AND (prm.artifact_id IS NULL AND prm.environment_id IS NULL AND prm.namespace_id IS NULL)
		OR
			prm.artifact_id = $4 AND (prm.service_id IS NULL AND prm.environment_id IS NULL AND prm.namespace_id IS NULL)
		OR
		    prm.artifact_id = $4 AND (prm.service_id IS NULL AND prm.environment_id = $2 AND prm.namespace_id IS NULL)
		OR
		    prm.artifact_id = $4 AND (prm.service_id IS NULL AND prm.environment_id IS NULL AND prm.namespace_id = $3)
		OR
		    prm.namespace_id = $3 AND (prm.artifact_id IS NULL AND prm.environment_id IS NULL AND prm.service_id IS NULL)
		OR
		    prm.environment_id = $2 AND (prm.artifact_id IS NULL AND prm.namespace_id IS NULL AND prm.service_id IS NULL)
		ORDER BY
		    prm.stacking_order
	`, serviceID, environmentID, namespaceID, artifactID)
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

	return prms, nil
}

func (r *Repo) PodResourcesStacked(prms []PodResourcesMap) (json.Text, error) {
	// Guard against no values set in the DB
	if len(prms) == 0 {
		return json.EmptyJSONText, nil
	}
	// Explicitly Sort the slice based on stacking Order (this is done on ORDER BY, but I want to be explicit about it)
	sort.Sort(PodResourcesByStackingOrder(prms))
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
	return json.StructToJson(targetPodResourceSettings)
}

func (r *Repo) HydrateDeployServicePodResource(ctx context.Context, svc DeployService) (json.Text, error) {
	// Get all of the matching records from the map table
	prms, err := r.PodResourcesMap(ctx, svc.ServiceID, svc.EnvironmentID, svc.NamespaceID, svc.ArtifactID)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	return r.PodResourcesStacked(prms)
}

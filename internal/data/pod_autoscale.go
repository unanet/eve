package data

import (
	"context"
	goJSON "encoding/json"
	"sort"

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

type PodAutoscaleByStackingOrder []PodAutoscaleMap

func (a PodAutoscaleByStackingOrder) Len() int { return len(a) }
func (a PodAutoscaleByStackingOrder) Less(i, j int) bool {
	return a[i].StackingOrder < a[j].StackingOrder
}
func (a PodAutoscaleByStackingOrder) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

func (r *Repo) PodAutoscaleMap(ctx context.Context, serviceID, environmentID, namespaceID int) ([]PodAutoscaleMap, error) {
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
		JOIN
			pod_autoscale_map pam on pa.id = pam.pod_autoscale_id
		WHERE
		    pam.service_id = $1 AND (pam.namespace_id IS NULL AND pam.environment_id IS NULL)
		OR  
			pam.environment_id = $2 AND (pam.service_id IS NULL AND pam.namespace_id IS NULL)
		OR  
			pam.namespace_id = $3 AND (pam.service_id IS NULL AND pam.environment_id IS NULL)		
		ORDER BY
		    pam.stacking_order
	`, serviceID, environmentID, namespaceID)
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
	return pads, nil
}

func (r *Repo) NamespacePodAutoscaleMap(ctx context.Context, namespaceID int) ([]PodAutoscaleMap, error) {
	log.Logger.Info("namespace pod autoscale map")
	return r.PodAutoscaleMap(ctx, 0, 0, namespaceID)
}

func (r *Repo) EnvironmentPodAutoscaleMap(ctx context.Context, environmentID int) ([]PodAutoscaleMap, error) {
	log.Logger.Info("environment pod autoscale map")
	return r.PodAutoscaleMap(ctx, 0, environmentID, 0)
}

func (r *Repo) PodAutoscaleStacked(pams []PodAutoscaleMap) (json.Text, error) {
	// Guard against no values set in the DB
	if len(pams) == 0 {
		return json.EmptyJSONText, nil
	}

	// Explicitly Sort the slice based on stacking Order
	// this is done on sql ORDER BY, but I want to be explicit about it
	sort.Sort(PodAutoscaleByStackingOrder(pams))
	//sort.SliceStable(pads, func(i, j int) bool {
	//	return pads[i].StackingOrder < pads[j].StackingOrder
	//})

	// Declare the "final" autoscale setting struct
	// this is the destination when we merge data on top of it
	targetAutoScaleSettings := make(map[string]interface{})

	// Iterate over the sorted slice
	// Unmarshal the JSON Bytes to a temp map structure
	// Merge the temp structure onto the target/dest map structure
	// rinse and repeat, until all have been merged on top (Highest StackOrder is the last to be merged...it "wins")
	for _, pad := range pams {
		var dataMap map[string]interface{}
		merr := goJSON.Unmarshal(pad.Data, &dataMap)
		if merr != nil {
			return nil, errors.Wrap(merr)
		}
		targetAutoScaleSettings = mergemap.Merge(targetAutoScaleSettings, dataMap)
	}
	// Serialize the final struct back to a Byte Slice (JSON.Text) and return to the caller
	return json.StructToJson(targetAutoScaleSettings)
}

func (r *Repo) HydrateDeployServicePodAutoscale(ctx context.Context, svc DeployService) (json.Text, error) {
	log.Logger.Info("hydrate deploy service pod autoscale")
	// Get all of the matching records from the map table
	pams, err := r.PodAutoscaleMap(ctx, svc.ServiceID, svc.EnvironmentID, svc.NamespaceID)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	return r.PodAutoscaleStacked(pams)
}

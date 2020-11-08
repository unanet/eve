package eve

import (
	"gitlab.unanet.io/devops/eve/pkg/json"
)

type PodResourcesMapSource struct {
	ArtifactID                 *int      `json:"artifact_id"`
	ServiceID                  *int      `json:"service_id"`
	EnvironmentID              *int      `json:"environment_id"`
	NamespaceID                *int      `json:"namespace_id"`
	Data                       json.Text `json:"data"`
	StackingOrder              int       `json:"stacking_order"`
	PodResourcesDescription    string    `json:"pr_description"`
	PodResourcesMapDescription string    `json:"prm_description"`
}

type PodResources struct {
	Sources []PodResourcesMapSource `json:"sources"`
	Data    json.Text               `json:"data"`
}
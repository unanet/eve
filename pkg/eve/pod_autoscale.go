package eve

import (
	"gitlab.unanet.io/devops/eve/pkg/json"
)

type PodAutoscaleMapSource struct {
	ServiceID                  *int      `json:"service_id,omitempty"`
	EnvironmentID              *int      `json:"environment_id,omitempty"`
	NamespaceID                *int      `json:"namespace_id,omitempty"`
	Data                       json.Text `json:"data"`
	StackingOrder              int       `json:"stacking_order"`
	PodAutoscaleDescription    string    `json:"pa_description"`
	PodAutoscaleMapDescription string    `json:"pam_description"`
}

type PodAutoscale struct {
	Sources []PodAutoscaleMapSource `json:"sources"`
	Data    json.Text               `json:"data"`
}

package eve

import "time"

type Namespace struct {
	ID               int                    `json:"id"`
	Name             string                 `json:"name"`
	Alias            string                 `json:"alias"`
	EnvironmentID    int                    `json:"environment_id"`
	EnvironmentName  string                 `json:"environment_name"`
	RequestedVersion string                 `json:"requested_version"`
	ExplicitDeploy   bool                   `json:"explicit_deploy"`
	ClusterID        int                    `json:"cluster_id"`
	Metadata         map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt        time.Time              `json:"created_at"`
	UpdatedAt        time.Time              `json:"updated_at"`
}
package eve

import (
	"time"
)

type Namespace struct {
	ID                 int                    `db:"id"`
	Name               string                 `db:"name"`
	Alias              string                 `db:"alias"`
	EnvironmentID      int                    `db:"environment_id"`
	EnvironmentName    string                 `db:"environment_name"`
	RequestedVersion   string                 `db:"requested_version"`
	ExplicitDeployOnly bool                   `db:"explicit_deploy_only"`
	ClusterID          int                    `db:"cluster_id"`
	Metadata           map[string]interface{} `db:"metadata,omitempty"`
	CreatedAt          time.Time              `db:"created_at"`
	UpdatedAt          time.Time              `db:"updated_at"`
}

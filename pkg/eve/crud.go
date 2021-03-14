package eve

import (
	"time"
)

type Environment struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Alias       string    `json:"alias,omitempty"`
	Description string    `json:"description"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Namespace struct {
	ID                 int                    `json:"id"`
	Name               string                 `json:"name"`
	Alias              string                 `json:"alias"`
	EnvironmentID      int                    `json:"environment_id"`
	EnvironmentName    string                 `json:"environment_name"`
	RequestedVersion   string                 `json:"requested_version"`
	ExplicitDeployOnly bool                   `json:"explicit_deploy_only"`
	ClusterID          int                    `json:"cluster_id"`
	Metadata           map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt          time.Time              `json:"created_at"`
	UpdatedAt          time.Time              `json:"updated_at"`
}

type Service struct {
	ID              int       `json:"id"`
	NamespaceID     int       `json:"namespace_id"`
	NamespaceName   string    `json:"namespace_name"`
	ArtifactID      int       `json:"artifact_id"`
	ArtifactName    string    `json:"artifact_name"`
	OverrideVersion string    `json:"override_version"`
	DeployedVersion string    `json:"deployed_version"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	Name            string    `json:"name"`
	StickySessions  bool      `json:"sticky_sessions"`
	Count           int       `json:"count"`
}

type Job struct {
	ID              int       `json:"id"`
	NamespaceID     int       `json:"namespace_id"`
	NamespaceName   string    `json:"namespace_name"`
	ArtifactID      int       `json:"artifact_id"`
	ArtifactName    string    `json:"artifact_name"`
	OverrideVersion string    `json:"override_version"`
	DeployedVersion string    `json:"deployed_version"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	Name            string    `json:"name"`
}

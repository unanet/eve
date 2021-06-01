package eve

import "time"

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
	Count           int       `json:"count"`
	ExplicitDeploy  bool      `json:"explicit_deploy"`
}

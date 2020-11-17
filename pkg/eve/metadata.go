package eve

import (
	"time"
)

type Metadata struct {
	ID          int                    `json:"id"`
	Description string                 `json:"description"`
	Value       map[string]interface{} `json:"value"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

type MetadataServiceMap struct {
	Description   string    `json:"description"`
	MetadataID    int       `json:"metadata_id"`
	EnvironmentID int       `json:"environment_id"`
	ArtifactID    int       `json:"artifact_id"`
	NamespaceID   int       `json:"namespace_id"`
	ServiceID     int       `json:"service_id"`
	StackingOrder int       `json:"stacking_order"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type MetadataJobMap struct {
	Description   string    `json:"description"`
	MetadataID    int       `json:"metadata_id"`
	EnvironmentID int       `json:"environment_id"`
	ArtifactID    int       `json:"artifact_id"`
	NamespaceID   int       `json:"namespace_id"`
	JobID         int       `json:"service_id"`
	StackingOrder int       `json:"stacking_order"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

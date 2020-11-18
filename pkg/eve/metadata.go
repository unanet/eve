package eve

import (
	"context"
	"errors"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type MetadataField map[string]interface{}

func (m MetadataField) ValidateWithContext(ctx context.Context) error {
	if m == nil {
		return nil
	}

	if _, ok := m[""]; ok {
		return errors.New("cannot have an empty value as a key for metadata")
	}

	return nil
}

type Metadata struct {
	ID          int           `json:"id"`
	Description string        `json:"description"`
	Value       MetadataField `json:"value"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
}

func (m Metadata) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, &m,
		validation.Field(&m.Description, validation.Required),
		validation.Field(&m.Value))
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

func (m MetadataServiceMap) environmentIDSet() int {
	if m.EnvironmentID > 0 {
		return 1
	} else {
		return 0
	}
}

func (m MetadataServiceMap) artifactIDSet() int {
	if m.ArtifactID > 0 {
		return 1
	} else {
		return 0
	}
}

func (m MetadataServiceMap) namespaceIDSet() int {
	if m.NamespaceID > 0 {
		return 1
	} else {
		return 0
	}
}

func (m MetadataServiceMap) serviceIDSet() int {
	if m.ServiceID > 0 {
		return 1
	} else {
		return 0
	}
}

func (m MetadataServiceMap) ValidateWithContext(ctx context.Context) error {
	if err := validation.ValidateStructWithContext(ctx, &m,
		validation.Field(&m.Description, validation.Required)); err != nil {
		return err
	}

	if m.EnvironmentID+m.ArtifactID+m.NamespaceID+m.ServiceID == 0 {
		return errors.New("you must set either service_id, environment_id, namespace_id, or artifact_id")
	}

	if m.serviceIDSet()+m.environmentIDSet()+m.namespaceIDSet() > 1 {
		return errors.New("you may only set one of the 3 fields: service_id, namespace_id, or environment_id")
	}

	if m.artifactIDSet()+m.serviceIDSet() > 1 {
		return errors.New("you may only set the artifact_id or service_id field")
	}

	return nil
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

func (m MetadataJobMap) environmentIDSet() int {
	if m.EnvironmentID > 0 {
		return 1
	} else {
		return 0
	}
}

func (m MetadataJobMap) artifactIDSet() int {
	if m.ArtifactID > 0 {
		return 1
	} else {
		return 0
	}
}

func (m MetadataJobMap) namespaceIDSet() int {
	if m.NamespaceID > 0 {
		return 1
	} else {
		return 0
	}
}

func (m MetadataJobMap) jobIDSet() int {
	if m.JobID > 0 {
		return 1
	} else {
		return 0
	}
}

func (m MetadataJobMap) ValidateWithContext(ctx context.Context) error {
	if err := validation.ValidateStructWithContext(ctx, &m,
		validation.Field(&m.Description, validation.Required)); err != nil {
		return err
	}

	if m.EnvironmentID+m.ArtifactID+m.NamespaceID+m.JobID == 0 {
		return errors.New("you must set either service_id, environment_id, namespace_id, or artifact_id")
	}

	if m.jobIDSet()+m.environmentIDSet()+m.namespaceIDSet() > 1 {
		return errors.New("you may only set one of the 3 fields: service_id, namespace_id, or environment_id")
	}

	if m.artifactIDSet()+m.jobIDSet() > 1 {
		return errors.New("you may only set the artifact_id or service_id field")
	}

	return nil
}

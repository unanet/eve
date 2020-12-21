package eve

import (
	"context"
	"errors"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Annotation struct {
	ID          int           `json:"id"`
	Description string        `json:"description"`
	Data        MetadataField `json:"data"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
}

func (m Annotation) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, &m,
		validation.Field(&m.Description, validation.Required),
		validation.Field(&m.Data))
}

type AnnotationServiceMap struct {
	Description   string    `json:"description"`
	AnnotationID  int       `json:"annotation_id"`
	EnvironmentID int       `json:"environment_id"`
	ArtifactID    int       `json:"artifact_id"`
	NamespaceID   int       `json:"namespace_id"`
	ServiceID     int       `json:"service_id"`
	StackingOrder int       `json:"stacking_order"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (m AnnotationServiceMap) environmentIDSet() int {
	if m.EnvironmentID > 0 {
		return 1
	} else {
		return 0
	}
}

func (m AnnotationServiceMap) artifactIDSet() int {
	if m.ArtifactID > 0 {
		return 1
	} else {
		return 0
	}
}

func (m AnnotationServiceMap) namespaceIDSet() int {
	if m.NamespaceID > 0 {
		return 1
	} else {
		return 0
	}
}

func (m AnnotationServiceMap) serviceIDSet() int {
	if m.ServiceID > 0 {
		return 1
	} else {
		return 0
	}
}

func (m AnnotationServiceMap) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, &m,
		validation.Field(&m.Description, validation.Required),
		validation.Field(&m.ServiceID, validation.By(func(value interface{}) error {
			if m.EnvironmentID+m.ArtifactID+m.NamespaceID+m.ServiceID == 0 {
				return errors.New("you must set either service_id, environment_id, namespace_id, or artifact_id")
			}
			return nil
		})),
		validation.Field(&m.ServiceID, validation.By(func(value interface{}) error {
			if m.serviceIDSet()+m.environmentIDSet()+m.namespaceIDSet() > 1 {
				return errors.New("you may only set one of the 3 fields: service_id, namespace_id, or environment_id")
			}
			return nil
		})),
		validation.Field(&m.ArtifactID, validation.By(func(value interface{}) error {
			if m.artifactIDSet()+m.serviceIDSet() > 1 {
				return errors.New("you may only set the artifact_id or service_id field")
			}
			return nil
		})))
}

type AnnotationJobMap struct {
	Description   string    `json:"description"`
	AnnotationID  int       `json:"annotation_id"`
	EnvironmentID int       `json:"environment_id"`
	ArtifactID    int       `json:"artifact_id"`
	NamespaceID   int       `json:"namespace_id"`
	JobID         int       `json:"service_id"`
	StackingOrder int       `json:"stacking_order"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (m AnnotationJobMap) environmentIDSet() int {
	if m.EnvironmentID > 0 {
		return 1
	} else {
		return 0
	}
}

func (m AnnotationJobMap) artifactIDSet() int {
	if m.ArtifactID > 0 {
		return 1
	} else {
		return 0
	}
}

func (m AnnotationJobMap) namespaceIDSet() int {
	if m.NamespaceID > 0 {
		return 1
	} else {
		return 0
	}
}

func (m AnnotationJobMap) jobIDSet() int {
	if m.JobID > 0 {
		return 1
	} else {
		return 0
	}
}

func (m AnnotationJobMap) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, &m,
		validation.Field(&m.Description, validation.Required),
		validation.Field(&m.JobID, validation.By(func(value interface{}) error {
			if m.EnvironmentID+m.ArtifactID+m.NamespaceID+m.JobID == 0 {
				return errors.New("you must set either job_id, environment_id, namespace_id, or artifact_id")
			}
			return nil
		})),
		validation.Field(&m.JobID, validation.By(func(value interface{}) error {
			if m.jobIDSet()+m.environmentIDSet()+m.namespaceIDSet() > 1 {
				return errors.New("you may only set one of the 3 fields: job_id, namespace_id, or environment_id")
			}
			return nil
		})),
		validation.Field(&m.ArtifactID, validation.By(func(value interface{}) error {
			if m.artifactIDSet()+m.jobIDSet() > 1 {
				return errors.New("you may only set the job_id or service_id field")
			}
			return nil
		})))
}

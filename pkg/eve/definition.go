package eve

import (
	"context"
	"errors"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"time"
)

type DefinitionType string

const (
	TDeployment DefinitionType = "appsv1.Deployment"
	TJob        DefinitionType = "batchv1.Job"
	TAutoScale  DefinitionType = "v2beta2.HorizontalPodAutoscaler"
)

type DefinitionSpec map[string]map[string]interface{}


type Definition struct {
	ID               int           `json:"id"`
	Description      string        `json:"description"`
	DefinitionTypeID int           `json:"definition_type_id"`
	Data             MetadataField `json:"data"`
	CreatedAt        time.Time     `json:"created_at"`
	UpdatedAt        time.Time     `json:"updated_at"`
}

func (d Definition) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, &d,
		validation.Field(&d.Description, validation.Required),
		validation.Field(&d.DefinitionTypeID, validation.Required),
		validation.Field(&d.Data))
}

type DefinitionServiceMap struct {
	Description   string    `json:"description"`
	DefinitionID  int       `json:"definition_id"`
	EnvironmentID int       `json:"environment_id"`
	ArtifactID    int       `json:"artifact_id"`
	NamespaceID   int       `json:"namespace_id"`
	ServiceID     int       `json:"service_id"`
	ClusterID     int       `json:"cluster_id"`
	StackingOrder int       `json:"stacking_order"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (m DefinitionServiceMap) environmentIDSet() int {
	if m.EnvironmentID > 0 {
		return 1
	} else {
		return 0
	}
}

func (m DefinitionServiceMap) artifactIDSet() int {
	if m.ArtifactID > 0 {
		return 1
	} else {
		return 0
	}
}

func (m DefinitionServiceMap) namespaceIDSet() int {
	if m.NamespaceID > 0 {
		return 1
	} else {
		return 0
	}
}

func (m DefinitionServiceMap) serviceIDSet() int {
	if m.ServiceID > 0 {
		return 1
	} else {
		return 0
	}
}

func (m DefinitionServiceMap) clusterIDSet() int {
	if m.ClusterID > 0 {
		return 1
	} else {
		return 0
	}
}

func (m DefinitionServiceMap) ValidateWithContext(ctx context.Context) error {
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

type DefinitionJobMap struct {
	Description   string    `json:"description"`
	DefinitionID  int       `json:"definition_id"`
	EnvironmentID int       `json:"environment_id"`
	ArtifactID    int       `json:"artifact_id"`
	NamespaceID   int       `json:"namespace_id"`
	ClusterID     int       `json:"cluster_id"`
	JobID         int       `json:"service_id"`
	StackingOrder int       `json:"stacking_order"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (m DefinitionJobMap) environmentIDSet() int {
	if m.EnvironmentID > 0 {
		return 1
	} else {
		return 0
	}
}

func (m DefinitionJobMap) artifactIDSet() int {
	if m.ArtifactID > 0 {
		return 1
	} else {
		return 0
	}
}

func (m DefinitionJobMap) namespaceIDSet() int {
	if m.NamespaceID > 0 {
		return 1
	} else {
		return 0
	}
}

func (m DefinitionJobMap) jobIDSet() int {
	if m.JobID > 0 {
		return 1
	} else {
		return 0
	}
}

func (m DefinitionJobMap) clusterIDSet() int {
	if m.ClusterID > 0 {
		return 1
	} else {
		return 0
	}
}

func (m DefinitionJobMap) ValidateWithContext(ctx context.Context) error {
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

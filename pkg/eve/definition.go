package eve

import (
	"context"
	"errors"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"strconv"
	"strings"
	"time"
)

// TODO: Remove once Defaults are applied to every service/job
func DefaultServiceResourceDef() DefinitionResult {
	return DefinitionResult{
		Class:   "",
		Version: "v1",
		Kind:    "Service",
		Order:   "main",
		Data: map[string]interface{}{
			"spec": map[string]interface{}{},
		},
	}
}

// TODO: Remove once Defaults are applied to every service/job
func DefaultDeploymentResourceDef() DefinitionResult {
	return DefinitionResult{
		Class:   "apps",
		Version: "v1",
		Kind:    "Deployment",
		Order:   "main",
		Data: map[string]interface{}{
			"spec": map[string]interface{}{},
		},
	}
}

// TODO: Remove once Defaults are applied to every service/job
func DefaultJobResourceDef() DefinitionResult {
	return DefinitionResult{
		Class:   "batch",
		Version: "v1",
		Kind:    "Job",
		Order:   "main",
		Data: map[string]interface{}{
			"spec": map[string]interface{}{},
		},
	}
}

type DefinitionResults []DefinitionResult

func (drs DefinitionResults) CRDs(order string) []DefinitionResult {
	var result = make([]DefinitionResult, 0)
	for _, dr := range drs {
		if dr.Order == order {
			result = append(result, dr)
		}
	}
	return result
}

type ResourceDefinition struct {
	APIVersion string      `json:"apiVersion"`
	Kind       string      `json:"kind"`
	Meta       interface{} `json:"metadata"`
	Spec       interface{} `json:"spec"`
}

type DefinitionResult struct {
	Class   string                 `json:"class"`
	Version string                 `json:"version"`
	Kind    string                 `json:"kind"`
	Order   string                 `json:"order"`
	Data    map[string]interface{} `json:"data"`
}

func (dr DefinitionResult) ToResourceDefinition() ResourceDefinition {
	return ResourceDefinition{
		APIVersion: dr.APIVersion(),
		Kind:       dr.Kind,
		Meta:       dr.Data["metadata"],
		Spec:       dr.Data["spec"],
	}
}

// TODO: Remove once Standard Values (Annotations and Labels) are migrated to all Def
func (dr *DefinitionResult) AnnotationKeys() []string {
	switch strings.ToLower(dr.Kind) {
	case "service":
		return []string{"metadata", "annotations"}
	}

	return []string{"spec", "template", "metadata", "annotations"}
}

// TODO: Migrate to the Definition with Templated values instead
func (dr *DefinitionResult) StandardAnnotations(eveDeployment DeploymentSpec) map[string]interface{} {
	switch strings.ToLower(dr.Kind) {
	case "deployment":
		// If the service has a metrics port we will set up scrape label here
		// TODO: remove after migration from eve service to definition
		if eveDeployment.GetMetricsPort() != 0 {
			return map[string]interface{}{
				"prometheus.io/scrape": "true",
				"prometheus.io/port":   strconv.Itoa(eveDeployment.GetMetricsPort()),
			}
		}
	}
	return map[string]interface{}{}
}

// TODO: Migrate to the Definition with Templated values instead
func (dr *DefinitionResult) StandardLabels(eveDeployment DeploymentSpec) map[string]interface{} {
	switch strings.ToLower(dr.Kind) {
	case "deployment":
		return map[string]interface{}{
			"app":     eveDeployment.GetName(),
			"version": eveDeployment.GetArtifact().AvailableVersion,
			"nuance":  eveDeployment.GetNuance(),
			"metrics": eveDeployment.Metrics(),
		}
	case "job":
		return map[string]interface{}{
			"job":     eveDeployment.GetName(),
			"version": eveDeployment.GetArtifact().AvailableVersion,
		}
	}

	return map[string]interface{}{}
}

func (dr *DefinitionResult) LabelKeys() []string {
	// Overrides
	switch strings.ToLower(dr.Kind) {
	case "service":
		return []string{"metadata", "labels"}
	}

	return []string{"spec", "template", "metadata", "labels"}
}

func (dr *DefinitionResult) APIVersion() string {
	// Some definitions dont have a Class (aka Group in K8s)
	// So we trim the leading "/" when there isn't a Class
	return strings.TrimPrefix(fmt.Sprintf("%s/%s", dr.Class, dr.Version),"/")
}

// Key is used to merge the data from slice to map in the service
// ex: main.apps.v1.Deployment
// we use it as a key in the map, and then split on the "." to extract the individual values
func (dr *DefinitionResult) Key() string {
	return fmt.Sprintf("%s.%s.%s.%s", dr.Order, dr.Class, dr.Version, dr.Kind)
}

// Resource is the just the lowercase, plural form of the Kind property
func (dr *DefinitionResult) Resource() string {
	if strings.HasSuffix(dr.Kind, "s") {
		return strings.ToLower(dr.Kind)
	}
	return strings.ToLower(fmt.Sprintf("%ss", dr.Kind))
}

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

package eve

import (
	"context"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type ReleaseType string

const (
	ReleaseTypeArtifact  ReleaseType = "artifact"
	ReleaseTypeNamespace ReleaseType = "namespace"
)
type Release struct {
	Type     ReleaseType `json:"type"`
	FromFeed string      `json:"from_feed"`
	ToFeed   string      `json:"to_feed"`
	Message  string      `json:"message,omitempty"`

	// Artifact Release
	Artifact string `json:"artifact,omitempty"`
	Version  string `json:"version,omitempty"`

	// Namespace Release
	Namespace   string `json:"namespace,omitempty"`
	Environment string `json:"environment,omitempty"`
}

func (r Release) ValidateWithContext(ctx context.Context) error {
	releaseNamespaceValidation := validation.NilOrNotEmpty.When(r.Type == ReleaseTypeNamespace)
	releaseArtifactValidation := validation.NilOrNotEmpty.When(r.Type == ReleaseTypeArtifact)

	return validation.ValidateStructWithContext(ctx, &r,
		validation.Field(&r.Type, validation.Required),

		validation.Field(&r.Namespace, releaseNamespaceValidation.Error("namespace is required when release type is artifact")),
		validation.Field(&r.Artifact, releaseArtifactValidation.Error("artifact is required when release type is artifact")),

		validation.Field(&r.FromFeed, validation.Required),
		validation.Field(&r.Environment, releaseNamespaceValidation.Error("environment is required when releasing a namespace")),
	)
}

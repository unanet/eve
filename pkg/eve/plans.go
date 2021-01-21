package eve

import (
	"context"
	"fmt"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	uuid "github.com/satori/go.uuid"
)

type StringList []string

func (s StringList) Contains(value string) bool {
	for _, a := range s {
		if a == value {
			return true
		}
	}
	return false
}

type ArtifactDefinition struct {
	ID               int    `json:"id"`
	Name             string `json:"name"`
	ArtifactName     string `json:"artifact_name"`
	RequestedVersion string `json:"requested_version,omitempty"`
	AvailableVersion string `json:"available_version"`
	ArtifactoryFeed  string `json:"artifactory_feed"`
	ArtifactoryPath  string `json:"artifactory_path"`
	FunctionPointer  string `json:"function_pointer"`
	FeedType         string `json:"feed_type"`
	Matched          bool   `json:"-"`
}

func (ad ArtifactDefinition) ArtifactoryRequestedVersion() string {
	if ad.RequestedVersion == "latest" {
		return "*"
	} else if ad.RequestedVersion == "" {
		return "*"
	} else if len(strings.Split(ad.RequestedVersion, ".")) < 4 {
		return ad.RequestedVersion + ".*"
	}
	return ad.RequestedVersion
}

type ArtifactDefinitions []*ArtifactDefinition

func (ad ArtifactDefinitions) ContainsVersion(name string, version string) bool {
	for _, x := range ad {
		if x.AvailableVersion == version && x.ArtifactName == name {
			return true
		}
	}
	return false
}

func (ad ArtifactDefinitions) Match(artifactID int, optName string, requestedVersion string) *ArtifactDefinition {
	for _, x := range ad {
		if x.Name != "" {
			if x.Name == optName && strings.HasPrefix(x.AvailableVersion, requestedVersion) {
				return x
			}
		} else if x.ID == artifactID && strings.HasPrefix(x.AvailableVersion, requestedVersion) {
			return x
		}
	}
	return nil
}

func (ad ArtifactDefinitions) UnMatched() ArtifactDefinitions {
	var unmatched ArtifactDefinitions
	for _, x := range ad {
		if !x.Matched {
			unmatched = append(unmatched, x)
		}
	}
	return unmatched
}

type DeploymentPlanOptions struct {
	Artifacts        ArtifactDefinitions `json:"artifacts"`
	ForceDeploy      bool                `json:"force_deploy"`
	User             string              `json:"user"`
	DryRun           bool                `json:"dry_run"`
	CallbackURL      string              `json:"callback_url"`
	Environment      string              `json:"environment"`
	NamespaceAliases StringList          `json:"namespaces,omitempty"`
	Messages         []string            `json:"messages,omitempty"`
	Type             PlanType            `json:"type"`
	DeploymentIDs    []uuid.UUID         `json:"deployment_ids,omitempty"`
	Metadata         MetadataField       `json:"metadata"`
}

func (po *DeploymentPlanOptions) PlanType() string {
	return string(po.Type)
}

func (po *DeploymentPlanOptions) Message(format string, a ...interface{}) {
	po.Messages = append(po.Messages, fmt.Sprintf(format, a...))
}

func (po DeploymentPlanOptions) HasArtifacts() bool {
	return len(po.Artifacts) > 0
}

func (po DeploymentPlanOptions) HasNamespaceAliases() bool {
	return len(po.NamespaceAliases) > 0
}

func (po DeploymentPlanOptions) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, &po,
		validation.Field(&po.Environment, validation.Required),
		validation.Field(
			&po.Type,
			validation.Required,
			validation.In(
				DeploymentPlanTypeApplication,
				DeploymentPlanTypeMigration,
				DeploymentPlanTypeJob,
				DeploymentPlanTypeRestart,
			),
		),
		validation.Field(&po.User, validation.Required))
}

type NamespacePlanOptions struct {
	NamespaceRequest  *NamespaceRequest   `json:"namespace"`
	Artifacts         ArtifactDefinitions `json:"artifacts"`
	ArtifactsSupplied bool                `json:"artifacts_supplied"`
	ForceDeploy       bool                `json:"force_deploy"`
	DryRun            bool                `json:"dry_run"`
	CallbackURL       string              `json:"callback_url"`
	EnvironmentID     int                 `json:"environment_id"`
	EnvironmentName   string              `json:"environment_name"`
	EnvironmentAlias  string              `json:"environment_alias"`
	Type              PlanType            `json:"type"`
	Metadata          MetadataField       `json:"metadata"`
}

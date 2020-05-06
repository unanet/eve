package eve

import (
	"fmt"
	"strings"

	uuid "github.com/satori/go.uuid"
)

type DeployArtifactResult string

const (
	DeployArtifactResultNoop      DeployArtifactResult = "noop"
	DeployArtifactResultSucceeded DeployArtifactResult = "succeeded"
	DeployArtifactResultFailed    DeployArtifactResult = "failed"
)

func ParseDeployArtifactResult(value string) DeployArtifactResult {
	switch strings.ToLower(value) {
	case "noop":
		return DeployArtifactResultNoop
	case "succeeded":
		return DeployArtifactResultSucceeded
	case "failed":
		return DeployArtifactResultFailed
	default:
		return DeployArtifactResultNoop
	}
}

type DeploymentPlanStatus string

const (
	DeploymentPlanStatusPending  DeploymentPlanStatus = "pending"
	DeploymentPlanStatusDryrun   DeploymentPlanStatus = "dryrun"
	DeploymentPlanStatusErrors   DeploymentPlanStatus = "errors"
	DeploymentPlanStatusComplete DeploymentPlanStatus = "complete"
)

type DeployArtifact struct {
	ArtifactID       int                    `json:"artifact_id"`
	ArtifactName     string                 `json:"artifact_name"`
	RequestedVersion string                 `json:"requested_version"`
	DeployedVersion  string                 `json:"deployed_version"`
	AvailableVersion string                 `json:"available_version"`
	Metadata         map[string]interface{} `json:"metadata"`
	ArtifactoryFeed  string                 `json:"artifactory_feed"`
	ArtifactoryPath  string                 `json:"artifactory_path"`
	ArtifactFnPtr    string                 `json:"artifact_fn"`
	Result           DeployArtifactResult   `json:"result"`
	Deploy           bool                   `json:"-"`
}

type DeployService struct {
	*DeployArtifact
	ServiceID int `json:"service_id"`
}

type DeployServices []*DeployService

func (ds DeployServices) ToDeploy() DeployServices {
	var list DeployServices
	for _, x := range ds {
		if x.Deploy {
			list = append(list, x)
		}
	}
	return list
}

type DeployMigration struct {
	*DeployArtifact
	DatabaseID   int    `json:"database_id"`
	DatabaseName string `json:"database_name"`
}

type DeployMigrations []*DeployMigration

func (ds DeployMigrations) ToDeploy() DeployMigrations {
	var list DeployMigrations
	for _, x := range ds {
		if x.Deploy {
			list = append(list, x)
		}
	}
	return list
}

type NamespaceRequest struct {
	ID          int    `json:"id"`
	Alias       string `json:"alias"`
	Name        string `json:"name"`
	ClusterID   int    `json:"cluster_id"`
	ClusterName string `json:"cluster_name"`
}

func (ns *NamespaceRequest) GetQueueGroupID() string {
	return fmt.Sprintf("deploy-%s", ns.Name)
}

type NamespaceRequests []*NamespaceRequest

func (n NamespaceRequests) ToIDs() []int {
	var ids []int
	for _, x := range n {
		ids = append(ids, x.ID)
	}
	return ids
}

type NSDeploymentPlan struct {
	DeploymentID    uuid.UUID            `json:"deploymend_id"`
	Namespace       *NamespaceRequest    `json:"namespace"`
	EnvironmentName string               `json:"environment_name"`
	Services        DeployServices       `json:"services,omitempty"`
	Migrations      DeployMigrations     `json:"migrations,omitempty"`
	Messages        []string             `json:"messages,omitempty"`
	SchQueueUrl     string               `json:"-"`
	CallbackURL     string               `json:"callback_url"`
	Status          DeploymentPlanStatus `json:"status"`
}

func (ns *NSDeploymentPlan) GroupID() string {
	return ns.Namespace.Name
}

func (ns *NSDeploymentPlan) NoopExist() bool {
	for _, x := range ns.Services {
		if x.Result == DeployArtifactResultNoop {
			return true
		}
	}

	for _, x := range ns.Migrations {
		if x.Result == DeployArtifactResultNoop {
			return true
		}
	}

	return true
}

func (ns *NSDeploymentPlan) Failed() bool {
	for _, x := range ns.Services {
		if x.Result == DeployArtifactResultFailed {
			return true
		}
	}

	for _, x := range ns.Migrations {
		if x.Result == DeployArtifactResultFailed {
			return true
		}
	}

	return true
}

func (ns *NSDeploymentPlan) Message(format string, a ...interface{}) {
	ns.Messages = append(ns.Messages, fmt.Sprintf(format, a...))
}

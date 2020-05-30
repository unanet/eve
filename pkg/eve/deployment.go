package eve

import (
	"fmt"
	"strings"

	uuid "github.com/satori/go.uuid"
)

type DeployArtifactResult string

const (
	ArtifactoryFeedTypeDocker = "docker"
)

const (
	DeployArtifactResultNoop    DeployArtifactResult = "noop"
	DeployArtifactResultSuccess DeployArtifactResult = "success"
	DeployArtifactResultFailed  DeployArtifactResult = "failed"
)

func (dar DeployArtifactResult) String() string {
	return string(dar)
}

func ParseDeployArtifactResult(value string) DeployArtifactResult {
	switch strings.ToLower(value) {
	case "noop":
		return DeployArtifactResultNoop
	case "success":
		return DeployArtifactResultSuccess
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

func (dps DeploymentPlanStatus) String() string {
	return string(dps)
}

type DeployArtifact struct {
	ArtifactID          int                    `json:"artifact_id"`
	ArtifactName        string                 `json:"artifact_name"`
	RequestedVersion    string                 `json:"requested_version"`
	DeployedVersion     string                 `json:"deployed_version"`
	AvailableVersion    string                 `json:"available_version"`
	ServiceAccount      string                 `json:"service_account"`
	ImageTag            string                 `json:"image_tag"`
	Metadata            map[string]interface{} `json:"metadata"`
	ArtifactoryFeed     string                 `json:"artifactory_feed"`
	ArtifactoryPath     string                 `json:"artifactory_path"`
	ArtifactFnPtr       string                 `json:"artifact_fn"`
	ArtifactoryFeedType string                 `json:"artifactory_feed_type"`
	Result              DeployArtifactResult   `json:"result"`
	RunAs               string                 `json:"run_as"`
	Deploy              bool                   `json:"-"`
}

func (da DeployArtifact) EvalImageTag() string {
	imageTag := da.ImageTag
	versionSplit := strings.Split(da.AvailableVersion, ".")
	replacementMap := make(map[string]string)
	replacementMap["$version"] = da.AvailableVersion
	for i, x := range versionSplit {
		replacementMap[fmt.Sprintf("$%d", i+1)] = x
	}
	for k, v := range replacementMap {
		imageTag = strings.Replace(imageTag, k, v, -1)
	}
	return imageTag
}

type DeployService struct {
	*DeployArtifact
	ServiceID   int    `json:"service_id"`
	ServicePort int    `json:"service_port"`
	MetricsPort int    `json:"metrics_port"`
	ServiceName string `json:"service_name"`
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

// ArtifactDeployResultMap is used to convert the array of artifacts results into a map
type ArtifactDeployResultMap map[DeployArtifactResult]DeployServices

// TopResultMap converts the array of results into a map by result
func (ds DeployServices) TopResultMap() ArtifactDeployResultMap {
	result := make(ArtifactDeployResultMap)

	for _, svc := range ds {
		switch svc.Result {
		case DeployArtifactResultFailed:
			result[DeployArtifactResultFailed] = append(result[DeployArtifactResultFailed], svc)
		case DeployArtifactResultSuccess:
			result[DeployArtifactResultSuccess] = append(result[DeployArtifactResultSuccess], svc)
		case DeployArtifactResultNoop:
			result[DeployArtifactResultNoop] = append(result[DeployArtifactResultNoop], svc)
		}
	}

	return result
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
	return fmt.Sprintf("deploy-ns-%d", ns.ID)
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
	DeploymentID     uuid.UUID            `json:"deployment_id"`
	Namespace        *NamespaceRequest    `json:"namespace"`
	EnvironmentName  string               `json:"environment_name"`
	EnvironmentAlias string               `json:"environment_alias"`
	Services         DeployServices       `json:"services,omitempty"`
	Migrations       DeployMigrations     `json:"migrations,omitempty"`
	Messages         []string             `json:"messages,omitempty"`
	SchQueueUrl      string               `json:"-"`
	CallbackURL      string               `json:"callback_url"`
	Status           DeploymentPlanStatus `json:"status"`
}

func (ns *NSDeploymentPlan) NothingToDeploy() bool {
	if len(ns.Services) == 0 && len(ns.Migrations) == 0 {
		return true
	}
	return false
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

	return false
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

	return false
}

func (ns *NSDeploymentPlan) Message(format string, a ...interface{}) {
	ns.Messages = append(ns.Messages, fmt.Sprintf(format, a...))
}

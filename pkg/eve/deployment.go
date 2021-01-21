package eve

import (
	"fmt"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"

	"gitlab.unanet.io/devops/eve/internal/data"
	"gitlab.unanet.io/devops/eve/pkg/queue"
)

const (
	// ArtifactoryFeedTypeDocker is exposed in eve (and not used) but used in eve-sch
	// ask Casey why this is? :)
	ArtifactoryFeedTypeDocker = "docker"
)

type PlanType string

func (t PlanType) Command() string {
	switch t {
	case DeploymentPlanTypeRestart:
		return queue.CommandRestartNamespace
	default:
		return queue.CommandDeployNamespace
	}
}

type DeploymentState string

const (
	DeploymentStateQueued    DeploymentState = "queued"
	DeploymentStateScheduled DeploymentState = "scheduled"
	DeploymentStateCompleted DeploymentState = "completed"
	DeploymentStateUnknown   DeploymentState = "unknown"
)

func ParseDeploymentState(value data.DeploymentState) DeploymentState {
	switch value {
	case data.DeploymentStateQueued:
		return DeploymentStateQueued
	case data.DeploymentStateScheduled:
		return DeploymentStateScheduled
	case data.DeploymentStateCompleted:
		return DeploymentStateCompleted
	default:
		return DeploymentStateUnknown
	}
}

const (
	DeploymentPlanTypeApplication PlanType = "application"
	DeploymentPlanTypeJob         PlanType = "job"
	DeploymentPlanTypeRestart     PlanType = "restart"
)

type DeployArtifactResult string

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
	DeploymentPlanStatusMessage  DeploymentPlanStatus = "message"
)

func (dps DeploymentPlanStatus) String() string {
	return string(dps)
}

type DeployArtifact struct {
	ArtifactID          int                  `json:"artifact_id"`
	ArtifactName        string               `json:"artifact_name"`
	RequestedVersion    string               `json:"requested_version"`
	DeployedVersion     string               `json:"deployed_version"`
	AvailableVersion    string               `json:"available_version"`
	ServiceAccount      string               `json:"service_account"`
	ImageTag            string               `json:"image_tag"`
	Labels              MetadataField        `json:"labels"`
	Annotations         MetadataField        `json:"annotations"`
	Metadata            MetadataField        `json:"metadata"`
	ArtifactoryFeed     string               `json:"artifactory_feed"`
	ArtifactoryPath     string               `json:"artifactory_path"`
	ArtifactFnPtr       string               `json:"artifact_fn"`
	ArtifactoryFeedType string               `json:"artifactory_feed_type"`
	Result              DeployArtifactResult `json:"result"`
	ExitCode            int                  `json:"exit_code"`
	RunAs               int                  `json:"run_as"`
	Deploy              bool                 `json:"-"`
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
	ServiceID        int    `json:"service_id"`
	ServicePort      int    `json:"service_port"`
	MetricsPort      int    `json:"metrics_port"`
	ServiceName      string `json:"service_name"`
	StickySessions   bool   `json:"sticky_sessions"`
	NodeGroup        string `json:"node_group"`
	Count            int    `json:"count"`
	LivelinessProbe  []byte `json:"liveliness_probe"`
	ReadinessProbe   []byte `json:"readiness_probe"`
	Autoscaling      []byte `json:"autoscaling"`
	PodResource      []byte `json:"pod_resource"`
	SuccessExitCodes string `json:"success_exit_codes"`
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

// ArtifactDeployResultMap is used to convert the array of deploy artifact results into a map
type ArtifactJobResultMap map[DeployArtifactResult]DeployJobs

// ToResultMap converts the array of results into a map by result
func (dj DeployJobs) ToResultMap() ArtifactJobResultMap {
	result := make(ArtifactJobResultMap)

	for _, job := range dj {
		switch job.Result {
		case DeployArtifactResultFailed:
			result[DeployArtifactResultFailed] = append(result[DeployArtifactResultFailed], job)
		case DeployArtifactResultSuccess:
			result[DeployArtifactResultSuccess] = append(result[DeployArtifactResultSuccess], job)
		case DeployArtifactResultNoop:
			result[DeployArtifactResultNoop] = append(result[DeployArtifactResultNoop], job)
		}
	}

	return result
}

// ArtifactDeployResultMap is used to convert the array of deploy artifact results into a map
type ArtifactServiceResultMap map[DeployArtifactResult]DeployServices

// ToResultMap converts the array of results into a map by result
func (ds DeployServices) ToResultMap() ArtifactServiceResultMap {
	result := make(ArtifactServiceResultMap)

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

type DeployJob struct {
	*DeployArtifact
	JobID            int    `json:"job_id"`
	JobName          string `json:"job_name"`
	NodeGroup        string `json:"node_group"`
	SuccessExitCodes string `json:"success_exit_codes"`
}

type DeployJobs []*DeployJob

func (dj DeployJobs) ToDeploy() DeployJobs {
	var list DeployJobs
	for _, x := range dj {
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
	Version     string `json:"version"`
}

func (ns *NamespaceRequest) GetQueueGroupID() string {
	return fmt.Sprintf("deploy-ns-%d", ns.ID)
}

func (ns *NamespaceRequest) SnakeAlias() string {
	return strings.Replace(ns.Alias, "-", "_", -1)
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
	DeploymentID      uuid.UUID            `json:"deployment_id"`
	Namespace         *NamespaceRequest    `json:"namespace"`
	EnvironmentName   string               `json:"environment_name"`
	EnvironmentAlias  string               `json:"environment_alias"`
	Services          DeployServices       `json:"services,omitempty"`
	Jobs              DeployJobs           `json:"jobs,omitempty"`
	Messages          []string             `json:"messages,omitempty"`
	SchQueueUrl       string               `json:"-"`
	CallbackURL       string               `json:"callback_url"`
	Status            DeploymentPlanStatus `json:"status"`
	MetadataOverrides MetadataField        `json:"metadata_overrides"`
	Type              PlanType             `json:"type"`
}

// DeploymentPlanType is a helper method to know what type of deployment plan (application,job,migration,restart)
func (ns *NSDeploymentPlan) DeploymentPlanType() string {
	return string(ns.Type)
}

func (ns *NSDeploymentPlan) NothingToDeploy() bool {
	if len(ns.Services) == 0 && len(ns.Jobs) == 0 {
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

	for _, x := range ns.Jobs {
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

	for _, x := range ns.Jobs {
		if x.Result == DeployArtifactResultFailed {
			return true
		}
	}

	return false
}

func (ns *NSDeploymentPlan) Message(format string, a ...interface{}) {
	ns.Messages = append(ns.Messages, fmt.Sprintf(format, a...))
}

func ToDeployment(d data.Deployment) Deployment {
	return Deployment{
		ID:            d.ID,
		EnvironmentID: d.EnvironmentID,
		NamespaceID:   d.NamespaceID,
		User:          d.User,
		State:         ParseDeploymentState(d.State),
		CreatedAt:     d.CreatedAt.Time,
		UpdatedAt:     d.UpdatedAt.Time,
	}
}

type Deployment struct {
	ID            uuid.UUID       `json:"id"`
	EnvironmentID int             `json:"environment_id"`
	NamespaceID   int             `json:"namespace_id"`
	User          string          `json:"user"`
	State         DeploymentState `json:"state"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
}

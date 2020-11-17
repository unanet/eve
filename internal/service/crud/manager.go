package crud

import (
	"context"

	"gitlab.unanet.io/devops/eve/internal/data"
	"gitlab.unanet.io/devops/eve/pkg/json"
)

type Repo interface {
	Environments(ctx context.Context) (data.Environments, error)
	EnvironmentByID(ctx context.Context, id int) (*data.Environment, error)
	EnvironmentByName(ctx context.Context, name string) (*data.Environment, error)
	UpdateEnvironment(ctx context.Context, environment *data.Environment) error

	Namespaces(ctx context.Context) (data.Namespaces, error)
	NamespaceByID(ctx context.Context, id int) (*data.Namespace, error)
	NamespaceByName(ctx context.Context, name string) (*data.Namespace, error)
	NamespacesByEnvironmentID(ctx context.Context, environmentID int) (data.Namespaces, error)
	NamespacesByEnvironmentName(ctx context.Context, environmentName string) (data.Namespaces, error)
	UpdateNamespace(ctx context.Context, namespace *data.Namespace) error

	ServiceByID(ctx context.Context, id int) (*data.Service, error)
	ServiceByName(ctx context.Context, name string, namespace string) (*data.Service, error)
	ServicesByNamespaceID(ctx context.Context, namespaceID int) ([]data.Service, error)
	ServicesByNamespaceName(ctx context.Context, namespaceName string) ([]data.Service, error)
	UpdateService(ctx context.Context, service *data.Service) error
	UpdateServiceMetadata(ctx context.Context, serviceID int, metadata map[string]interface{}) error
	DeleteServiceMetadataKey(ctx context.Context, serviceID int, key string) error

	FeedByAliasAndType(ctx context.Context, alias, feedType string) (*data.Feed, error)
	NextFeedByPromotionOrderType(ctx context.Context, promotionOrder int, feedType string) (*data.Feed, error)
	PreviousFeedByPromotionOrderType(ctx context.Context, promotionOrder int, feedType string) (*data.Feed, error)

	ArtifactByName(ctx context.Context, name string) (*data.Artifact, error)
	ArtifactByID(ctx context.Context, id int) (*data.Artifact, error)

	PodAutoscaleMap(ctx context.Context, serviceID, environmentID, namespaceID int) ([]data.PodAutoscaleMap, error)
	PodAutoscaleStacked(pams []data.PodAutoscaleMap) (json.Text, error)
	NamespacePodAutoscaleMap(ctx context.Context, namespaceID int) ([]data.PodAutoscaleMap, error)
	EnvironmentPodAutoscaleMap(ctx context.Context, environmentID int) ([]data.PodAutoscaleMap, error)
	HydrateDeployServicePodAutoscale(ctx context.Context, svc data.DeployService) (json.Text, error)

	PodResourcesMap(ctx context.Context, serviceID, environmentID, namespaceID, artifactID int) ([]data.PodResourcesMap, error)
	PodResourcesStacked(prms []data.PodResourcesMap) (json.Text, error)
	NamespacePodResourcesMap(ctx context.Context, namespaceID int) ([]data.PodResourcesMap, error)
	EnvironmentPodResourcesMap(ctx context.Context, environmentID int) ([]data.PodResourcesMap, error)
	ArtifactPodResourcesMap(ctx context.Context, artifactID, environmentID, namespaceID int) ([]data.PodResourcesMap, error)
	HydrateDeployServicePodResource(ctx context.Context, svc data.DeployService) (json.Text, error)

	UpsertMergeMetadata(ctx context.Context, m *data.Metadata) error
	UpsertMetadata(ctx context.Context, m *data.Metadata) error
	UpsertMetadataServiceMap(ctx context.Context, msm *data.MetadataServiceMap) error
	Metadata(ctx context.Context) ([]data.Metadata, error)
	ServiceMetadata(ctx context.Context, serviceID int) ([]data.MetadataService, error)
	DeleteMetadataKey(ctx context.Context, metadataID int, key string) (*data.Metadata, error)
	DeleteMetadata(ctx context.Context, metadataID int) error
	GetMetadata(ctx context.Context, metadataID int) (*data.Metadata, error)
	GetMetadataByDescription(ctx context.Context, description string) (*data.Metadata, error)
	DeleteMetadataServiceMap(ctx context.Context, metadataID int, mapDescription string) error
}

func NewManager(r Repo) *Manager {
	return &Manager{
		repo: r,
	}
}

type Manager struct {
	repo Repo
}

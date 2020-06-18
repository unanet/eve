package crud

import (
	"context"

	"gitlab.unanet.io/devops/eve/internal/data"
)

type Repo interface {
	Environments(ctx context.Context) (data.Environments, error)
	EnvironmentByID(ctx context.Context, id int) (*data.Environment, error)
	EnvironmentByName(ctx context.Context, name string) (*data.Environment, error)

	Namespaces(ctx context.Context) (data.Namespaces, error)
	NamespaceByID(ctx context.Context, id int) (*data.Namespace, error)
	NamespaceByName(ctx context.Context, name string) (*data.Namespace, error)
	NamespacesByEnvironmentID(ctx context.Context, environmentID int) (data.Namespaces, error)
	NamespacesByEnvironmentName(ctx context.Context, environmentName string) (data.Namespaces, error)

	ServiceByID(ctx context.Context, id int) (*data.Service, error)
	ServiceByName(ctx context.Context, name string, namespace string) (*data.Service, error)
	ServicesByNamespaceID(ctx context.Context, namespaceID int) ([]data.Service, error)
	ServicesByNamespaceName(ctx context.Context, namespaceName string) ([]data.Service, error)
	UpdateService(ctx context.Context, service *data.Service) error
	UpdateServiceMetadataKey(ctx context.Context, serviceID int, key string, value string) error
	DeleteServiceMetadataKey(ctx context.Context, serviceID int, key string) error
}

func NewManager(r Repo) *Manager {
	return &Manager{
		repo: r,
	}
}

type Manager struct {
	repo Repo
}

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
}

func NewManager(r Repo) *Manager {
	return &Manager{
		repo: r,
	}
}

type Manager struct {
	repo Repo
}

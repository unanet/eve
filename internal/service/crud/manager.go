package crud

import (
	"context"

	"gitlab.unanet.io/devops/eve/internal/data"
)

type Repo interface {
	Environments(ctx context.Context) (data.Environments, error)
}

func NewManager(r Repo) *Manager {
	return &Manager{
		repo: r,
	}
}

type Manager struct {
	repo Repo
}

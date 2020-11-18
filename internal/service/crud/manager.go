package crud

import (
	"gitlab.unanet.io/devops/eve/internal/data"
)

func NewManager(r *data.Repo) *Manager {
	return &Manager{
		repo: r,
	}
}

type Manager struct {
	repo *data.Repo
}

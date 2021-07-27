package crud

import (
	"context"

	uuid "github.com/satori/go.uuid"
	"github.com/unanet/go/pkg/errors"

	"github.com/unanet/eve/internal/service"
	"github.com/unanet/eve/pkg/eve"
)

func (m *Manager) Deployment(ctx context.Context, id string) (*eve.Deployment, error) {
	uID, err := uuid.FromString(id)
	if err != nil {
		return nil, errors.NewRestError(400, "invalid deployment id")
	}

	d, err := m.repo.DeploymentByID(ctx, uID)
	if err != nil {
		return nil, service.CheckForNotFoundError(err)
	}

	deployment := eve.ToDeployment(*d)
	return &deployment, nil
}

package service

import (
	"context"

	uuid "github.com/satori/go.uuid"

	"gitlab.unanet.io/devops/eve/internal/cloud/queue"
	"gitlab.unanet.io/devops/eve/internal/data"
	"gitlab.unanet.io/devops/eve/pkg/errors"
)

type DeploymentQueue struct {
	worker   QueueWorker
	repo     DeploymentQueueRepo
	schQueue QWriter
}

type DeploymentQueueRepo interface {
	Deployment(ctx context.Context, id uuid.UUID) (*data.Deployment, error)
}

func NewDeploymentQueue(worker QueueWorker, repo DeploymentQueueRepo, schQueue QWriter) *DeploymentQueue {
	return &DeploymentQueue{
		worker:   worker,
		repo:     repo,
		schQueue: schQueue,
	}
}

func (dq *DeploymentQueue) Start() {
	go func() {
		dq.worker.Start(queue.HandlerFunc(dq.handleMessage))
	}()
}

func (dq *DeploymentQueue) Stop() {
	dq.worker.Stop()
}

func (dq *DeploymentQueue) handleMessage(ctx context.Context, m *queue.M) error {
	deployment, err := dq.repo.Deployment(ctx, m.ID)
	if err != nil {
		return errors.Wrap(err)
	}
	switch deployment.State {
	// This means it hasn't been send to the scheduler yet
	case data.DeploymentStateQueued:

	// This means it came back from the scheduler
	case data.DeploymentStateScheduled:

	// WTF we shouldn't hit this case in here
	case data.DeploymentStateCompleted:

	// Also we should hit this case
	default:
	}

	return nil
}

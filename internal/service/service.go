package service

import (
	"context"

	"gitlab.unanet.io/devops/eve/pkg/queue"
)

type QWriter interface {
	Message(ctx context.Context, m *queue.M) error
}

type QueueWorker interface {
	Start(queue.Handler)
	Stop()
	DeleteMessage(ctx context.Context, m *queue.M) error
	// Message sends a message to a different queue given a url, not this one
	Message(ctx context.Context, qUrl string, m *queue.M) error
}

type StringList []string

func (s StringList) Contains(value string) bool {
	for _, a := range s {
		if a == value {
			return true
		}
	}
	return false
}

package service

import (
	"gitlab.unanet.io/devops/eve/internal/cloud/queue"
)

type QWriter interface {
	Message(m *queue.M) error
}

type QueueWorker interface {
	Start(queue.Handler)
	Stop()
	DeleteMessage(m *queue.M) error
}

type M map[string]interface{}

// mergeKeys stomps on the keys in the left map if they exist in the right map
func MergeMetadata(left, right M) M {
	for key, rightVal := range right {
		left[key] = rightVal
	}
	return left
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

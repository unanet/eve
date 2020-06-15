package crud

import (
	"gitlab.unanet.io/devops/eve/pkg/json"
)

type Environment struct {
	ID       int       `json:"id"`
	Name     string    `json:"name"`
	Alias    string    `json:"alias"`
	Metadata json.Text `json:"metadata,omitempty"`
}

func (m *Manager) Environments() []Environment {
	return nil
}

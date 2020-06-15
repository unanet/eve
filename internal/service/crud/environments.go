package crud

import (
	"context"

	"gitlab.unanet.io/devops/eve/internal/data"
	"gitlab.unanet.io/devops/eve/pkg/errors"
	"gitlab.unanet.io/devops/eve/pkg/eve"
)

func fromDataEnvironment(environment data.Environment) eve.Environment {
	return eve.Environment{
		ID:       environment.ID,
		Name:     environment.Name,
		Alias:    environment.Alias,
		Metadata: environment.Metadata.AsMap(),
	}
}

func fromDataEnvironments(environments data.Environments) []eve.Environment {
	var list []eve.Environment
	for _, x := range environments {
		list = append(list, fromDataEnvironment(x))
	}
	return list
}

func (m *Manager) Environments(ctx context.Context) ([]eve.Environment, error) {
	dataEnvironments, err := m.repo.Environments(ctx)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	return fromDataEnvironments(dataEnvironments), nil
}

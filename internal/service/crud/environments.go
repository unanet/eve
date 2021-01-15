package crud

import (
	"context"
	"strconv"

	"gitlab.unanet.io/devops/go/pkg/errors"

	"gitlab.unanet.io/devops/eve/internal/data"
	"gitlab.unanet.io/devops/eve/internal/service"
	"gitlab.unanet.io/devops/eve/pkg/eve"
)

func fromDataEnvironment(environment data.Environment) eve.Environment {
	return eve.Environment{
		ID:          environment.ID,
		Name:        environment.Name,
		Alias:       environment.Alias,
		Description: environment.Description,
		UpdatedAt:   environment.UpdatedAt.Time,
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

func (m *Manager) Environment(ctx context.Context, id string) (*eve.Environment, error) {
	var dEnvironment *data.Environment
	if intID, err := strconv.Atoi(id); err == nil {
		dEnvironment, err = m.repo.EnvironmentByID(ctx, intID)
		if err != nil {
			return nil, service.CheckForNotFoundError(err)
		}
	} else {
		dEnvironment, err = m.repo.EnvironmentByName(ctx, id)
		if err != nil {
			return nil, service.CheckForNotFoundError(err)
		}
	}

	environment := fromDataEnvironment(*dEnvironment)
	return &environment, nil
}

func (m *Manager) UpdateEnvironment(ctx context.Context, e *eve.Environment) (*eve.Environment, error) {
	dEnvironment := toDataEnvironment(*e)
	err := m.repo.UpdateEnvironment(ctx, &dEnvironment)
	if err != nil {
		return nil, service.CheckForNotFoundError(err)
	}

	e2 := fromDataEnvironment(dEnvironment)
	return &e2, nil
}

func toDataEnvironment(environment eve.Environment) data.Environment {
	return data.Environment{
		ID:          environment.ID,
		Name:        environment.Name,
		Alias:       environment.Alias,
		Description: environment.Description,
	}
}

package crud

import (
	"context"
	"strconv"

	"gitlab.unanet.io/devops/eve/internal/data"
	"gitlab.unanet.io/devops/eve/pkg/errors"
	"gitlab.unanet.io/devops/eve/pkg/eve"
)

func fromDataNamespace(namespace data.Namespace) eve.Namespace {
	return eve.Namespace{
		ID:                 namespace.ID,
		Name:               namespace.Name,
		Alias:              namespace.Alias,
		EnvironmentID:      namespace.EnvironmentID,
		EnvironmentName:    namespace.EnvironmentName,
		RequestedVersion:   namespace.RequestedVersion,
		ExplicitDeployOnly: namespace.ExplicitDeployOnly,
		ClusterID:          namespace.ClusterID,
		Metadata:           namespace.Metadata.AsMap(),
		CreatedAt:          namespace.CreatedAt.Time,
		UpdatedAt:          namespace.UpdatedAt.Time,
	}
}

func fromDataNamespaces(namespaces data.Namespaces) []eve.Namespace {
	var list []eve.Namespace
	for _, x := range namespaces {
		list = append(list, fromDataNamespace(x))
	}
	return list
}

func (m *Manager) Namespaces(ctx context.Context) ([]eve.Namespace, error) {
	dataNamespaces, err := m.repo.Namespaces(ctx)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	return fromDataNamespaces(dataNamespaces), nil
}

func (m *Manager) NamespacesByEnvironment(ctx context.Context, environmentID string) ([]eve.Namespace, error) {
	var dNamespaces []data.Namespace
	if intID, err := strconv.Atoi(environmentID); err == nil {
		dNamespaces, err = m.repo.NamespacesByEnvironmentID(ctx, intID)
		if err != nil {
			return nil, errors.Wrap(err)
		}
	} else {
		dNamespaces, err = m.repo.NamespacesByEnvironmentName(ctx, environmentID)
		if err != nil {
			return nil, errors.Wrap(err)
		}
	}

	return fromDataNamespaces(dNamespaces), nil
}

func (m *Manager) Namespace(ctx context.Context, id string) (*eve.Namespace, error) {
	var dNamespace *data.Namespace
	if intID, err := strconv.Atoi(id); err == nil {
		dNamespace, err = m.repo.NamespaceByID(ctx, intID)
		if err != nil {
			return nil, errors.Wrap(err)
		}
	} else {
		dNamespace, err = m.repo.NamespaceByName(ctx, id)
		if err != nil {
			return nil, errors.Wrap(err)
		}
	}

	namespace := fromDataNamespace(*dNamespace)
	return &namespace, nil
}

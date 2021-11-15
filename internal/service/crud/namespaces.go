package crud

import (
	"context"
	"strconv"

	"github.com/unanet/go/pkg/errors"

	"github.com/unanet/eve/internal/data"
	"github.com/unanet/eve/internal/service"
	"github.com/unanet/eve/pkg/eve"
)

func fromDataNamespace(namespace data.Namespace) eve.Namespace {
	return eve.Namespace{
		ID:               namespace.ID,
		Name:             namespace.Name,
		Alias:            namespace.Alias,
		EnvironmentID:    namespace.EnvironmentID,
		EnvironmentName:  namespace.EnvironmentName,
		RequestedVersion: namespace.RequestedVersion,
		ExplicitDeploy:   namespace.ExplicitDeploy,
		ClusterID:        namespace.ClusterID,
		CreatedAt:        namespace.CreatedAt.Time,
		UpdatedAt:        namespace.UpdatedAt.Time,
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
			return nil, service.CheckForNotFoundError(err)
		}
	} else {
		dNamespaces, err = m.repo.NamespacesByEnvironmentName(ctx, environmentID)
		if err != nil {
			return nil, service.CheckForNotFoundError(err)
		}
	}

	return fromDataNamespaces(dNamespaces), nil
}

// Namespace id can be an int or the namespace name
func (m *Manager) Namespace(ctx context.Context, id string) (*eve.Namespace, error) {
	var dNamespace *data.Namespace
	if intID, err := strconv.Atoi(id); err == nil {
		dNamespace, err = m.repo.NamespaceByID(ctx, intID)
		if err != nil {
			return nil, service.CheckForNotFoundError(err)
		}
	} else {
		dNamespace, err = m.repo.NamespaceByName(ctx, id)
		if err != nil {
			return nil, service.CheckForNotFoundError(err)
		}
	}

	namespace := fromDataNamespace(*dNamespace)
	return &namespace, nil
}

func (m *Manager) UpdateNamespace(ctx context.Context, n *eve.Namespace) (*eve.Namespace, error) {
	dNamespace := toDataNamespace(*n)
	err := m.repo.UpdateNamespace(ctx, &dNamespace)
	if err != nil {
		return nil, service.CheckForNotFoundError(err)
	}

	n2 := fromDataNamespace(dNamespace)
	return &n2, nil
}

func (m *Manager) CreateNamespace(ctx context.Context, model *eve.Namespace) error {

	dbNamespace := toDataNamespace(*model)
	if err := m.repo.CreateNamespace(ctx, &dbNamespace); err != nil {
		return errors.Wrap(err)
	}

	model.ID = dbNamespace.ID

	return nil
}

func (m *Manager) DeleteNamespace(ctx context.Context, id int) (err error) {

	if err := m.repo.DeleteNamespace(ctx, id); err != nil {
		return service.CheckForNotFoundError(err)
	}

	return nil
}

func toDataNamespace(namespace eve.Namespace) data.Namespace {
	return data.Namespace{
		ID:               namespace.ID,
		Name:             namespace.Name,
		Alias:            namespace.Alias,
		EnvironmentID:    namespace.EnvironmentID,
		EnvironmentName:  namespace.EnvironmentName,
		RequestedVersion: namespace.RequestedVersion,
		ExplicitDeploy:   namespace.ExplicitDeploy,
		ClusterID:        namespace.ClusterID,
	}
}

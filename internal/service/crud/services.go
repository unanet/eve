package crud

import (
	"context"
	"strconv"

	"gitlab.unanet.io/devops/eve/internal/data"
	"gitlab.unanet.io/devops/eve/internal/service"
	"gitlab.unanet.io/devops/eve/pkg/eve"
)

func fromDataService(service data.Service) eve.Service {
	return eve.Service{
		ID:              service.ID,
		NamespaceID:     service.NamespaceID,
		NamespaceName:   service.NamespaceName,
		ArtifactID:      service.ArtifactID,
		ArtifactName:    service.ArtifactName,
		OverrideVersion: service.OverrideVersion.String,
		DeployedVersion: service.DeployedVersion.String,
		Metadata:        service.Metadata.AsMap(),
		CreatedAt:       service.CreatedAt.Time,
		UpdatedAt:       service.UpdatedAt.Time,
		Name:            service.Name,
		StickySessions:  service.StickySessions,
		Count:           service.Count,
	}
}

func fromDataServices(services []data.Service) []eve.Service {
	var list []eve.Service
	for _, x := range services {
		list = append(list, fromDataService(x))
	}
	return list
}

func (m *Manager) ServicesByNamespace(ctx context.Context, namespaceID string) ([]eve.Service, error) {
	var dServices []data.Service
	if intID, err := strconv.Atoi(namespaceID); err == nil {
		dServices, err = m.repo.ServicesByNamespaceID(ctx, intID)
		if err != nil {
			return nil, service.CheckForNotFoundError(err)
		}
	} else {
		dServices, err = m.repo.ServicesByNamespaceName(ctx, namespaceID)
		if err != nil {
			return nil, service.CheckForNotFoundError(err)
		}
	}

	return fromDataServices(dServices), nil
}

func (m *Manager) Service(ctx context.Context, id string, namespace string) (*eve.Service, error) {
	var dService *data.Service
	if intID, err := strconv.Atoi(id); err == nil {
		dService, err = m.repo.ServiceByID(ctx, intID)
		if err != nil {
			return nil, service.CheckForNotFoundError(err)
		}
	} else {
		dService, err = m.repo.ServiceByName(ctx, id, namespace)
		if err != nil {
			return nil, service.CheckForNotFoundError(err)
		}
	}

	service := fromDataService(*dService)
	return &service, nil
}

func (m *Manager) UpdateService(ctx context.Context, serviceID int, version string) error {
	return nil
}

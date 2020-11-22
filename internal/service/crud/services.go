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
		CreatedAt:       service.CreatedAt.Time,
		UpdatedAt:       service.UpdatedAt.Time,
		Name:            service.Name,
		StickySessions:  service.StickySessions,
		NodeGroup:       service.NodeGroup,
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

func toDataService(service eve.Service) data.Service {
	s := data.Service{
		ID:             service.ID,
		NamespaceID:    service.NamespaceID,
		NamespaceName:  service.NamespaceName,
		ArtifactID:     service.ArtifactID,
		ArtifactName:   service.ArtifactName,
		Name:           service.Name,
		StickySessions: service.StickySessions,
		NodeGroup:      service.NodeGroup,
		Count:          service.Count,
	}

	if service.OverrideVersion != "" {
		s.OverrideVersion.String = service.OverrideVersion
		s.OverrideVersion.Valid = true
	}

	if service.DeployedVersion != "" {
		s.DeployedVersion.String = service.DeployedVersion
		s.DeployedVersion.Valid = true
	}

	return s
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

	s := fromDataService(*dService)
	return &s, nil
}

func (m *Manager) UpdateService(ctx context.Context, s *eve.Service) (*eve.Service, error) {
	dService := toDataService(*s)

	err := m.repo.UpdateService(ctx, &dService)
	if err != nil {
		return nil, service.CheckForNotFoundError(err)
	}

	s2 := fromDataService(dService)
	return &s2, nil
}

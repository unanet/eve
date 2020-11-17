package crud

import (
	"context"
	"strconv"

	"gitlab.unanet.io/devops/eve/internal/data"
	"gitlab.unanet.io/devops/eve/internal/service"
	"gitlab.unanet.io/devops/eve/pkg/errors"
	"gitlab.unanet.io/devops/eve/pkg/eve"
	"gitlab.unanet.io/devops/eve/pkg/json"
)

func toDataMetadataServiceMap(m eve.MetadataServiceMap) data.MetadataServiceMap {
	return data.MetadataServiceMap{
		Description:   m.Description,
		MetadataID:    m.MetadataID,
		EnvironmentID: m.EnvironmentID,
		ArtifactID:    m.ArtifactID,
		NamespaceID:   m.NamespaceID,
		ServiceID:     m.ServiceID,
		StackingOrder: m.StackingOrder,
	}
}

func toDataMetadata(m eve.Metadata) data.Metadata {
	return data.Metadata{
		ID:          m.ID,
		Description: m.Description,
		Value:       json.FromMap(m.Value),
	}
}

func fromDataMetadata(m data.Metadata) eve.Metadata {
	return eve.Metadata{
		ID:          m.ID,
		Description: m.Description,
		Value:       m.Value.AsMap(),
		CreatedAt:   m.CreatedAt.Time,
		UpdatedAt:   m.UpdatedAt.Time,
	}
}

func fromDataMetadataList(metadata []data.Metadata) []eve.Metadata {
	var list []eve.Metadata
	for _, x := range metadata {
		list = append(list, fromDataMetadata(x))
	}
	return list
}

func fromDataMetadataServiceToMetadata(m data.MetadataService) eve.Metadata {
	return eve.Metadata{
		ID:          m.MetadataID,
		Description: m.MetadataDescription,
		Value:       m.Metadata.AsMap(),
		CreatedAt:   m.CreatedAt.Time,
		UpdatedAt:   m.UpdatedAt.Time,
	}
}

func fromDataMetadataServiceListToMetadataList(m []data.MetadataService) []eve.Metadata {
	var list []eve.Metadata
	for _, x := range m {
		list = append(list, fromDataMetadataServiceToMetadata(x))
	}
	return list
}

func (m Manager) Metadata(ctx context.Context, serviceID string, namespaceID string) ([]eve.Metadata, error) {
	if len(serviceID) != 0 {
		s, err := m.Service(ctx, serviceID, namespaceID)
		if err != nil {
			return nil, service.CheckForNotFoundError(err)
		}
		sMetadata, err := m.repo.ServiceMetadata(ctx, s.ID)
		if err != nil {
			return nil, errors.Wrap(err)
		}

		return fromDataMetadataServiceListToMetadataList(sMetadata), nil
	} else {
		metadata, err := m.repo.Metadata(ctx)
		if err != nil {
			return nil, errors.Wrap(err)
		}
		return fromDataMetadataList(metadata), nil
	}
}

func (m Manager) CreateMetadata(ctx context.Context, metadata *eve.Metadata) error {
	dataMetadata := toDataMetadata(*metadata)
	err := m.repo.UpsertMetadata(ctx, &dataMetadata)
	if err != nil {
		return errors.Wrap(err)
	}

	metadata.UpdatedAt = dataMetadata.UpdatedAt.Time
	metadata.CreatedAt = dataMetadata.CreatedAt.Time
	metadata.ID = dataMetadata.ID
	return nil
}

func (m Manager) UpsertMergeMetadata(ctx context.Context, metadata *eve.Metadata) error {
	dataMetadata := toDataMetadata(*metadata)
	err := m.repo.UpsertMergeMetadata(ctx, &dataMetadata)
	if err != nil {
		return errors.Wrap(err)
	}

	metadata.UpdatedAt = dataMetadata.UpdatedAt.Time
	metadata.CreatedAt = dataMetadata.CreatedAt.Time
	metadata.ID = dataMetadata.ID
	return nil
}

func (m Manager) DeleteMetadataKey(ctx context.Context, id int, key string) (eve.Metadata, error) {
	metadata, err := m.repo.DeleteMetadataKey(ctx, id, key)
	if err != nil {
		return eve.Metadata{}, service.CheckForNotFoundError(err)
	}

	return fromDataMetadata(*metadata), nil
}

func (m *Manager) GetMetadata(ctx context.Context, id string) (*eve.Metadata, error) {
	var metadata *data.Metadata
	if intID, err := strconv.Atoi(id); err == nil {
		metadata, err = m.repo.GetMetadata(ctx, intID)
		if err != nil {
			return nil, service.CheckForNotFoundError(err)
		}
	} else {
		metadata, err = m.repo.GetMetadataByDescription(ctx, id)
		if err != nil {
			return nil, service.CheckForNotFoundError(err)
		}
	}

	r := fromDataMetadata(*metadata)
	return &r, nil
}

func (m *Manager) DeleteMetadata(ctx context.Context, id int) error {
	err := m.repo.DeleteMetadata(ctx, id)
	if err != nil {
		return service.CheckForNotFoundError(err)
	}
	return nil
}

func (m *Manager) UpsertMetadataServiceMap(ctx context.Context, serviceMap *eve.MetadataServiceMap) error {
	dataMetadataServiceMap := toDataMetadataServiceMap(*serviceMap)
	err := m.repo.UpsertMetadataServiceMap(ctx, &dataMetadataServiceMap)
	if err != nil {
		return errors.Wrap(err)
	}

	serviceMap.UpdatedAt = dataMetadataServiceMap.UpdatedAt.Time
	serviceMap.CreatedAt = dataMetadataServiceMap.CreatedAt.Time
	return nil
}

func (m *Manager) DeleteMetadataServiceMap(ctx context.Context, metadataID int, description string) error {
	err := m.repo.DeleteMetadataServiceMap(ctx, metadataID, description)
	if err != nil {
		return service.CheckForNotFoundError(err)
	}

	return nil
}

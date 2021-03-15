package crud

import (
	"context"
	"strconv"

	"gitlab.unanet.io/devops/go/pkg/errors"
	"gitlab.unanet.io/devops/go/pkg/json"

	"gitlab.unanet.io/devops/eve/internal/data"
	"gitlab.unanet.io/devops/eve/internal/service"
	"gitlab.unanet.io/devops/eve/pkg/eve"
)

func toDataLabelServiceMap(m eve.LabelServiceMap) data.LabelServiceMap {
	dm := data.LabelServiceMap{
		Description:   m.Description,
		LabelID:       m.LabelID,
		StackingOrder: m.StackingOrder,
	}

	if m.EnvironmentID != 0 {
		dm.EnvironmentID.Int32 = int32(m.EnvironmentID)
		dm.EnvironmentID.Valid = true
	}

	if m.ArtifactID != 0 {
		dm.ArtifactID.Int32 = int32(m.ArtifactID)
		dm.ArtifactID.Valid = true
	}

	if m.NamespaceID != 0 {
		dm.NamespaceID.Int32 = int32(m.NamespaceID)
		dm.NamespaceID.Valid = true
	}

	if m.ServiceID != 0 {
		dm.ServiceID.Int32 = int32(m.ServiceID)
		dm.ServiceID.Valid = true
	}

	return dm
}

func toDataLabelJobMap(m eve.LabelJobMap) data.LabelJobMap {
	dm := data.LabelJobMap{
		Description:   m.Description,
		LabelID:       m.LabelID,
		StackingOrder: m.StackingOrder,
	}

	if m.EnvironmentID != 0 {
		dm.EnvironmentID.Int32 = int32(m.EnvironmentID)
		dm.EnvironmentID.Valid = true
	}

	if m.ArtifactID != 0 {
		dm.ArtifactID.Int32 = int32(m.ArtifactID)
		dm.ArtifactID.Valid = true
	}

	if m.NamespaceID != 0 {
		dm.NamespaceID.Int32 = int32(m.NamespaceID)
		dm.NamespaceID.Valid = true
	}

	if m.JobID != 0 {
		dm.JobID.Int32 = int32(m.JobID)
		dm.JobID.Valid = true
	}

	return dm
}

func toDataLabel(m eve.Label) data.Label {
	return data.Label{
		ID:          m.ID,
		Description: m.Description,
		Data:        json.FromMapOrEmpty(m.Data),
	}
}

func fromDataLabel(m data.Label) eve.Label {
	return eve.Label{
		ID:          m.ID,
		Description: m.Description,
		Data:        m.Data.AsMapOrEmpty(),
		CreatedAt:   m.CreatedAt.Time,
		UpdatedAt:   m.UpdatedAt.Time,
	}
}

func fromDataLabelList(labels []data.Label) []eve.Label {
	var list []eve.Label
	for _, x := range labels {
		list = append(list, fromDataLabel(x))
	}
	return list
}

func fromDataLabelServiceMaps(m []data.LabelServiceMap) []eve.LabelServiceMap {
	var list []eve.LabelServiceMap
	for _, x := range m {
		list = append(list, fromDataLabelServiceMap(x))
	}
	return list
}

func fromDataLabelJobMap(m data.LabelJobMap) eve.LabelJobMap {
	return eve.LabelJobMap{
		Description:   m.Description,
		LabelID:       m.LabelID,
		EnvironmentID: int(m.EnvironmentID.Int32),
		ArtifactID:    int(m.ArtifactID.Int32),
		NamespaceID:   int(m.NamespaceID.Int32),
		JobID:         int(m.JobID.Int32),
		StackingOrder: m.StackingOrder,
		CreatedAt:     m.CreatedAt.Time,
		UpdatedAt:     m.UpdatedAt.Time,
	}
}

func fromDataLabelServiceMap(m data.LabelServiceMap) eve.LabelServiceMap {
	return eve.LabelServiceMap{
		Description:   m.Description,
		LabelID:       m.LabelID,
		EnvironmentID: int(m.EnvironmentID.Int32),
		ArtifactID:    int(m.ArtifactID.Int32),
		NamespaceID:   int(m.NamespaceID.Int32),
		ServiceID:     int(m.ServiceID.Int32),
		StackingOrder: m.StackingOrder,
		CreatedAt:     m.CreatedAt.Time,
		UpdatedAt:     m.UpdatedAt.Time,
	}
}

func fromDataLabelService(m data.LabelService) eve.LabelServiceMap {
	return eve.LabelServiceMap{
		Description:   m.MapDescription,
		LabelID:       m.LabelID,
		EnvironmentID: int(m.MapEnvironmentID.Int32),
		ArtifactID:    int(m.MapArtifactID.Int32),
		NamespaceID:   int(m.MapNamespaceID.Int32),
		ServiceID:     int(m.MapServiceID.Int32),
		StackingOrder: m.StackingOrder,
		CreatedAt:     m.CreatedAt.Time,
		UpdatedAt:     m.UpdatedAt.Time,
	}
}

func fromDataLabelServiceList(m []data.LabelService) []eve.LabelServiceMap {
	var list []eve.LabelServiceMap
	for _, x := range m {
		list = append(list, fromDataLabelService(x))
	}
	return list
}

func fromDataLabelJobMaps(maps []data.LabelJobMap) []eve.LabelJobMap {
	var list []eve.LabelJobMap
	for _, x := range maps {
		list = append(list, fromDataLabelJobMap(x))
	}
	return list
}

func (m Manager) Labels(ctx context.Context) ([]eve.Label, error) {
	labels, err := m.repo.Labels(ctx)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	return fromDataLabelList(labels), nil
}

func (m Manager) CreateLabel(ctx context.Context, label *eve.Label) error {
	dataLabel := toDataLabel(*label)
	err := m.repo.UpsertLabel(ctx, &dataLabel)
	if err != nil {
		return errors.Wrap(err)
	}

	label.UpdatedAt = dataLabel.UpdatedAt.Time
	label.CreatedAt = dataLabel.CreatedAt.Time
	label.ID = dataLabel.ID
	return nil
}

func (m Manager) UpsertMergeLabel(ctx context.Context, label *eve.Label) error {
	dataLabel := toDataLabel(*label)
	err := m.repo.UpsertMergeLabel(ctx, &dataLabel)
	if err != nil {
		return errors.Wrap(err)
	}

	label.UpdatedAt = dataLabel.UpdatedAt.Time
	label.CreatedAt = dataLabel.CreatedAt.Time
	label.ID = dataLabel.ID
	label.Data = dataLabel.Data.AsMapOrEmpty()
	return nil
}

func (m *Manager) GetLabel(ctx context.Context, id string) (*eve.Label, error) {
	var label *data.Label
	if intID, err := strconv.Atoi(id); err == nil {
		label, err = m.repo.GetLabel(ctx, intID)
		if err != nil {
			return nil, service.CheckForNotFoundError(err)
		}
	} else {
		label, err = m.repo.GetLabelByDescription(ctx, id)
		if err != nil {
			return nil, service.CheckForNotFoundError(err)
		}
	}

	r := fromDataLabel(*label)
	return &r, nil
}

func (m *Manager) JobLabelMaps(ctx context.Context, id int) ([]eve.LabelJobMap, error) {
	maps, err := m.repo.JobLabelMaps(ctx, id)
	if err != nil {
		return nil, service.CheckForNotFoundError(err)
	}

	return fromDataLabelJobMaps(maps), nil
}

func (m *Manager) ServiceLabelMaps(ctx context.Context, id int) ([]eve.LabelServiceMap, error) {
	maps, err := m.repo.ServiceLabels(ctx, id)
	if err != nil {
		return nil, service.CheckForNotFoundError(err)
	}

	return fromDataLabelServiceList(maps), nil
}

func (m *Manager) ServiceLabelMapsByLabelID(ctx context.Context, id int) ([]eve.LabelServiceMap, error) {
	maps, err := m.repo.ServiceLabelMapsByLabelID(ctx, id)
	if err != nil {
		return nil, service.CheckForNotFoundError(err)
	}

	return fromDataLabelServiceMaps(maps), nil
}

func (m *Manager) JobLabelMapsByLabelID(ctx context.Context, id int) ([]eve.LabelJobMap, error) {
	maps, err := m.repo.JobLabelMapsByLabelID(ctx, id)
	if err != nil {
		return nil, service.CheckForNotFoundError(err)
	}

	return fromDataLabelJobMaps(maps), nil
}

func (m *Manager) UpsertLabelJobMap(ctx context.Context, e *eve.LabelJobMap) error {
	dataLabelJobMap := toDataLabelJobMap(*e)
	err := m.repo.UpsertLabelJobMap(ctx, &dataLabelJobMap)
	if err != nil {
		return errors.Wrap(err)
	}

	e.UpdatedAt = dataLabelJobMap.UpdatedAt.Time
	e.CreatedAt = dataLabelJobMap.CreatedAt.Time
	return nil
}

func (m *Manager) UpsertLabelServiceMap(ctx context.Context, serviceMap *eve.LabelServiceMap) error {
	dataLabelServiceMap := toDataLabelServiceMap(*serviceMap)
	err := m.repo.UpsertLabelServiceMap(ctx, &dataLabelServiceMap)
	if err != nil {
		return errors.Wrap(err)
	}

	serviceMap.UpdatedAt = dataLabelServiceMap.UpdatedAt.Time
	serviceMap.CreatedAt = dataLabelServiceMap.CreatedAt.Time
	return nil
}

func (m *Manager) DeleteLabel(ctx context.Context, id int) error {
	err := m.repo.DeleteLabel(ctx, id)
	if err != nil {
		return service.CheckForNotFoundError(err)
	}
	return nil
}

func (m Manager) DeleteLabelKey(ctx context.Context, id int, key string) (eve.Label, error) {
	label, err := m.repo.DeleteLabelKey(ctx, id, key)
	if err != nil {
		return eve.Label{}, service.CheckForNotFoundError(err)
	}

	return fromDataLabel(*label), nil
}

func (m *Manager) DeleteLabelJobMap(ctx context.Context, labelID int, description string) error {
	err := m.repo.DeleteLabelJobMap(ctx, labelID, description)
	if err != nil {
		return service.CheckForNotFoundError(err)
	}

	return nil
}

func (m *Manager) DeleteLabelServiceMap(ctx context.Context, labelID int, description string) error {
	err := m.repo.DeleteLabelServiceMap(ctx, labelID, description)
	if err != nil {
		return service.CheckForNotFoundError(err)
	}

	return nil
}

func (m *Manager) ServiceLabel(ctx context.Context, id int) (eve.MetadataField, error) {
	metadata, err := m.repo.ServiceLabels(ctx, id)
	if err != nil {
		return nil, service.CheckForNotFoundError(err)
	}

	var collectedMetadata []eve.MetadataField
	for _, x := range metadata {
		collectedMetadata = append(collectedMetadata, x.Data.AsMapOrEmpty())
	}

	return m.mergeMetadata(collectedMetadata), nil
}

func (m *Manager) JobLabel(ctx context.Context, id int) (eve.MetadataField, error) {
	metadata, err := m.repo.JobLabels(ctx, id)
	if err != nil {
		return nil, service.CheckForNotFoundError(err)
	}

	var collectedMetadata []eve.MetadataField
	for _, x := range metadata {
		collectedMetadata = append(collectedMetadata, x.Data.AsMapOrEmpty())
	}

	return m.mergeMetadata(collectedMetadata), nil
}

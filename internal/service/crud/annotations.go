package crud

import (
	"context"
	"strconv"

	"gitlab.unanet.io/devops/eve/internal/data"
	"gitlab.unanet.io/devops/eve/internal/service"
	"gitlab.unanet.io/devops/eve/pkg/eve"
	"gitlab.unanet.io/devops/go/pkg/errors"
	"gitlab.unanet.io/devops/go/pkg/json"
)

func toDataAnnotationServiceMap(m eve.AnnotationServiceMap) data.AnnotationServiceMap {
	dm := data.AnnotationServiceMap{
		Description:   m.Description,
		AnnotationID:  m.AnnotationID,
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

func toDataAnnotationJobMap(m eve.AnnotationJobMap) data.AnnotationJobMap {
	dm := data.AnnotationJobMap{
		Description:   m.Description,
		AnnotationID:  m.AnnotationID,
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

func toDataAnnotation(m eve.Annotation) data.Annotation {
	return data.Annotation{
		ID:          m.ID,
		Description: m.Description,
		Data:        json.FromMap(m.Data),
	}
}

func fromDataAnnotation(m data.Annotation) eve.Annotation {
	return eve.Annotation{
		ID:          m.ID,
		Description: m.Description,
		Data:        m.Data.AsMap(),
		CreatedAt:   m.CreatedAt.Time,
		UpdatedAt:   m.UpdatedAt.Time,
	}
}

func fromDataAnnotationList(annotations []data.Annotation) []eve.Annotation {
	var list []eve.Annotation
	for _, x := range annotations {
		list = append(list, fromDataAnnotation(x))
	}
	return list
}

func fromDataAnnotationServiceToAnnotation(m data.AnnotationService) eve.Annotation {
	return eve.Annotation{
		ID:          m.AnnotationID,
		Description: m.AnnotationDescription,
		Data:        m.Data.AsMap(),
		CreatedAt:   m.CreatedAt.Time,
		UpdatedAt:   m.UpdatedAt.Time,
	}
}

func fromDataAnnotationServiceListToAnnotationList(m []data.AnnotationService) []eve.Annotation {
	var list []eve.Annotation
	for _, x := range m {
		list = append(list, fromDataAnnotationServiceToAnnotation(x))
	}
	return list
}

func fromDataAnnotationServiceMaps(m []data.AnnotationServiceMap) []eve.AnnotationServiceMap {
	var list []eve.AnnotationServiceMap
	for _, x := range m {
		list = append(list, fromDataAnnotationServiceMap(x))
	}
	return list
}

func fromDataAnnotationJobMap(m data.AnnotationJobMap) eve.AnnotationJobMap {
	return eve.AnnotationJobMap{
		Description:   m.Description,
		AnnotationID:  m.AnnotationID,
		EnvironmentID: int(m.EnvironmentID.Int32),
		ArtifactID:    int(m.ArtifactID.Int32),
		NamespaceID:   int(m.NamespaceID.Int32),
		JobID:         int(m.JobID.Int32),
		StackingOrder: m.StackingOrder,
		CreatedAt:     m.CreatedAt.Time,
		UpdatedAt:     m.UpdatedAt.Time,
	}
}

func fromDataAnnotationServiceMap(m data.AnnotationServiceMap) eve.AnnotationServiceMap {
	return eve.AnnotationServiceMap{
		Description:   m.Description,
		AnnotationID:  m.AnnotationID,
		EnvironmentID: int(m.EnvironmentID.Int32),
		ArtifactID:    int(m.ArtifactID.Int32),
		NamespaceID:   int(m.NamespaceID.Int32),
		ServiceID:     int(m.ServiceID.Int32),
		StackingOrder: m.StackingOrder,
		CreatedAt:     m.CreatedAt.Time,
		UpdatedAt:     m.UpdatedAt.Time,
	}
}

func fromDataAnnotationService(m data.AnnotationService) eve.AnnotationServiceMap {
	return eve.AnnotationServiceMap{
		Description:   m.MapDescription,
		AnnotationID:  m.AnnotationID,
		EnvironmentID: int(m.MapEnvironmentID.Int32),
		ArtifactID:    int(m.MapArtifactID.Int32),
		NamespaceID:   int(m.MapNamespaceID.Int32),
		ServiceID:     int(m.MapServiceID.Int32),
		StackingOrder: m.StackingOrder,
		CreatedAt:     m.CreatedAt.Time,
		UpdatedAt:     m.UpdatedAt.Time,
	}
}

func fromDataAnnotationServiceList(m []data.AnnotationService) []eve.AnnotationServiceMap {
	var list []eve.AnnotationServiceMap
	for _, x := range m {
		list = append(list, fromDataAnnotationService(x))
	}
	return list
}

func fromDataAnnotationJobMaps(maps []data.AnnotationJobMap) []eve.AnnotationJobMap {
	var list []eve.AnnotationJobMap
	for _, x := range maps {
		list = append(list, fromDataAnnotationJobMap(x))
	}
	return list
}

func (m Manager) Annotations(ctx context.Context) ([]eve.Annotation, error) {
	annotations, err := m.repo.Annotations(ctx)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	return fromDataAnnotationList(annotations), nil
}

func (m Manager) CreateAnnotation(ctx context.Context, annotation *eve.Annotation) error {
	dataAnnotation := toDataAnnotation(*annotation)
	err := m.repo.UpsertAnnotation(ctx, &dataAnnotation)
	if err != nil {
		return errors.Wrap(err)
	}

	annotation.UpdatedAt = dataAnnotation.UpdatedAt.Time
	annotation.CreatedAt = dataAnnotation.CreatedAt.Time
	annotation.ID = dataAnnotation.ID
	return nil
}

func (m Manager) UpsertMergeAnnotation(ctx context.Context, annotation *eve.Annotation) error {
	dataAnnotation := toDataAnnotation(*annotation)
	err := m.repo.UpsertMergeAnnotation(ctx, &dataAnnotation)
	if err != nil {
		return errors.Wrap(err)
	}

	annotation.UpdatedAt = dataAnnotation.UpdatedAt.Time
	annotation.CreatedAt = dataAnnotation.CreatedAt.Time
	annotation.ID = dataAnnotation.ID
	annotation.Data = dataAnnotation.Data.AsMap()
	return nil
}

func (m *Manager) GetAnnotation(ctx context.Context, id string) (*eve.Annotation, error) {
	var annotation *data.Annotation
	if intID, err := strconv.Atoi(id); err == nil {
		annotation, err = m.repo.GetAnnotation(ctx, intID)
		if err != nil {
			return nil, service.CheckForNotFoundError(err)
		}
	} else {
		annotation, err = m.repo.GetAnnotationByDescription(ctx, id)
		if err != nil {
			return nil, service.CheckForNotFoundError(err)
		}
	}

	r := fromDataAnnotation(*annotation)
	return &r, nil
}

func (m *Manager) JobAnnotationMaps(ctx context.Context, id int) ([]eve.AnnotationJobMap, error) {
	maps, err := m.repo.JobAnnotationMaps(ctx, id)
	if err != nil {
		return nil, service.CheckForNotFoundError(err)
	}

	return fromDataAnnotationJobMaps(maps), nil
}

func (m *Manager) ServiceAnnotationMaps(ctx context.Context, id int) ([]eve.AnnotationServiceMap, error) {
	maps, err := m.repo.ServiceAnnotations(ctx, id)
	if err != nil {
		return nil, service.CheckForNotFoundError(err)
	}

	return fromDataAnnotationServiceList(maps), nil
}

func (m *Manager) ServiceAnnotationMapsByAnnotationID(ctx context.Context, id int) ([]eve.AnnotationServiceMap, error) {
	maps, err := m.repo.ServiceAnnotationMapsByAnnotationID(ctx, id)
	if err != nil {
		return nil, service.CheckForNotFoundError(err)
	}

	return fromDataAnnotationServiceMaps(maps), nil
}

func (m *Manager) JobAnnotationMapsByAnnotationID(ctx context.Context, id int) ([]eve.AnnotationJobMap, error) {
	maps, err := m.repo.JobAnnotationMapsByAnnotationID(ctx, id)
	if err != nil {
		return nil, service.CheckForNotFoundError(err)
	}

	return fromDataAnnotationJobMaps(maps), nil
}

func (m *Manager) UpsertAnnotationJobMap(ctx context.Context, e *eve.AnnotationJobMap) error {
	dataAnnotationJobMap := toDataAnnotationJobMap(*e)
	err := m.repo.UpsertAnnotationJobMap(ctx, &dataAnnotationJobMap)
	if err != nil {
		return errors.Wrap(err)
	}

	e.UpdatedAt = dataAnnotationJobMap.UpdatedAt.Time
	e.CreatedAt = dataAnnotationJobMap.CreatedAt.Time
	return nil
}

func (m *Manager) UpsertAnnotationServiceMap(ctx context.Context, serviceMap *eve.AnnotationServiceMap) error {
	dataAnnotationServiceMap := toDataAnnotationServiceMap(*serviceMap)
	err := m.repo.UpsertAnnotationServiceMap(ctx, &dataAnnotationServiceMap)
	if err != nil {
		return errors.Wrap(err)
	}

	serviceMap.UpdatedAt = dataAnnotationServiceMap.UpdatedAt.Time
	serviceMap.CreatedAt = dataAnnotationServiceMap.CreatedAt.Time
	return nil
}

func (m *Manager) DeleteAnnotation(ctx context.Context, id int) error {
	err := m.repo.DeleteAnnotation(ctx, id)
	if err != nil {
		return service.CheckForNotFoundError(err)
	}
	return nil
}

func (m Manager) DeleteAnnotationKey(ctx context.Context, id int, key string) (eve.Annotation, error) {
	annotation, err := m.repo.DeleteAnnotationKey(ctx, id, key)
	if err != nil {
		return eve.Annotation{}, service.CheckForNotFoundError(err)
	}

	return fromDataAnnotation(*annotation), nil
}

func (m *Manager) DeleteAnnotationJobMap(ctx context.Context, annotationID int, description string) error {
	err := m.repo.DeleteAnnotationJobMap(ctx, annotationID, description)
	if err != nil {
		return service.CheckForNotFoundError(err)
	}

	return nil
}

func (m *Manager) DeleteAnnotationServiceMap(ctx context.Context, annotationID int, description string) error {
	err := m.repo.DeleteAnnotationServiceMap(ctx, annotationID, description)
	if err != nil {
		return service.CheckForNotFoundError(err)
	}

	return nil
}

func (m *Manager) ServiceAnnotation(ctx context.Context, id int) (eve.MetadataField, error) {
	metadata, err := m.repo.ServiceAnnotations(ctx, id)
	if err != nil {
		return nil, service.CheckForNotFoundError(err)
	}

	var collectedMetadata []eve.MetadataField
	for _, x := range metadata {
		collectedMetadata = append(collectedMetadata, x.Data.AsMap())
	}

	return m.mergeMetadata(collectedMetadata), nil
}

func (m *Manager) JobAnnotation(ctx context.Context, id int) (eve.MetadataField, error) {
	metadata, err := m.repo.JobAnnotations(ctx, id)
	if err != nil {
		return nil, service.CheckForNotFoundError(err)
	}

	var collectedMetadata []eve.MetadataField
	for _, x := range metadata {
		collectedMetadata = append(collectedMetadata, x.Data.AsMap())
	}

	return m.mergeMetadata(collectedMetadata), nil
}

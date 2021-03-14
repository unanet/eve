package crud

import (
	"context"
	gojson "encoding/json"
	goerrors "github.com/pkg/errors"
	"gitlab.unanet.io/devops/eve/internal/data"
	"gitlab.unanet.io/devops/eve/internal/service"
	"gitlab.unanet.io/devops/eve/pkg/eve"
	"gitlab.unanet.io/devops/go/pkg/errors"
	"gitlab.unanet.io/devops/go/pkg/jmerge"
	"gitlab.unanet.io/devops/go/pkg/json"
	"strconv"
)

func toDataDefinitionServiceMap(m eve.DefinitionServiceMap) data.DefinitionServiceMap {
	dm := data.DefinitionServiceMap{
		Description:   m.Description,
		DefinitionID:  m.DefinitionID,
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

	if m.ClusterID != 0 {
		dm.ClusterID.Int32 = int32(m.ClusterID)
		dm.ClusterID.Valid = true
	}

	if m.ServiceID != 0 {
		dm.ServiceID.Int32 = int32(m.ServiceID)
		dm.ServiceID.Valid = true
	}

	return dm
}

func toDataDefinitionJobMap(m eve.DefinitionJobMap) data.DefinitionJobMap {
	dm := data.DefinitionJobMap{
		Description:   m.Description,
		DefinitionID:  m.DefinitionID,
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

	if m.ClusterID != 0 {
		dm.ClusterID.Int32 = int32(m.ClusterID)
		dm.ClusterID.Valid = true
	}

	if m.JobID != 0 {
		dm.JobID.Int32 = int32(m.JobID)
		dm.JobID.Valid = true
	}

	return dm
}

func toDataDefinition(m eve.Definition) data.Definition {
	return data.Definition{
		ID:               m.ID,
		Description:      m.Description,
		DefinitionTypeID: m.DefinitionTypeID,
		Data:             json.FromMapOrEmpty(m.Data),
	}
}

func fromDataDefinition(m data.Definition) eve.Definition {
	return eve.Definition{
		ID:               m.ID,
		Description:      m.Description,
		DefinitionTypeID: m.DefinitionTypeID,
		Data:             m.Data.AsMapOrEmpty(),
		CreatedAt:        m.CreatedAt.Time,
		UpdatedAt:        m.UpdatedAt.Time,
	}
}

func fromDataDefinitionList(defs []data.Definition) []eve.Definition {
	var list []eve.Definition
	for _, x := range defs {
		list = append(list, fromDataDefinition(x))
	}
	return list
}

func fromDataDefinitionServiceToDefinition(m data.DefinitionService) eve.Definition {
	return eve.Definition{
		ID:          m.DefinitionID,
		Description: m.DefinitionDescription,
		Data:        m.Data.AsMapOrEmpty(),
		CreatedAt:   m.CreatedAt.Time,
		UpdatedAt:   m.UpdatedAt.Time,
	}
}

func fromDataDefinitionServiceListToDefinitionList(m []data.DefinitionService) []eve.Definition {
	var list []eve.Definition
	for _, x := range m {
		list = append(list, fromDataDefinitionServiceToDefinition(x))
	}
	return list
}

func fromDataDefinitionServiceMaps(m []data.DefinitionServiceMap) []eve.DefinitionServiceMap {
	var list []eve.DefinitionServiceMap
	for _, x := range m {
		list = append(list, fromDataDefinitionServiceMap(x))
	}
	return list
}

func fromDataDefinitionJobMap(m data.DefinitionJobMap) eve.DefinitionJobMap {
	return eve.DefinitionJobMap{
		Description:   m.Description,
		DefinitionID:  m.DefinitionID,
		EnvironmentID: int(m.EnvironmentID.Int32),
		ArtifactID:    int(m.ArtifactID.Int32),
		NamespaceID:   int(m.NamespaceID.Int32),
		JobID:         int(m.JobID.Int32),
		ClusterID:     int(m.ClusterID.Int32),
		StackingOrder: m.StackingOrder,
		CreatedAt:     m.CreatedAt.Time,
		UpdatedAt:     m.UpdatedAt.Time,
	}
}

func fromDataDefinitionJobMaps(m []data.DefinitionJobMap) []eve.DefinitionJobMap {
	var list []eve.DefinitionJobMap
	for _, x := range m {
		list = append(list, fromDataDefinitionJobMap(x))
	}
	return list
}

func fromDataDefinitionServiceMap(m data.DefinitionServiceMap) eve.DefinitionServiceMap {
	return eve.DefinitionServiceMap{
		Description:   m.Description,
		DefinitionID:  m.DefinitionID,
		EnvironmentID: int(m.EnvironmentID.Int32),
		ArtifactID:    int(m.ArtifactID.Int32),
		NamespaceID:   int(m.NamespaceID.Int32),
		ServiceID:     int(m.ServiceID.Int32),
		ClusterID:     int(m.ClusterID.Int32),
		StackingOrder: m.StackingOrder,
		CreatedAt:     m.CreatedAt.Time,
		UpdatedAt:     m.UpdatedAt.Time,
	}
}

func fromDataDefinitionService(m data.DefinitionService) eve.DefinitionServiceMap {
	return eve.DefinitionServiceMap{
		Description:   m.MapDescription,
		DefinitionID:  m.DefinitionID,
		EnvironmentID: int(m.MapEnvironmentID.Int32),
		ArtifactID:    int(m.MapArtifactID.Int32),
		NamespaceID:   int(m.MapNamespaceID.Int32),
		ServiceID:     int(m.MapServiceID.Int32),
		ClusterID:     int(m.MapClusterID.Int32),
		StackingOrder: m.StackingOrder,
		CreatedAt:     m.CreatedAt.Time,
		UpdatedAt:     m.UpdatedAt.Time,
	}
}

func fromDataDefinitionServiceList(m []data.DefinitionService) []eve.DefinitionServiceMap {
	var list []eve.DefinitionServiceMap
	for _, x := range m {
		list = append(list, fromDataDefinitionService(x))
	}
	return list
}

func (m Manager) Definitions(ctx context.Context) ([]eve.Definition, error) {
	defs, err := m.repo.Definition(ctx)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	return fromDataDefinitionList(defs), nil
}

func (m Manager) UpsertMergeDefinition(ctx context.Context, def *eve.Definition) error {
	dataDefinition := toDataDefinition(*def)
	err := m.repo.UpsertMergeDefinition(ctx, &dataDefinition)
	if err != nil {
		return errors.Wrap(err)
	}

	def.UpdatedAt = dataDefinition.UpdatedAt.Time
	def.CreatedAt = dataDefinition.CreatedAt.Time
	def.ID = dataDefinition.ID
	def.Data = dataDefinition.Data.AsMapOrEmpty()
	return nil
}

func (m Manager) CreateDefinition(ctx context.Context, def *eve.Definition) error {
	dataDefinition := toDataDefinition(*def)
	err := m.repo.UpsertDefinition(ctx, &dataDefinition)
	if err != nil {
		return errors.Wrap(err)
	}

	def.UpdatedAt = dataDefinition.UpdatedAt.Time
	def.CreatedAt = dataDefinition.CreatedAt.Time
	def.ID = dataDefinition.ID
	return nil
}

func (m Manager) DeleteDefinitionKey(ctx context.Context, id int, key string) (eve.Definition, error) {
	definition, err := m.repo.DeleteDefinitionKey(ctx, id, key)
	if err != nil {
		return eve.Definition{}, service.CheckForNotFoundError(err)
	}

	return fromDataDefinition(*definition), nil
}

func (m Manager) DeleteDefinition(ctx context.Context, id int) error {
	err := m.repo.DeleteDefinition(ctx, id)
	if err != nil {
		return service.CheckForNotFoundError(err)
	}
	return nil
}

func (m Manager) GetDefinition(ctx context.Context, id string) (*eve.Definition, error) {
	var definition *data.Definition
	if intID, err := strconv.Atoi(id); err == nil {
		definition, err = m.repo.GetDefinition(ctx, intID)
		if err != nil {
			return nil, service.CheckForNotFoundError(err)
		}
	} else {
		definition, err = m.repo.GetDefinitionByDescription(ctx, id)
		if err != nil {
			return nil, service.CheckForNotFoundError(err)
		}
	}

	r := fromDataDefinition(*definition)
	return &r, nil
}

func (m Manager) UpsertDefinitionServiceMap(ctx context.Context, serviceMap *eve.DefinitionServiceMap) error {
	dataDefinitionServiceMap := toDataDefinitionServiceMap(*serviceMap)
	err := m.repo.UpsertDefinitionServiceMap(ctx, &dataDefinitionServiceMap)
	if err != nil {
		return errors.Wrap(err)
	}

	serviceMap.UpdatedAt = dataDefinitionServiceMap.UpdatedAt.Time
	serviceMap.CreatedAt = dataDefinitionServiceMap.CreatedAt.Time
	return nil
}

func (m Manager) DeleteDefinitionServiceMap(ctx context.Context, definitionID int, description string) error {
	err := m.repo.DeleteDefinitionServiceMap(ctx, definitionID, description)
	if err != nil {
		return service.CheckForNotFoundError(err)
	}

	return nil
}

func (m Manager) ServiceDefinitionMapsByDefinitionID(ctx context.Context, id int) ([]eve.DefinitionServiceMap, error) {
	maps, err := m.repo.ServiceDefinitionMapsByDefinitionID(ctx, id)
	if err != nil {
		return nil, service.CheckForNotFoundError(err)
	}

	return fromDataDefinitionServiceMaps(maps), nil
}

func (m Manager) UpsertDefinitionJobMap(ctx context.Context, e *eve.DefinitionJobMap) error {
	dataDefinitionJobMap := toDataDefinitionJobMap(*e)
	err := m.repo.UpsertDefinitionJobMap(ctx, &dataDefinitionJobMap)
	if err != nil {
		return errors.Wrap(err)
	}

	e.UpdatedAt = dataDefinitionJobMap.UpdatedAt.Time
	e.CreatedAt = dataDefinitionJobMap.CreatedAt.Time
	return nil
}

func (m Manager) DeleteDefinitionJobMap(ctx context.Context, definitionID int, description string) error {
	err := m.repo.DeleteDefinitionJobMap(ctx, definitionID, description)
	if err != nil {
		return service.CheckForNotFoundError(err)
	}

	return nil
}

func (m Manager) JobDefinitionMapsByDefinitionID(ctx context.Context, id int) ([]eve.DefinitionJobMap, error) {
	maps, err := m.repo.JobDefinitionMapsByDefinitionID(ctx, id)
	if err != nil {
		return nil, service.CheckForNotFoundError(err)
	}

	return fromDataDefinitionJobMaps(maps), nil
}

func (m *Manager) JobDefinitionData(ctx context.Context, id int) (eve.MetadataField, error) {
	definitionData, err := m.repo.JobDefinition(ctx, id)
	if err != nil {
		return nil, service.CheckForNotFoundError(err)
	}

	var definitionSpecs []eve.DefinitionSpec
	for _, x := range definitionData {
		var defSpec = make(eve.DefinitionSpec)
		var d = make(map[string]interface{})
		if err := gojson.Unmarshal(x.Data, &d); err != nil {
			return nil, errors.Wrapf("failed to parse the deployment definition: %s", err)
		}
		defSpec[x.DefinitionType] = d
		definitionSpecs = append(definitionSpecs, defSpec)
	}

	return m.mergeDefinitionData(definitionSpecs)
}

func (m *Manager) ServiceDefinitionData(ctx context.Context, id int) (eve.MetadataField, error) {
	definitionData, err := m.repo.ServiceDefinition(ctx, id)
	if err != nil {
		return nil, service.CheckForNotFoundError(err)
	}

	var definitionSpecs []eve.DefinitionSpec
	for _, x := range definitionData {
		var defSpec = make(eve.DefinitionSpec)
		var defSpecData = make(map[string]interface{})
		if err := gojson.Unmarshal(x.Data, &defSpecData); err != nil {
			return nil, errors.Wrapf("failed to parse the deployment definition: %s", err)
		}
		defSpec[x.DefinitionType] = defSpecData
		definitionSpecs = append(definitionSpecs, defSpec)
	}

	return m.mergeDefinitionData(definitionSpecs)
}

// eve.DefinitionSpec: map[string]map[string]interface{}
// Def Specs: []eve.DefinitionSpec
// Def Spec: eve.DefinitionSpec
// Def Data: map[string]interface{}
//
// Example: defSpec1["appsv1.Deployment"] = '{"spec": {"template": {"spec": {"nodeSelector": {"node-group": "shared"}}}}}'
// Example: defSpec2["appsv1.Deployment"] = '{"spec": {"template": {"spec": {"containers": [{"livenessProbe": {"httpGet": {"path": "/analytics-api/Api.asmx", "port": 8080}, "periodSeconds": 10, "initialDelaySeconds": 30}}]}}}}'
// Example: defSpec3["appsv1.Deployment"] = '{"spec": {"template": {"spec": {"containers": [{"readinessProbe": {"httpGet": {"path": "/analytics-api/Api.asmx", "port": 8080 }, "periodSeconds": 10, "initialDelaySeconds": 45}}]}}}}'
// Example: defSpec4["v2beta2.HorizontalPodAutoscaler"] = '{"spec": {"minReplicas":2, "maxReplicas": 10}'
//
// Example: defSpecs = []eve.DefinitionSpec{defSpec1,defSpec2,defSpec3,defSpec4}
type result struct {
	defs map[string]interface{}
}

func (m Manager) mergeDefinitionData(defSpecs []eve.DefinitionSpec) (eve.MetadataField, error) {

	r := result{defs: make(map[string]interface{})}

	for _, defSpec := range defSpecs {
		resultSpecData := make(map[string]interface{})
		for specType, specData := range defSpec {
			existingSpecData, ok := r.defs[specType]
			if !ok {
				resultSpecData = jmerge.Merge(resultSpecData, specData)
			} else {
				datamap, ok := existingSpecData.(map[string]interface{})
				if !ok {
					return nil, goerrors.New("failed to cast existing spec data back to map interface")
				}
				resultSpecData = jmerge.Merge(datamap, specData)
			}
			r.defs[specType] = resultSpecData
		}
	}

	return r.defs, nil
}

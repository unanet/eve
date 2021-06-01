package crud

import (
	"context"
	"gitlab.unanet.io/devops/eve/internal/data"
	"gitlab.unanet.io/devops/eve/pkg/eve"
)


func (m *Manager) DefinitionTypes(ctx context.Context) (models []eve.DefinitionType, err error)  {
	dbArtifacts, err := m.repo.DefinitionTypes(ctx)
	if err != nil {
		return nil, err
	}

	return fromDataDefinitionTypeList(dbArtifacts), err
}

func (m *Manager) CreateDefinitionType(ctx context.Context, model *eve.DefinitionType) error  {
	dbModel := toDataDefinitionType(*model)
	if err := m.repo.CreateDefinitionType(ctx, &dbModel); err != nil {
		return err
	}

	model.ID = dbModel.ID
	return nil
}

func (m *Manager) UpdateDefinitionType(ctx context.Context, model *eve.DefinitionType) (err error)  {

	dbModel := toDataDefinitionType(*model)
	if err := m.repo.UpdateDefinitionType(ctx, &dbModel); err != nil {
		return err
	}

	model.CreatedAt = dbModel.CreatedAt.Time

	return nil
}

func (m *Manager) DeleteDefinitionType(ctx context.Context, id int) (err error)  {
	return m.repo.DeleteDefinitionType(ctx, id)
}

func fromDataDefinitionTypeList(artifacts []data.DefinitionType) []eve.DefinitionType {
	var list []eve.DefinitionType
	for _, x := range artifacts {
		list = append(list, fromDataDefinitionTypeToDefinitionType(x))
	}
	return list
}

func fromDataDefinitionTypeToDefinitionType(dbM data.DefinitionType) eve.DefinitionType {
	return eve.DefinitionType{
		ID:              dbM.ID,
		Name:            dbM.Name,
		Description:     dbM.Description,
		Class:           dbM.Class,
		Version:         dbM.Version,
		Kind:            dbM.Kind,
		DefinitionOrder: dbM.DefinitionOrder,
		CreatedAt:       dbM.CreatedAt.Time,
		UpdatedAt:       dbM.UpdatedAt.Time,
	}
}

func toDataDefinitionType(dbM eve.DefinitionType) data.DefinitionType {
	return data.DefinitionType{
		ID:              dbM.ID,
		Name:            dbM.Name,
		Description:     dbM.Description,
		Class:           dbM.Class,
		Version:         dbM.Version,
		Kind:            dbM.Kind,
		DefinitionOrder: dbM.DefinitionOrder,
	}
}

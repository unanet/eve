package crud

import (
	"context"
	"gitlab.unanet.io/devops/eve/internal/data"
	"gitlab.unanet.io/devops/eve/pkg/eve"
	"gitlab.unanet.io/devops/go/pkg/errors"
	"gitlab.unanet.io/devops/go/pkg/json"
)


func (m *Manager) MetadataHistory(ctx context.Context) (models []eve.MetadataHistory, err error)  {
	dbResults, err := m.repo.MetadataHistory(ctx)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	return fromDataMetadataHistoryList(dbResults), nil
}

func fromDataMetadataHistory(dbModel data.MetadataHistory) eve.MetadataHistory {
	deletedTime := &dbModel.Deleted.Time
	if !dbModel.Deleted.Valid {
		deletedTime = nil
	}

	return eve.MetadataHistory{
		MetadataId: dbModel.MetadataId,
		Description: dbModel.Description,
		Value: dbModel.Value.AsMapOrEmpty(),
		Created: dbModel.Created.Time,
		CreatedBy: dbModel.CreatedBy,
		Deleted: deletedTime,
		DeletedBy: dbModel.DeletedBy,
	}
}

func fromDataMetadataHistoryList(dbModels []data.MetadataHistory) []eve.MetadataHistory {
	var list []eve.MetadataHistory
	for _, x := range dbModels {
		list = append(list, fromDataMetadataHistory(x))
	}
	return list
}

func toDataMetadataHistory(model eve.MetadataHistory) data.MetadataHistory {
	return data.MetadataHistory{
		MetadataId: model.MetadataId,
		Description: model.Description,
		Value: json.FromMapOrEmpty(model.Value),
		CreatedBy: model.CreatedBy,
		DeletedBy: model.DeletedBy,
	}
}


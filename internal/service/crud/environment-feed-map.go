package crud

import (
	"context"
	"github.com/unanet/go/pkg/errors"

	"github.com/unanet/eve/internal/data"
	"github.com/unanet/eve/pkg/eve"
)

func (m *Manager) EnvironmentFeedMaps(ctx context.Context) (models []eve.EnvironmentFeedMap, err error) {

	dbModels, err := m.repo.EnvironmentFeedMaps(ctx)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	return fromDataEnvironmentFeedMapList(dbModels), err
}

func (m *Manager) CreateEnvironmentFeedMap(ctx context.Context, model *eve.EnvironmentFeedMap) error {
	dbModel := toDataEnvironmentFeedMap(*model)
	if err := m.repo.CreateEnvironmentFeedMap(ctx, &dbModel); err != nil {
		return errors.Wrap(err)
	}

	return nil
}

func (m *Manager) UpdateEnvironmentFeedMap(ctx context.Context, model *eve.EnvironmentFeedMap) (err error) {

	dbModel := toDataEnvironmentFeedMap(*model)
	if err := m.repo.UpdateEnvironmentFeedMap(ctx, &dbModel); err != nil {
		return errors.Wrap(err)
	}

	return nil
}

func (m *Manager) DeleteEnvironmentFeedMap(ctx context.Context, model *eve.EnvironmentFeedMap) (err error) {
	return m.repo.DeleteEnvironmentFeedMap(ctx, model.EnvironmentID, model.FeedID)
}

func fromDataEnvironmentFeedMapList(feedMaps []data.EnvironmentFeedMap) []eve.EnvironmentFeedMap {
	var list []eve.EnvironmentFeedMap
	for _, x := range feedMaps {
		list = append(list, fromEnvironmentFeedMap(x))
	}
	return list
}

func fromEnvironmentFeedMap(dbModel data.EnvironmentFeedMap) eve.EnvironmentFeedMap {
	return eve.EnvironmentFeedMap{
		FeedID:        dbModel.FeedID,
		EnvironmentID: dbModel.EnvironmentID,
	}
}

func toDataEnvironmentFeedMap(m eve.EnvironmentFeedMap) data.EnvironmentFeedMap {
	return data.EnvironmentFeedMap{
		FeedID:        m.FeedID,
		EnvironmentID: m.EnvironmentID,
	}
}

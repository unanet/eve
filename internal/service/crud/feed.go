package crud

import (
	"context"
	"github.com/unanet/eve/internal/data"
	"github.com/unanet/eve/pkg/eve"
	"github.com/unanet/go/pkg/errors"
)

func (m *Manager) Feeds(ctx context.Context) (models []eve.Feed, err error) {
	dbModels, err := m.repo.Feeds(ctx)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	return fromDataFeedList(dbModels), nil
}

func (m *Manager) CreateFeed(ctx context.Context, model *eve.Feed) error {
	dbModel := toDataFeed(*model)
	if err := m.repo.CreateFeed(ctx, &dbModel); err != nil {
		return errors.Wrap(err)
	}

	model.ID = dbModel.ID

	return nil
}

func (m *Manager) UpdateFeed(ctx context.Context, model *eve.Feed) (err error) {
	dbModel := toDataFeed(*model)
	if err := m.repo.UpdateFeed(ctx, &dbModel); err != nil {
		return err
	}

	return nil
}
func (m *Manager) DeleteFeed(ctx context.Context, id int) (err error) {
	return m.repo.DeleteFeed(ctx, id)
}

func fromDataFeed(dbModel data.Feed) eve.Feed {
	return eve.Feed{
		ID:             dbModel.ID,
		Name:           dbModel.Name,
		PromotionOrder: dbModel.PromotionOrder,
		FeedType:       dbModel.FeedType,
		Alias:          dbModel.Alias,
	}
}

func fromDataFeedList(dbModels data.Feeds) []eve.Feed {
	var list []eve.Feed
	for _, x := range dbModels {
		list = append(list, fromDataFeed(x))
	}
	return list
}

func toDataFeed(model eve.Feed) data.Feed {
	return data.Feed{
		ID:             model.ID,
		Name:           model.Name,
		PromotionOrder: model.PromotionOrder,
		FeedType:       model.FeedType,
		Alias:          model.Alias,
	}
}

package data

import (
	"context"
	"database/sql"

	"gitlab.unanet.io/devops/eve/pkg/errors"
)

type Feed struct {
	ID             string `db:"id"`
	Name           string `db:"name"`
	PromotionOrder int    `db:"promotion_order"`
	FeedType       string `db:"feed_type"`
}

type Feeds []Feed

func (r *Repo) FeedByName(ctx context.Context, name string) (*Feed, error) {
	var feed Feed

	row := r.db.QueryRowxContext(ctx, "select * from feed where name = $1", name)
	err := row.StructScan(&feed)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NotFoundErrorf("feed with name: %s, not found", name)
		}
		return nil, errors.Wrap(err)
	}

	return &feed, nil
}

func (r *Repo) FeedByID(ctx context.Context, id int) (*Feed, error) {
	var feed Feed

	row := r.db.QueryRowxContext(ctx, "select * from feed where id = $1", id)
	err := row.StructScan(&feed)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NotFoundErrorf("feed with id: %v, not found", id)
		}
		return nil, errors.Wrap(err)
	}

	return &feed, nil
}

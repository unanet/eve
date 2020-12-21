package data

import (
	"context"
	"database/sql"

	"gitlab.unanet.io/devops/go/pkg/errors"
)

type Feed struct {
	ID             int    `db:"id"`
	Name           string `db:"name"`
	PromotionOrder int    `db:"promotion_order"`
	FeedType       string `db:"feed_type"`
	Alias          string `db:"alias"`
}

type Feeds []Feed

func (r *Repo) FeedByAliasAndType(ctx context.Context, alias, feedType string) (*Feed, error) {
	var feed Feed

	// QA and INT are currently shared
	// if a user wants QA they really mean Int
	if alias == "qa" {
		alias = "int"
	}

	row := r.db.QueryRowxContext(ctx, "select * from feed where feed_type = $1 AND alias = $2", feedType, alias)
	err := row.StructScan(&feed)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NotFoundErrorf("feed with feed_type: %v and alias: %v, not found", feedType, alias)
		}
		return nil, errors.Wrap(err)
	}

	return &feed, nil
}

func (r *Repo) NextFeedByPromotionOrderType(ctx context.Context, promotionOrder int, feedType string) (*Feed, error) {
	var feed Feed

	row := r.db.QueryRowxContext(ctx, "select * from feed where feed_type = $1 AND promotion_order > $2 order by promotion_order asc limit 1;", feedType, promotionOrder)

	err := row.StructScan(&feed)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NotFoundErrorf("feed with feed_type: %v and promotion_order greater than: %v, not found", feedType, promotionOrder)
		}
		return nil, errors.Wrap(err)
	}

	return &feed, nil
}

func (r *Repo) PreviousFeedByPromotionOrderType(ctx context.Context, promotionOrder int, feedType string) (*Feed, error) {
	var feed Feed

	row := r.db.QueryRowxContext(ctx, "select * from feed where alias <> '' AND feed_type = $1 AND promotion_order < $2 order by promotion_order desc limit 1;", feedType, promotionOrder)

	err := row.StructScan(&feed)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NotFoundErrorf("feed with feed_type: %v and promotion_order less than: %v, not found", feedType, promotionOrder)
		}
		return nil, errors.Wrap(err)
	}

	return &feed, nil
}

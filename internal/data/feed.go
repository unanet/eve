package data

import (
	"context"
	"database/sql"
	"github.com/unanet/go/pkg/errors"
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

func (r *Repo) Feeds(ctx context.Context) ([]Feed, error) {
	rows, err := r.db.QueryxContext(ctx, `
		select 
			id,
			name,
			promotion_order,
			feed_type,
			alias
		from feed`)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	defer rows.Close()

	var cc []Feed
	for rows.Next() {
		if rows.Err() != nil {
			return nil, errors.Wrap(err)
		}

		var c Feed
		err = rows.StructScan(&c)
		if err != nil {
			return nil, errors.Wrap(err)
		}
		cc = append(cc, c)
	}

	return cc, nil
}

func (r *Repo) CreateFeed(ctx context.Context, model *Feed) error {
	err := r.db.QueryRowxContext(ctx, `
	INSERT INTO feed(id, name, promotion_order, feed_type, alias)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`, model.ID, model.Name, model.PromotionOrder, model.FeedType, model.Alias).
		StructScan(model)

	if err != nil {
		return errors.Wrap(err)
	}

	return nil
}

func (r *Repo) UpdateFeed(ctx context.Context, model *Feed) error {
	result, err := r.db.ExecContext(ctx, `
		update feed set 
			name = $2,
		    promotion_order = $3,
			feed_type = $4,
			alias = $5
		where id = $1
	`,
		model.ID,
		model.Name,
		model.PromotionOrder,
		model.FeedType,
		model.Alias)
	if err != nil {
		return errors.Wrap(err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err)
	}

	if affected == 0 {
		return errors.NotFoundf("feed id: %d not found", model.ID)
	}
	return nil
}

func (r *Repo) DeleteFeed(ctx context.Context, id int) error {
	return r.deleteByID(ctx, "feed", id)
}

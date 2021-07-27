package data

import (
	"context"
	"fmt"
	"github.com/unanet/go/pkg/errors"
)

type EnvironmentFeedMap struct {
	EnvironmentID int `db:"environment_id"`
	FeedID        int `db:"feed_id"`
}

func (r *Repo) EnvironmentFeedMaps(ctx context.Context) ([]EnvironmentFeedMap, error) {
	rows, err := r.db.QueryxContext(ctx, fmt.Sprintf(`
		select 
			environment_id,
			feed_id
		from environment_feed_map`))
	if err != nil {
		return nil, errors.Wrap(err)
	}
	defer rows.Close()

	var ss []EnvironmentFeedMap
	for rows.Next() {
		if rows.Err() != nil {
			return nil, errors.Wrap(err)
		}

		var s EnvironmentFeedMap
		err = rows.StructScan(&s)
		if err != nil {
			return nil, errors.Wrap(err)
		}
		ss = append(ss, s)
	}

	return ss, nil
}

func (r *Repo) CreateEnvironmentFeedMap(ctx context.Context, model *EnvironmentFeedMap) error {
	err := r.db.QueryRowxContext(ctx, `
	INSERT INTO environment_feed_map(environment_id, feed_id)
		VALUES ($1, $2)
	RETURNING environment_id, feed_id
	`,
		model.EnvironmentID,
		model.FeedID).
		StructScan(model)

	if err != nil {
		return errors.Wrap(err)
	}

	return nil
}

func (r *Repo) UpdateEnvironmentFeedMap(ctx context.Context, model *EnvironmentFeedMap) error {
	result, err := r.db.ExecContext(ctx, `
		update environment_feed_map set 
			environment_id = $1,
		    feed_id = $2
		where environment_id = $1
		and feed_id = $2
	`,
		model.EnvironmentID,
		model.FeedID,
	)
	if err != nil {
		return errors.Wrap(err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err)
	}

	if affected == 0 {
		return errors.NotFoundf("%s with environment id: %d and feed id: %d not found", "environment_feed_map", model.EnvironmentID, model.FeedID)
	}
	return nil
}

func (r *Repo) DeleteEnvironmentFeedMap(ctx context.Context, environmentID int, feedID int) error {
	whereQuery := fmt.Sprintf(`environment_id = %d AND feed_id = %d`,
		environmentID,
		feedID)

	return r.deleteWithQuery(ctx, "environment_feed_map", whereQuery)
}

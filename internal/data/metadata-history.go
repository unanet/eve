package data

import (
	"context"
	"database/sql"
	"gitlab.unanet.io/devops/go/pkg/errors"
	"gitlab.unanet.io/devops/go/pkg/json"
)

type MetadataHistory struct {
	MetadataId  int          `db:"metadata_id"`
	Description string       `db:"description"`
	Value       json.Object  `db:"value"`
	Created     sql.NullTime `db:"created"`
	CreatedBy   string       `db:"created_by"`
	Deleted     sql.NullTime `db:"deleted"`
	DeletedBy   *string      `db:"deleted_by"`
}

func (r *Repo) MetadataHistory(ctx context.Context) ([]MetadataHistory, error) {
	rows, err := r.db.QueryxContext(ctx, `
		SELECT 
			metadata_id,
			description,
			value,
			created,
			created_by,
			deleted,
			deleted_by
		FROM metadata_history 
		LIMIT 100`) // Arbitrary value to limit by

	if err != nil {
		return nil, errors.Wrap(err)
	}
	defer rows.Close()

	var mm []MetadataHistory
	for rows.Next() {
		if rows.Err() != nil {
			return nil, errors.Wrap(err)
		}

		var m MetadataHistory
		err = rows.StructScan(&m)
		if err != nil {
			return nil, errors.Wrap(err)
		}

		mm = append(mm, m)
	}

	return mm, nil
}

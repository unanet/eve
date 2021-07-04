package data

import (
	"context"
	"database/sql"
	"gitlab.unanet.io/devops/go/pkg/errors"
	"time"
)

type DefinitionType struct {
	ID              int          `db:"id"`
	Name            string       `db:"name"`
	Description     string       `db:"description"`
	CreatedAt       sql.NullTime `db:"created_at"`
	UpdatedAt       sql.NullTime `db:"updated_at"`
	Class           string       `db:"class"`
	Version         string       `db:"version"`
	Kind            string       `db:"kind"`
	DefinitionOrder string       `db:"definition_order"`
}

func (r *Repo) DefinitionTypes(ctx context.Context) ([]DefinitionType, error) {
	rows, err := r.db.QueryxContext(ctx, `
		select 
			id,
			name,
			description,
			created_at,
			updated_at,
			class,
			version,
			kind,
			definition_order
		from definition_type`)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	defer rows.Close()

	var ss []DefinitionType
	for rows.Next() {
		if rows.Err() != nil {
			return nil, errors.Wrap(err)
		}

		var s DefinitionType
		err = rows.StructScan(&s)
		if err != nil {
			return nil, errors.Wrap(err)
		}
		ss = append(ss, s)
	}

	return ss, nil
}

func (r *Repo) CreateDefinitionType(ctx context.Context, model *DefinitionType) error {
	model.CreatedAt.Time = time.Now().UTC()
	model.CreatedAt.Valid = true

	err := r.db.QueryRowxContext(ctx, `
	INSERT INTO definition_type(name, description, class, version, kind, definition_order, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	RETURNING id
	`,
		model.Name,
		model.Description,
		model.Class,
		model.Version,
		model.Kind,
		model.DefinitionOrder,
		model.CreatedAt).
		StructScan(model)

	if err != nil {
		return errors.Wrap(err)
	}

	return nil
}

func (r *Repo) UpdateDefinitionType(ctx context.Context, m *DefinitionType) error {
	m.UpdatedAt.Time = time.Now().UTC()
	m.UpdatedAt.Valid = true

	result, err := r.db.ExecContext(ctx, `
		update definition_type set 
			name = $2,
			description = $3, 
			class = $4,
			version = $5, 
			kind = $6,
			definition_order = $7,
			updated_at = $8
		where id = $1
		RETURNING created_at
	`,
		m.ID,
		m.Name,
		m.Description,
		m.Class,
		m.Version,
		m.Kind,
		m.DefinitionOrder,
		m.UpdatedAt,
	)
	if err != nil {
		return errors.Wrap(err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err)
	}

	if affected == 0 {
		return errors.NotFoundf("definition type by id: %d not found", m.ID)
	}
	return nil
}

func (r *Repo) DeleteDefinitionType(ctx context.Context, id int) error {
	return r.deleteByID(ctx, "definition_type", id)
}

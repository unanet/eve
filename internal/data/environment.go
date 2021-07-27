package data

import (
	"context"
	"database/sql"
	goErrors "errors"
	"time"

	"github.com/unanet/go/pkg/errors"
)

type Environment struct {
	ID          int          `db:"id"`
	Name        string       `db:"name"`
	Alias       string       `db:"alias"`
	Description string       `db:"description"`
	UpdatedAt   sql.NullTime `db:"updated_at"`
}

type Environments []Environment

func (r *Repo) EnvironmentByName(ctx context.Context, name string) (*Environment, error) {
	var environment Environment

	row := r.db.QueryRowxContext(ctx, `
		select id,
		       name,
		       alias,
		       description,
		       updated_at
		from environment where name = $1
		`, name)
	err := row.StructScan(&environment)
	if err != nil {
		if goErrors.Is(err, sql.ErrNoRows) {
			return nil, NotFoundErrorf("environment with name: %s, not found", name)
		}
		return nil, errors.Wrap(err)
	}

	return &environment, nil
}

func (r *Repo) EnvironmentByID(ctx context.Context, id int) (*Environment, error) {
	var environment Environment

	row := r.db.QueryRowxContext(ctx, `
		select id,
		       name,
		       alias,
		       description,
		       updated_at
		from environment where id = $1
		`, id)
	err := row.StructScan(&environment)
	if err != nil {
		if goErrors.Is(err, sql.ErrNoRows) {
			return nil, NotFoundErrorf("environment with id: %d, not found", id)
		}
		return nil, errors.Wrap(err)
	}

	return &environment, nil
}

func (r *Repo) Environments(ctx context.Context) (Environments, error) {
	rows, err := r.db.QueryxContext(ctx, `
		select id, 
		       name,
		       alias,
		       description 
		from environment order by name
		`)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	var environments []Environment
	for rows.Next() {
		var environment Environment
		err = rows.StructScan(&environment)
		if err != nil {
			return nil, errors.Wrap(err)
		}
		environments = append(environments, environment)
	}

	return environments, nil
}

func (r *Repo) UpdateEnvironment(ctx context.Context, environment *Environment) error {
	environment.UpdatedAt.Time = time.Now().UTC()
	environment.UpdatedAt.Valid = true
	result, err := r.db.ExecContext(ctx, `
		update environment set 
			description = $1,
			updated_at = $2
		where id = $3
	`,
		environment.Description,
		environment.UpdatedAt,
		environment.ID)
	if err != nil {
		return errors.Wrap(err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err)
	}

	if affected == 0 {
		return errors.NotFoundf("environment id: %d not found", environment.ID)
	}
	return nil
}

func (r *Repo) CreateEnvironment(ctx context.Context, model *Environment) error {
	err := r.db.QueryRowxContext(ctx, `
	INSERT INTO environment(id, name, alias, description)
		VALUES ($1, $2, $3, $4)
	`,
		model.ID,
		model.Name,
		model.Alias,
		model.Description).
		StructScan(model)

	if err != nil {
		return errors.Wrap(err)
	}

	return nil
}

func (r *Repo) DeleteEnvironment(ctx context.Context, id int) error {
	return r.deleteByID(ctx, "environment", id)
}

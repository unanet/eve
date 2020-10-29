package data

import (
	"context"
	"database/sql"
	goErrors "errors"
	"time"

	"gitlab.unanet.io/devops/eve/pkg/errors"
	"gitlab.unanet.io/devops/eve/pkg/json"
)

type Environment struct {
	ID          int          `db:"id"`
	Name        string       `db:"name"`
	Alias       string       `db:"alias"`
	Description string       `db:"description"`
	Metadata    json.Text    `db:"metadata"`
	UpdatedAt   sql.NullTime `db:"updated_at"`
}

type Environments []Environment

func (r *Repo) EnvironmentByName(ctx context.Context, name string) (*Environment, error) {
	var environment Environment

	row := r.db.QueryRowxContext(ctx, "select * from environment where name = $1", name)
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

	row := r.db.QueryRowxContext(ctx, "select * from environment where id = $1", id)
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
	rows, err := r.db.QueryxContext(ctx, "select id, name, description from environment order by name")
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
			metadata = $1,
			description = $2,
			updated_at = $3
		where id = $4
	`,
		environment.Metadata,
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

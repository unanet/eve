package data

import (
	"context"
	"database/sql"
	"time"

	uuid "github.com/satori/go.uuid"

	"gitlab.unanet.io/devops/eve/pkg/errors"
)

type Deployment struct {
	ID               uuid.UUID      `db:"id"`
	EnvironmentID    int            `db:"environment_id"`
	NamespaceID      int            `db:"namespace_id"`
	ReqID            string         `db:"req_id"`
	PlanOptions      JSONText       `db:"plan_options"`
	S3PlanLocation   sql.NullString `db:"s3_plan_location"`
	S3ResultLocation sql.NullString `db:"s3_result_location"`
	CreatedAt        sql.NullTime   `db:"created_at"`
	UpdatedAt        sql.NullTime   `db:"updated_at"`
}

func (r *Repo) CreateDeployment(ctx context.Context, d *Deployment) error {
	now := time.Now().UTC()
	d.CreatedAt = sql.NullTime{
		Time:  now,
		Valid: true,
	}
	d.UpdatedAt = sql.NullTime{
		Time:  now,
		Valid: true,
	}

	err := r.db.QueryRowxContext(ctx, `
	
	insert into deployment(environment_id, namespace_id, req_id, plan_options, s3_plan_location, s3_result_location, created_at, updated_at) 
		values ($1, $2, $3, $4, $5, $6, $7, $8)
		returning (id)
	
	`, d.EnvironmentID, d.NamespaceID, d.ReqID, d.PlanOptions, d.S3PlanLocation, d.S3ResultLocation, d.CreatedAt, d.UpdatedAt).
		Scan(&d.ID)

	if err != nil {
		return errors.Wrap(err)
	}

	return nil
}

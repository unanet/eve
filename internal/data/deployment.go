package data

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	uuid "github.com/satori/go.uuid"

	"gitlab.unanet.io/devops/eve/pkg/errors"
)

type DeploymentState string

const (
	DeploymentStateQueued    DeploymentState = "queued"
	DeploymentStateScheduled DeploymentState = "scheduled"
	DeploymentStateCompleted DeploymentState = "completed"
)

type Deployment struct {
	ID               uuid.UUID       `db:"id"`
	EnvironmentID    int             `db:"environment_id"`
	NamespaceID      int             `db:"namespace_id"`
	ReqID            string          `db:"req_id"`
	PlanOptions      JSONText        `db:"plan_options"`
	S3PlanLocation   sql.NullString  `db:"s3_plan_location"`
	S3ResultLocation sql.NullString  `db:"s3_result_location"`
	CreatedAt        sql.NullTime    `db:"created_at"`
	UpdatedAt        sql.NullTime    `db:"updated_at"`
	State            DeploymentState `db:"state"`
}

func (r *Repo) Deployment(ctx context.Context, id uuid.UUID) (*Deployment, error) {
	var deployment Deployment

	row := r.db.QueryRowxContext(ctx, "select * from deployment where id = $1", id)
	err := row.StructScan(&deployment)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, NotFoundErrorf("environment with id: %d, not found", id)
		}
		return nil, errors.Wrap(err)
	}

	return &deployment, nil
}

func (r *Repo) UpdateDeploymentMessageIDTx(ctx context.Context, tx driver.Tx, id uuid.UUID, messageID string) error {
	sTx, ok := tx.(*sqlx.Tx)
	if !ok {
		return fmt.Errorf("could not cast tx to sqlx.Tx")
	}

	_, err := sTx.ExecContext(ctx, "update deployment set message_id = $1 where id = $2", messageID, id)
	if err != nil {
		return errors.Wrap(err)
	}
	return nil
}

func (r *Repo) CreateDeploymentTx(ctx context.Context, d *Deployment) (driver.Tx, error) {
	tx, err := r.db.BeginTxx(ctx, &sql.TxOptions{
		Isolation: 0,
		ReadOnly:  false,
	})

	if err != nil {
		return nil, errors.Wrap(err)
	}

	now := time.Now().UTC()
	d.CreatedAt = sql.NullTime{
		Time:  now,
		Valid: true,
	}
	d.UpdatedAt = sql.NullTime{
		Time:  now,
		Valid: true,
	}

	err = tx.QueryRowxContext(ctx, `
	
	insert into deployment(environment_id, namespace_id, req_id, plan_options, s3_plan_location, s3_result_location, state, created_at, updated_at) 
		values ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		returning (id)
	
	`, d.EnvironmentID, d.NamespaceID, d.ReqID, d.PlanOptions, d.S3PlanLocation, d.S3ResultLocation, DeploymentStateQueued, d.CreatedAt, d.UpdatedAt).
		Scan(&d.ID)

	if err != nil {
		return nil, errors.Wrap(err)
	}

	return tx, nil
}

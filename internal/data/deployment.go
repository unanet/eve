package data

import (
	"context"
	"database/sql"
	"time"

	uuid "github.com/satori/go.uuid"

	"gitlab.unanet.io/devops/go/pkg/errors"
	"gitlab.unanet.io/devops/go/pkg/json"
)

type DeploymentState string

const (
	DeploymentStateQueued    DeploymentState = "queued"
	DeploymentStateScheduled DeploymentState = "scheduled"
	DeploymentStateCompleted DeploymentState = "completed"
)

type Deployment struct {
	ID            uuid.UUID       `db:"id"`
	EnvironmentID int             `db:"environment_id"`
	NamespaceID   int             `db:"namespace_id"`
	MessageID     sql.NullString  `db:"message_id"`
	ReceiptHandle sql.NullString  `db:"receipt_handle"`
	ReqID         string          `db:"req_id"`
	PlanOptions   json.Text       `db:"plan_options"`
	PlanLocation  json.Text       `db:"plan_location"`
	State         DeploymentState `db:"state"`
	User          string          `db:"user"`
	CreatedAt     sql.NullTime    `db:"created_at"`
	UpdatedAt     sql.NullTime    `db:"updated_at"`
}

func (r *Repo) UpdateDeploymentMessageID(ctx context.Context, id uuid.UUID, messageID string) error {
	result, err := r.db.ExecContext(ctx, "update deployment set message_id = $1, updated_at = $2 where id = $3", messageID, time.Now().UTC(), id)
	if err != nil {
		return errors.Wrap(err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err)
	}

	if affected == 0 {
		return errors.Wrapf("the following id: %s was not found to update in deployment table", id)
	}
	return nil
}

func (r *Repo) UpdateDeploymentPlanLocation(ctx context.Context, id uuid.UUID, location json.Text) error {
	result, err := r.db.ExecContext(ctx, "update deployment set plan_location = $1, state = $2, updated_at = $3 where id = $4",
		location, DeploymentStateScheduled, time.Now().UTC(), id)
	if err != nil {
		return errors.Wrap(err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err)
	}

	if affected == 0 {
		return errors.Wrapf("the following id: %s was not found to update in deployment table", id)
	}
	return nil
}

func (r *Repo) UpdateDeploymentResult(ctx context.Context, id uuid.UUID) (*Deployment, error) {
	var deployment Deployment

	row := r.db.QueryRowxContext(ctx, `
		update deployment set state = $1, updated_at = $2 where id = $3
		returning *
		`, DeploymentStateCompleted, time.Now().UTC(), id)

	err := row.StructScan(&deployment)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	return &deployment, nil
}

func (r *Repo) DeploymentByID(ctx context.Context, id uuid.UUID) (*Deployment, error) {
	var deployment Deployment

	row := r.db.QueryRowxContext(ctx, "select * from deployment where id = $1", id)
	err := row.StructScan(&deployment)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	return &deployment, nil
}

func (r *Repo) UpdateDeploymentReceiptHandle(ctx context.Context, id uuid.UUID, receiptHandle string) (*Deployment, error) {
	var deployment Deployment
	row := r.db.QueryRowxContext(ctx, `
		update deployment set receipt_handle = $1, updated_at = $2 where id = $3
		returning *
	`, receiptHandle, time.Now().UTC(), id)
	err := row.StructScan(&deployment)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	return &deployment, nil
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
	
	insert into deployment(environment_id, namespace_id, req_id, plan_options, plan_location, state, "user", created_at, updated_at) 
		values ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		returning (id)
	
	`, d.EnvironmentID, d.NamespaceID, d.ReqID, d.PlanOptions, d.PlanLocation, DeploymentStateQueued, d.User, d.CreatedAt, d.UpdatedAt).
		Scan(&d.ID)

	if err != nil {
		return errors.Wrap(err)
	}

	return nil
}

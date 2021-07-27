package data

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	uuid "github.com/satori/go.uuid"

	"github.com/unanet/go/pkg/errors"
	"github.com/unanet/go/pkg/json"
)

type DeploymentCronState string

const (
	DeploymentCronStateIdle    DeploymentCronState = "idle"
	DeploymentCronStateRunning DeploymentCronState = "running"
)

type DeploymentCronJob struct {
	ID          uuid.UUID           `db:"id"`
	Description string              `db:"description"`
	PlanOptions json.Object         `db:"plan_options"`
	Schedule    string              `db:"schedule"`
	LastRun     sql.NullTime        `db:"last_run"`
	State       DeploymentCronState `db:"state"`
	Disabled    bool                `db:"disabled"`
	Order       int                 `db:"exec_order"`
}

type DeploymentCronJobs []*DeploymentCronJob

func (r *Repo) getDeploymentCronJobs(ctx context.Context, tx *sqlx.Tx) (DeploymentCronJobs, error) {
	rows, err := tx.QueryxContext(ctx, "select * from deployment_cron where state = 'idle' and disabled = false order by exec_order for update")
	if err != nil {
		return nil, errors.Wrap(err)
	}
	var cronJobs DeploymentCronJobs
	for rows.Next() {
		var cron DeploymentCronJob
		err = rows.StructScan(&cron)
		if err != nil {
			return nil, errors.Wrap(err)
		}
		cronJobs = append(cronJobs, &cron)
	}

	return cronJobs, nil
}

type ExecFn func(ctx context.Context, query string, args ...interface{}) (sql.Result, error)

func (r *Repo) UpdateDeploymentCronState(ctx context.Context, id uuid.UUID, state DeploymentCronState) error {
	return r.updateDeploymentCronState(ctx, nil, id, state)
}

func (r *Repo) updateDeploymentCronState(ctx context.Context, tx *sqlx.Tx, id uuid.UUID, state DeploymentCronState) error {
	var f ExecFn
	if tx != nil {
		f = tx.ExecContext
	} else {
		f = r.db.ExecContext
	}
	result, err := f(ctx, "update deployment_cron set state = $1  where id = $2", state, id)
	if err != nil {
		return errors.Wrap(err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err)
	}

	if affected == 0 {
		return errors.Wrapf("the following id: %s was not found to update in deployment_cron table", id)
	}
	return nil
}

func (r *Repo) insertDeploymentCronJobs(ctx context.Context, tx *sqlx.Tx, cronID uuid.UUID, ids []uuid.UUID) error {
	baseSql := "insert into deployment_cron_job (deployment_cron_id, deployment_id) values "
	var values []string
	args := []interface{}{cronID}
	for i, x := range ids {
		values = append(values, fmt.Sprintf("($1, $%d)", i+2))
		args = append(args, x)
	}
	sqlStr := baseSql + strings.Join(values, ",")
	_, err := tx.ExecContext(ctx, sqlStr, args...)
	if err != nil {
		return errors.Wrap(err)
	}
	return nil
}

func (r *Repo) ScheduleDeploymentCronJobs(ctx context.Context, schedule func(context.Context, *DeploymentCronJob) ([]uuid.UUID, error)) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return errors.Wrap(err)
	}

	jobs, err := r.getDeploymentCronJobs(ctx, tx)
	if err != nil {
		return errors.WrapTx(tx, err)
	}

	for _, x := range jobs {
		ids, e := schedule(ctx, x)
		if e != nil {
			return errors.WrapTx(tx, e)
		}
		if ids == nil || len(ids) == 0 {
			continue
		}
		err = r.insertDeploymentCronJobs(ctx, tx, x.ID, ids)
		if err != nil {
			return errors.WrapTx(tx, err)
		}

		err = r.updateDeploymentCronState(ctx, tx, x.ID, DeploymentCronStateRunning)
		if err != nil {
			return errors.WrapTx(tx, err)
		}
	}
	err = tx.Commit()
	if err != nil {
		return errors.WrapTx(tx, err)
	}
	return nil
}

func (r *Repo) UpdateFinishedJobs(ctx context.Context) error {
	now := time.Now().UTC()
	_, err := r.db.ExecContext(ctx, `
		update deployment_cron set state = 'idle', last_run = $1
		where state = 'running' and
		      (select count(*) from deployment_cron_job as dcj
		    		left join deployment d on dcj.deployment_id = d.id
		    		where d.state != 'completed' and dcj.deployment_cron_id = deployment_cron.id) = 0 
	`, now)
	if err != nil {
		return errors.Wrap(err)
	}

	return nil
}

func (r *Repo) DeploymentCronJobs(ctx context.Context) ([]DeploymentCronJob, error) {
	rows, err := r.db.QueryxContext(ctx, `
		select 
			id,
			description,
			plan_options,
			schedule,
			last_run,
			state,
			disabled,
			exec_order
		from deployment_cron`)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	defer rows.Close()

	var ss []DeploymentCronJob
	for rows.Next() {
		if rows.Err() != nil {
			return nil, errors.Wrap(err)
		}

		var s DeploymentCronJob
		err = rows.StructScan(&s)
		if err != nil {
			return nil, errors.Wrap(err)
		}

		ss = append(ss, s)
	}

	return ss, nil
}

func (r *Repo) CreateDeploymentCronJob(ctx context.Context, m *DeploymentCronJob) error {

	now := time.Now().UTC()
	m.LastRun = sql.NullTime{
		Time:  now,
		Valid: true,
	}

	err := r.db.QueryRowxContext(ctx, `
	INSERT INTO deployment_cron(plan_options, schedule, state, last_run, disabled, description, exec_order)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`,
		m.PlanOptions,
		m.Schedule,
		DeploymentCronStateIdle,
		m.LastRun,
		m.Disabled,
		m.Description,
		m.Order,
	).StructScan(m)

	if err != nil {
		return errors.Wrap(err)
	}

	return nil
}

func (r *Repo) UpdateDeploymentCronJob(ctx context.Context, m *DeploymentCronJob) error {

	result, err := r.db.ExecContext(ctx, `
		update deployment_cron set plan_options = $2, schedule = $3, state = $4, disabled = $5, description = $6, exec_order = $7
		where id = $1
		RETURNING last_run
	`,
		m.ID,
		m.PlanOptions,
		m.Schedule,
		m.State,
		m.Disabled,
		m.Description,
		m.Order)
	if err != nil {
		return errors.Wrap(err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err)
	}

	if affected == 0 {
		return errors.NotFoundf("deployment cron by id: %s not found", m.ID)
	}
	return nil
}

func (r *Repo) DeleteDeploymentCronJob(ctx context.Context, id string) error {
	return r.deleteWithQuery(ctx, "deployment_cron", fmt.Sprintf("id = '%s'", id))
}

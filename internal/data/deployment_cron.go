package data

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	uuid "github.com/satori/go.uuid"

	"gitlab.unanet.io/devops/eve/pkg/errors"
	"gitlab.unanet.io/devops/eve/pkg/json"
)

type DeploymentCronState string

const (
	DeploymentCronStateIdle    DeploymentCronState = "idle"
	DeploymentCronStateRunning DeploymentCronState = "running"
)

type DeploymentCronJob struct {
	ID          uuid.UUID           `db:"id"`
	Description string              `db:"description"`
	PlanOptions json.Text           `db:"plan_options"`
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
		ids, err := schedule(ctx, x)
		if err != nil {
			return errors.WrapTx(tx, err)
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

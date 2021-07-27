package plans

import (
	"context"
	"encoding/json"
	"time"

	"github.com/robfig/cron/v3"
	uuid "github.com/satori/go.uuid"
	"github.com/unanet/go/pkg/errors"
	"github.com/unanet/go/pkg/log"
	"go.uber.org/zap"

	"github.com/unanet/eve/internal/data"
	"github.com/unanet/eve/pkg/eve"
)

type DeploymentCronRepo interface {
	ScheduleDeploymentCronJobs(ctx context.Context, schedule func(context.Context, *data.DeploymentCronJob) ([]uuid.UUID, error)) error
	UpdateFinishedJobs(ctx context.Context) error
}

type DeploymentQueuer interface {
	QueuePlan(ctx context.Context, options *eve.DeploymentPlanOptions) error
}

type DeploymentCron struct {
	log     *zap.Logger
	timeout time.Duration
	ctx     context.Context
	cancel  context.CancelFunc
	done    chan bool
	repo    DeploymentCronRepo
	dq      DeploymentQueuer
}

func NewDeploymentCron(repo DeploymentCronRepo, dq DeploymentQueuer, timeout time.Duration) *DeploymentCron {
	ctx, cancel := context.WithCancel(context.Background())
	return &DeploymentCron{
		repo:    repo,
		log:     log.Logger,
		ctx:     ctx,
		cancel:  cancel,
		dq:      dq,
		done:    make(chan bool),
		timeout: timeout,
	}
}

func (dc *DeploymentCron) Start() {
	go dc.start()
	dc.log.Info("deployment cron started")
}

func (dc *DeploymentCron) scheduler(ctx context.Context, job *data.DeploymentCronJob) ([]uuid.UUID, error) {
	schedule, err := cron.ParseStandard(job.Schedule)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	nextTime := schedule.Next(job.LastRun.Time)
	if nextTime.After(time.Now().UTC()) {
		return nil, nil
	}

	var options eve.DeploymentPlanOptions
	err = json.Unmarshal(job.PlanOptions, &options)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	err = dc.dq.QueuePlan(ctx, &options)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	return options.DeploymentIDs, nil
}

func (dc *DeploymentCron) run(ctx context.Context) error {

	err := dc.repo.ScheduleDeploymentCronJobs(ctx, dc.scheduler)
	if err != nil {
		return errors.Wrap(err)
	}

	err = dc.repo.UpdateFinishedJobs(ctx)
	if err != nil {
		return errors.Wrap(err)
	}

	return nil
}

func (dc *DeploymentCron) start() {
	for {
		select {
		case <-dc.ctx.Done():
			dc.log.Info("deployment cron stopped")
			close(dc.done)
			return
		default:
			ctx, _ := context.WithTimeout(context.WithValue(context.Background(), log.RequestIDKey, log.GetNextRequestID()), dc.timeout)
			err := dc.run(ctx)
			if err != nil {
				dc.log.Error("an error occurred in the deployment cron scheduler", zap.Error(err))
			}
		}

		time.Sleep(15 * time.Second)
	}
}

func (dc *DeploymentCron) Stop() {
	dc.cancel()
	<-dc.done
}

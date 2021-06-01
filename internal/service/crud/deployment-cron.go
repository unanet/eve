package crud

import (
	"context"
	uuid "github.com/satori/go.uuid"
	"gitlab.unanet.io/devops/eve/internal/data"
	"gitlab.unanet.io/devops/eve/pkg/eve"
	"gitlab.unanet.io/devops/go/pkg/json"
)

func (m *Manager) DeploymentCronJobs(ctx context.Context) (models []eve.DeploymentCronJob, err error) {
	dbDeploymentCronJobs, err := m.repo.DeploymentCronJobs(ctx)
	if err != nil {
		return nil, err
	}

	return fromDataDeploymentCronJobList(dbDeploymentCronJobs), err
}

func (m *Manager) CreateDeploymentCronJob(ctx context.Context, model *eve.DeploymentCronJob) error {

	dbDeploymentCron := toDataDeploymentCronJob(*model)
	if err := m.repo.CreateDeploymentCronJob(ctx, &dbDeploymentCron); err != nil {
		return err
	}

	model.ID = dbDeploymentCron.ID.String()

	return nil
}

func (m *Manager) UpdateDeploymentCronJob(ctx context.Context, model *eve.DeploymentCronJob) (err error) {
	dbModel := toDataDeploymentCronJob(*model)
	if err := m.repo.UpdateDeploymentCronJob(ctx, &dbModel); err != nil {
		return err
	}

	return nil
}

func (m *Manager) DeleteDeploymentCronJob(ctx context.Context, id string) (err error) {
	return m.repo.DeleteDeploymentCronJob(ctx, id)
}


func fromDataDeploymentCronJobList(crons []data.DeploymentCronJob) []eve.DeploymentCronJob {
	var list []eve.DeploymentCronJob
	for _, x := range crons {
		list = append(list, fromDataDeploymentCronJob(x))
	}
	return list
}

func fromDataDeploymentCronJob(dbModel data.DeploymentCronJob) eve.DeploymentCronJob {
	return eve.DeploymentCronJob{
		ID: dbModel.ID.String(),
		Description: dbModel.Description,
		PlanOptions: dbModel.PlanOptions.AsMapOrEmpty(),
		Schedule: dbModel.Schedule,
		LastRun: dbModel.LastRun.Time,
		State: dbModel.State,
		Disabled: dbModel.Disabled,
		Order: dbModel.Order,
	}
}

func toDataDeploymentCronJob(m eve.DeploymentCronJob) data.DeploymentCronJob {
	// Swallow the error if we throw one :(
	id, _ := uuid.FromString(m.ID)

	return data.DeploymentCronJob{
		ID: id,
		Description: m.Description,
		PlanOptions: json.FromMapOrEmpty(m.PlanOptions),
		Schedule: m.Schedule,
		// Omit last run
		State: m.State,
		Disabled: m.Disabled,
		Order: m.Order,
	}
}

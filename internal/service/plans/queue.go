package plans

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/unanet/eve/internal/data"
	"github.com/unanet/eve/internal/service/crud"
	"github.com/unanet/eve/pkg/eve"
	"github.com/unanet/eve/pkg/queue"
	"github.com/unanet/go/pkg/errors"
	"go.uber.org/zap"
)

type QWriter interface {
	Message(ctx context.Context, m *queue.M) error
}

type QueueWorker interface {
	Start(queue.Handler)
	Stop()
	DeleteMessage(ctx context.Context, m *queue.M) error
	// Message sends a message to a different queue given a url, not this one
	Message(ctx context.Context, qUrl string, m *queue.M) error
}

type HttpCallback interface {
	Post(ctx context.Context, url string, body interface{}) error
}

func fromDataService(s data.DeployService) *eve.DeployService {
	return &eve.DeployService{
		ServiceID:        s.ServiceID,
		ServicePort:      s.ServicePort,
		MetricsPort:      s.MetricsPort,
		ServiceName:      s.ServiceName,
		Count:            s.Count,
		SuccessExitCodes: s.SuccessExitCodes,
		DeployArtifact: &eve.DeployArtifact{
			ArtifactID:       s.ArtifactID,
			ArtifactName:     s.ArtifactName,
			RequestedVersion: s.RequestedVersion,
			DeployedVersion:  s.DeployedVersion.String,
			ImageTag:         s.ImageTag,
			Result:           eve.DeployArtifactResultNoop,
		},
		Nuance: strconv.Itoa(int(time.Now().Unix())),
	}
}

func fromDataServices(services data.DeployServices) eve.DeployServices {
	var list eve.DeployServices
	for _, x := range services {
		list = append(list, fromDataService(x))
	}
	return list
}

func fromDataJob(j data.DeployJob) *eve.DeployJob {
	return &eve.DeployJob{
		JobID:            j.JobID,
		JobName:          j.JobName,
		SuccessExitCodes: j.SuccessExitCodes,
		Nuance:           strconv.Itoa(int(time.Now().Unix())),
		DeployArtifact: &eve.DeployArtifact{
			ArtifactID:       j.ArtifactID,
			ArtifactName:     j.ArtifactName,
			RequestedVersion: j.RequestedVersion,
			DeployedVersion:  j.DeployedVersion.String,
			ImageTag:         j.ImageTag,
			Result:           eve.DeployArtifactResultNoop,
		},
	}
}

func fromDataJobs(d data.DeployJobs) eve.DeployJobs {
	var list eve.DeployJobs
	for _, x := range d {
		list = append(list, fromDataJob(x))
	}
	return list
}

type messageLogger func(format string, a ...interface{})

type Queue struct {
	worker     QueueWorker
	repo       *data.Repo
	uploader   eve.CloudUploader
	callback   HttpCallback
	downloader eve.CloudDownloader
	crud       *crud.Manager
}

func NewQueue(
	worker QueueWorker,
	repo *data.Repo,
	crud *crud.Manager,
	uploader eve.CloudUploader,
	downloader eve.CloudDownloader,
	httpCallBack HttpCallback) *Queue {
	return &Queue{
		worker:     worker,
		repo:       repo,
		crud:       crud,
		uploader:   uploader,
		downloader: downloader,
		callback:   httpCallBack,
	}
}

func (dq *Queue) Logger(ctx context.Context) *zap.Logger {
	return queue.GetLogger(ctx)
}

func (dq *Queue) Start() {
	go func() {
		dq.worker.Start(queue.HandlerFunc(dq.handleMessage))
	}()
}

func (dq *Queue) Stop() {
	dq.worker.Stop()
}

func (dq *Queue) matchArtifact(a *eve.DeployArtifact, optName string, options eve.NamespacePlanOptions, logger messageLogger) {
	// match services to be deployed
	// we need to pass in the service/database name here to match if it was supplied as we should only match services/databases that were specified
	match := options.Artifacts.Match(a.ArtifactID, optName, a.RequestedVersion)
	if match == nil {
		return
	}
	match.Matched = true
	if a.DeployedVersion == match.AvailableVersion && !options.ForceDeploy {
		if options.ArtifactsSupplied {
			if len(optName) > 0 {
				logger("service: %s, version: %s, is already up to date", optName, a.DeployedVersion)
			} else {
				logger("artifact: %s, version: %s, is already up to date", a.ArtifactName, a.DeployedVersion)
			}
		}
		return
	}
	a.AvailableVersion = match.AvailableVersion
	a.ArtifactoryPath = match.ArtifactoryPath
	a.ArtifactoryFeed = match.ArtifactoryFeed
	a.ArtifactoryFeedType = match.FeedType
	if a.AvailableVersion == "" || (a.DeployedVersion == a.AvailableVersion && !options.ForceDeploy) {
		return
	}
	a.Deploy = true
}

func (dq *Queue) setupNSDeploymentPlan(ctx context.Context, deploymentID uuid.UUID, options eve.NamespacePlanOptions) (*eve.NSDeploymentPlan, error) {
	cluster, err := dq.repo.ClusterByID(ctx, options.NamespaceRequest.ClusterID)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	plan := eve.NSDeploymentPlan{
		Namespace:         options.NamespaceRequest,
		EnvironmentName:   options.EnvironmentName,
		EnvironmentAlias:  options.EnvironmentAlias,
		CallbackURL:       options.CallbackURL,
		SchQueueUrl:       cluster.SchQueueUrl,
		DeploymentID:      deploymentID,
		Type:              options.Type,
		MetadataOverrides: options.Metadata,
	}

	plan.Namespace.ClusterName = cluster.Name

	if options.DryRun == true {
		plan.Status = eve.DeploymentPlanStatusDryrun
	} else {
		plan.Status = eve.DeploymentPlanStatusPending
	}

	return &plan, nil
}

func (dq *Queue) createServicesDeployment(ctx context.Context, deploymentID uuid.UUID, options eve.NamespacePlanOptions) (*eve.NSDeploymentPlan, error) {
	nSDeploymentPlan, err := dq.setupNSDeploymentPlan(ctx, deploymentID, options)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	dataServices, err := dq.repo.DeployedServicesByNamespaceID(ctx, options.NamespaceRequest.ID)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	services := fromDataServices(dataServices)
	for _, x := range services {
		metadata, err := dq.crud.ServiceMetadata(ctx, x.ServiceID)
		if err != nil {
			return nil, errors.Wrap(err)
		}
		x.Metadata = metadata

		definitions, err := dq.crud.ServiceDefinitionResults(ctx, x.ServiceID)
		if err != nil {
			return nil, errors.Wrap(err)
		}

		defBytes, err := json.Marshal(definitions)
		if err != nil {
			return nil, errors.Wrap(err)
		}

		x.Definition = defBytes

		dq.matchArtifact(x.DeployArtifact, x.ServiceName, options, nSDeploymentPlan.Message)
	}
	// Trap the restart command, since we don't care about matching a service (we just want to restart whatever version is currently deployed)
	if options.ArtifactsSupplied && options.Type != eve.DeploymentPlanTypeRestart {
		unmatched := options.Artifacts.UnMatched()
		for _, x := range unmatched {
			nSDeploymentPlan.Message("unmatched service: %s:%s", x.Name, x.AvailableVersion)
		}
	}
	nSDeploymentPlan.Services = services.ToDeploy()

	return nSDeploymentPlan, nil
}

func (dq *Queue) createJobsDeployment(ctx context.Context, deploymentID uuid.UUID, options eve.NamespacePlanOptions) (*eve.NSDeploymentPlan, error) {
	nSDeploymentPlan, err := dq.setupNSDeploymentPlan(ctx, deploymentID, options)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	dataJobs, err := dq.repo.DeployedJobsByNamespaceID(ctx, options.NamespaceRequest.ID)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	jobs := fromDataJobs(dataJobs)
	for _, x := range jobs {
		metadata, mErr := dq.crud.JobMetadata(ctx, x.JobID)
		if mErr != nil {
			return nil, errors.Wrap(mErr)
		}
		x.Metadata = metadata

		definition, dErr := dq.crud.JobDefinitionResults(ctx, x.JobID)
		if dErr != nil {
			return nil, errors.Wrap(dErr)
		}

		defBytes, _ := json.Marshal(definition)

		x.Definition = defBytes

		dq.matchArtifact(x.DeployArtifact, x.JobName, options, nSDeploymentPlan.Message)
	}
	if options.ArtifactsSupplied {
		unmatched := options.Artifacts.UnMatched()
		for _, x := range unmatched {
			nSDeploymentPlan.Message("unmatched job: %s:%s", x.Name, x.AvailableVersion)
		}
	}
	nSDeploymentPlan.Jobs = jobs.ToDeploy()
	return nSDeploymentPlan, nil
}

func (dq *Queue) rollbackError(ctx context.Context, m *queue.M, err error) error {
	qerr := dq.worker.DeleteMessage(ctx, m)
	if qerr != nil {
		dq.Logger(ctx).Error("an error occurred while trying to remove the message due to an error", zap.Any("queue_message", m), zap.Error(qerr))
	}
	return errors.Wrap(err)
}

func (dq *Queue) scheduleDeployment(ctx context.Context, m *queue.M) error {
	deployment, err := dq.repo.UpdateDeploymentReceiptHandle(ctx, m.ID, m.ReceiptHandle)
	if err != nil {
		return dq.rollbackError(ctx, m, err)
	}

	var options eve.NamespacePlanOptions
	err = json.Unmarshal(deployment.PlanOptions, &options)
	if err != nil {
		return dq.rollbackError(ctx, m, err)
	}

	var nsDeploymentPlan *eve.NSDeploymentPlan
	switch options.Type {
	case eve.DeploymentPlanTypeApplication, eve.DeploymentPlanTypeRestart:
		nsDeploymentPlan, err = dq.createServicesDeployment(ctx, deployment.ID, options)
	case eve.DeploymentPlanTypeJob:
		nsDeploymentPlan, err = dq.createJobsDeployment(ctx, deployment.ID, options)
	}
	if err != nil {
		return dq.rollbackError(ctx, m, err)
	}

	if len(options.CallbackURL) > 0 {
		if cErr := dq.callback.Post(ctx, options.CallbackURL, nsDeploymentPlan); cErr != nil {
			dq.Logger(ctx).Warn("schedule deployment callback failed",
				zap.Error(cErr),
				zap.String("plan_callback_url", nsDeploymentPlan.CallbackURL),
				zap.String("options_callback_url", options.CallbackURL),
			)
		}
	}

	if options.DryRun || nsDeploymentPlan.NothingToDeploy() {
		dq.Logger(ctx).Info("message deleted")
		err = dq.worker.DeleteMessage(ctx, m)
		if err != nil {
			return dq.rollbackError(ctx, m, err)
		}
		dq.Logger(ctx).Info("updating scheduled deployment", zap.Any("id", deployment.ID))
		_, err = dq.repo.UpdateDeploymentResult(ctx, deployment.ID)
		if err != nil {
			return errors.Wrap(err)
		}
		return nil
	}

	mBody, err := eve.MarshalNSDeploymentPlanToS3LocationBody(ctx, dq.uploader, nsDeploymentPlan)
	if err != nil {
		return errors.Wrap(err)
	}

	err = dq.worker.Message(ctx, nsDeploymentPlan.SchQueueUrl, &queue.M{
		ID:      deployment.ID,
		GroupID: nsDeploymentPlan.Namespace.GetQueueGroupID(),
		Body:    mBody,
		Command: nsDeploymentPlan.Type.Command(),
	})
	if err != nil {
		return dq.rollbackError(ctx, m, err)
	}

	err = dq.repo.UpdateDeploymentPlanLocation(ctx, deployment.ID, mBody)
	if err != nil {
		return dq.rollbackError(ctx, m, err)
	}

	return nil
}

func (dq *Queue) handleMessage(ctx context.Context, m *queue.M) error {
	switch m.Command {
	// This means it hasn't been sent to the scheduler yet
	case queue.CommandScheduleDeployment:
		return dq.scheduleDeployment(ctx, m)

	// This means it came back from the scheduler
	case queue.CommandUpdateDeployment:
		return dq.updateDeployment(ctx, m)

	case queue.CommandCallbackMessage:
		return dq.callbackMessage(ctx, m)
	default:
		return errors.Wrapf("unrecognized command: %s", m.Command)
	}
}

func (dq *Queue) updateDeployment(ctx context.Context, m *queue.M) error {
	dq.Logger(ctx).Info("updating message deployment", zap.Any("id", m.ID))
	deployment, err := dq.repo.UpdateDeploymentResult(ctx, m.ID)
	if err != nil {
		return errors.Wrap(err)
	}

	plan, err := eve.UnMarshalNSDeploymentFromS3LocationBody(ctx, dq.downloader, m.Body)
	if err != nil {
		return errors.Wrap(err)
	}

	for _, x := range plan.Services {
		if x.Result != eve.DeployArtifactResultSuccess {
			continue
		}

		err = dq.repo.UpdateDeployedServiceVersion(ctx, x.ServiceID, x.AvailableVersion)
		if err != nil {
			return errors.Wrap(err)
		}
	}

	for _, x := range plan.Jobs {
		if x.Result != eve.DeployArtifactResultSuccess {
			continue
		}

		err = dq.repo.UpdateDeployedJobVersion(ctx, x.JobID, x.AvailableVersion)
		if err != nil {
			return errors.Wrap(err)
		}
	}

	if len(plan.CallbackURL) > 0 {
		if cErr := dq.callback.Post(ctx, plan.CallbackURL, plan); cErr != nil {
			dq.Logger(ctx).Warn("update deployment callback failed",
				zap.Error(cErr),
				zap.String("plan_callback_url", plan.CallbackURL),
			)
		}
	}

	// Here we are deleting the original deploy message which unblocks deployments for a namespace in an environment
	// We will need to add some additional logic to this to account for certain scenarios where we should
	// Still Delete the Message that triggers this updateDeployment (like an error that returns not found or already deleted)
	err = dq.worker.DeleteMessage(ctx, &queue.M{
		ID:            deployment.ID,
		ReceiptHandle: deployment.ReceiptHandle.String,
	})
	if err != nil {
		return errors.Wrap(err)
	}

	err = dq.worker.DeleteMessage(ctx, m)
	if err != nil {
		return errors.Wrap(err)
	}

	return nil
}

func (dq *Queue) callbackMessage(ctx context.Context, m *queue.M) error {
	defer func() {
		err := dq.worker.DeleteMessage(ctx, m)
		if err != nil {
			dq.Logger(ctx).Error("sqs message removal failed", zap.Error(err))
		}
	}()

	var cm eve.CallbackMessage
	err := json.Unmarshal(m.Body, &cm)
	if err != nil {
		return errors.Wrap(err)
	}

	// curl -D '{ "messages": []}' http://eve-sch.eve:3000/callback?deployment_id=
	d, err := dq.repo.DeploymentByID(ctx, m.ID)
	if err != nil {
		dq.Logger(ctx).Warn("an error occurred trying to get the deployment from the db", zap.String("id", m.ID.String()), zap.Error(errors.Wrap(err)))
		return nil
	}

	if d.State == data.DeploymentStateCompleted {
		dq.Logger(ctx).Warn("message callback came in for a deployment that's already completed, skipping...", zap.String("id", d.ID.String()))
		return nil
	}
	var options eve.NamespacePlanOptions
	err = json.Unmarshal(d.PlanOptions, &options)
	if err != nil {
		return errors.Wrap(err)
	}

	dcm := eve.DeploymentCallbackMessage{
		DeploymentID: m.ID,
		Status:       eve.DeploymentPlanStatusMessage,
		Type:         options.Type,
		Messages:     cm.Messages,
	}

	if len(options.CallbackURL) > 0 {
		if cErr := dq.callback.Post(ctx, options.CallbackURL, dcm); cErr != nil {
			dq.Logger(ctx).Warn("callback message callback failed", zap.String("callback_url", options.CallbackURL), zap.Error(cErr), zap.String("id", d.ID.String()))
		}
	} else {
		dq.Logger(ctx).Warn("message callback came in for a deployment without a registered callback, skipping...", zap.String("id", d.ID.String()))
	}
	return nil
}

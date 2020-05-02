package service

import (
	"context"
	"encoding/json"
	"fmt"

	uuid "github.com/satori/go.uuid"

	"gitlab.unanet.io/devops/eve/internal/cloud/queue"
	"gitlab.unanet.io/devops/eve/internal/data"
	"gitlab.unanet.io/devops/eve/pkg/errors"
)

type DeploymentQueueRepo interface {
	UpdateDeploymentReceiptHandle(ctx context.Context, id uuid.UUID, receiptHandle string) (*data.Deployment, error)
	DeployedServicesByNamespaceID(ctx context.Context, namespaceID int) (data.Services, error)
	DeployedDatabaseInstancesByNamespaceID(ctx context.Context, namespaceID int) (data.DatabaseInstances, error)
	UpdateDeploymentS3PlanLocation(ctx context.Context, id uuid.UUID, location string) error
	UpdateDeploymentS3ResultLocation(ctx context.Context, id uuid.UUID, location string) (*data.Deployment, error)
}

type CloudUploader interface {
	UploadText(ctx context.Context, key string, body string) (string, error)
}

type HttpCallback interface {
	Post(ctx context.Context, url string) error
}

type DeployArtifact struct {
	ArtifactID       int    `json:"artifact_id"`
	ArtifactName     string `json:"artifact_name"`
	RequestedVersion string `json:"requested_version"`
	DeployedVersion  string `json:"deployed_version"`
	AvailableVersion string `json:"available_version"`
	Metadata         M      `json:"metadata"`
	Deploy           bool   `json:"-"`
}

type DeployService struct {
	*DeployArtifact
	ServiceID int `json:"service_id"`
}

type DeployServices []*DeployService

func (ds DeployServices) ToDeploy() DeployServices {
	var list DeployServices
	for _, x := range ds {
		if x.Deploy {
			list = append(list, x)
		}
	}
	return list
}

func fromDataService(s data.Service) *DeployService {
	return &DeployService{
		ServiceID: s.ServiceID,
		DeployArtifact: &DeployArtifact{
			ArtifactID:       s.ArtifactID,
			ArtifactName:     s.ArtifactName,
			RequestedVersion: s.RequestedVersion,
			DeployedVersion:  s.DeployedVersion.String,
			Metadata:         s.Metadata.AsMap(),
		},
	}
}

func fromDataServices(services data.Services) DeployServices {
	var list DeployServices
	for _, x := range services {
		list = append(list, fromDataService(x))
	}
	return list
}

type DeployMigration struct {
	*DeployArtifact
	DatabaseID   int    `json:"database_id"`
	DatabaseName string `json:"database_name"`
}

type DeployMigrations []*DeployMigration

func (ds DeployMigrations) ToDeploy() DeployMigrations {
	var list DeployMigrations
	for _, x := range ds {
		if x.Deploy {
			list = append(list, x)
		}
	}
	return list
}

func fromDataDatabaseInstance(s data.DatabaseInstance) *DeployMigration {
	return &DeployMigration{
		DatabaseID:   s.DatabaseID,
		DatabaseName: s.DatabaseName,
		DeployArtifact: &DeployArtifact{
			ArtifactID:       s.ArtifactID,
			ArtifactName:     s.ArtifactName,
			RequestedVersion: s.RequestedVersion,
			DeployedVersion:  s.DeployedVersion.String,
			Metadata:         s.Metadata.AsMap(),
		},
	}
}

func fromDataDatabaseInstances(d data.DatabaseInstances) DeployMigrations {
	var list DeployMigrations
	for _, x := range d {
		list = append(list, fromDataDatabaseInstance(x))
	}
	return list
}

type NSDeploymentPlan struct {
	Namespace       *NamespaceRequest `json:"namespace"`
	EnvironmentName string            `json:"environment_name"`
	Services        DeployServices    `json:"services,omitempty"`
	Migrations      DeployMigrations  `json:"migrations,omitempty"`
	Messages        []string          `json:"messages,omitempty"`
}

func (ns *NSDeploymentPlan) GroupID() string {
	return ns.Namespace.Name
}

func (ns *NSDeploymentPlan) Message(format string, a ...interface{}) {
	ns.Messages = append(ns.Messages, fmt.Sprintf(format, a...))
}

type messageLogger func(format string, a ...interface{})

type DeploymentQueue struct {
	worker   QueueWorker
	repo     DeploymentQueueRepo
	schQueue QWriter
	uploader CloudUploader
}

func NewDeploymentQueue(worker QueueWorker,
	repo DeploymentQueueRepo,
	schQueue QWriter,
	uploader CloudUploader) *DeploymentQueue {
	return &DeploymentQueue{
		worker:   worker,
		repo:     repo,
		schQueue: schQueue,
		uploader: uploader,
	}
}

func (dq *DeploymentQueue) Start() {
	go func() {
		dq.worker.Start(queue.HandlerFunc(dq.handleMessage))
	}()
}

func (dq *DeploymentQueue) Stop() {
	dq.worker.Stop()
}

func (dq *DeploymentQueue) matchArtifact(a *DeployArtifact, options NamespacePlanOptions, logger messageLogger) {
	// match services to be deployed
	match := options.Artifacts.Match(a.ArtifactID, a.RequestedVersion)
	if match == nil {
		return
	}
	match.Matched = true
	if a.DeployedVersion == match.AvailableVersion && !options.ForceDeploy {
		if options.ArtifactsSupplied {
			logger("artifact: %s, version: %s, is already up to date", a.ArtifactName, a.DeployedVersion)
		}
		return
	}
	a.AvailableVersion = match.AvailableVersion
	if a.AvailableVersion == "" || (a.DeployedVersion == a.AvailableVersion && !options.ForceDeploy) {
		return
	}
	a.Deploy = true
}

func (dq *DeploymentQueue) setupNSDeploymentPlan(options NamespacePlanOptions) *NSDeploymentPlan {
	return &NSDeploymentPlan{
		Namespace:       options.NamespaceRequest,
		EnvironmentName: options.EnvironmentName,
	}
}

func (dq *DeploymentQueue) createServicesDeployment(ctx context.Context, options NamespacePlanOptions) (*NSDeploymentPlan, error) {
	nSDeploymentPlan := dq.setupNSDeploymentPlan(options)
	dataServices, err := dq.repo.DeployedServicesByNamespaceID(ctx, options.NamespaceRequest.ID)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	services := fromDataServices(dataServices)
	for _, x := range services {
		dq.matchArtifact(x.DeployArtifact, options, nSDeploymentPlan.Message)
	}
	if options.ArtifactsSupplied {
		unmatched := options.Artifacts.UnMatched()
		for _, x := range unmatched {
			nSDeploymentPlan.Message("unmatched service: %s", x.Name)
		}
	}
	nSDeploymentPlan.Services = services.ToDeploy()
	return nSDeploymentPlan, nil
}

func (dq *DeploymentQueue) createMigrationsDeployment(ctx context.Context, options NamespacePlanOptions) (*NSDeploymentPlan, error) {
	nSDeploymentPlan := dq.setupNSDeploymentPlan(options)
	dataDatabaseInstances, err := dq.repo.DeployedDatabaseInstancesByNamespaceID(ctx, options.NamespaceRequest.ID)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	migrations := fromDataDatabaseInstances(dataDatabaseInstances)
	for _, x := range migrations {
		dq.matchArtifact(x.DeployArtifact, options, nSDeploymentPlan.Message)
	}
	if options.ArtifactsSupplied {
		unmatched := options.Artifacts.UnMatched()
		for _, x := range unmatched {
			nSDeploymentPlan.Message("unmatched service: %s", x.Name)
		}
	}
	nSDeploymentPlan.Migrations = migrations.ToDeploy()
	return nSDeploymentPlan, nil
}

func (dq *DeploymentQueue) scheduleDeployment(ctx context.Context, m *queue.M) error {
	deployment, err := dq.repo.UpdateDeploymentReceiptHandle(ctx, m.ID, m.ReceiptHandle)
	if err != nil {
		return errors.Wrap(err)
	}
	var options NamespacePlanOptions
	err = json.Unmarshal(deployment.PlanOptions, &options)
	if err != nil {
		return errors.Wrap(err)
	}
	var nsDeploymentPlan *NSDeploymentPlan
	if options.Type == DeploymentPlanTypeApplication {
		nsDeploymentPlan, err = dq.createServicesDeployment(ctx, options)
	} else {
		nsDeploymentPlan, err = dq.createMigrationsDeployment(ctx, options)
	}

	nsDeploymentPlanText, err := data.StructToJSONText(nsDeploymentPlan)
	if err != nil {
		return errors.Wrap(err)
	}

	//if len(options.CallbackURL) == 0 {
	//	err := dq.callback.Post(ctx, options.CallbackURL)
	//	if err != nil {
	//		log.Logger.Warn("callback failed", zap.String("callback_url", options.CallbackURL), zap.String("req_id", queue.GetReqID(ctx)))
	//	}
	//}

	if options.DryRun {
		err := dq.worker.DeleteMessage(ctx, m)
		if err != nil {
			return errors.Wrap(err)
		}
		return nil
	}

	location, err := dq.uploader.UploadText(ctx, fmt.Sprintf("plan-%s", deployment.ID), nsDeploymentPlanText.String())
	err = dq.schQueue.Message(&queue.M{
		ID:      deployment.ID,
		ReqID:   queue.GetReqID(ctx),
		GroupID: nsDeploymentPlan.GroupID(),
		Body:    location,
	})
	if err != nil {
		return errors.Wrap(err)
	}
	err = dq.repo.UpdateDeploymentS3PlanLocation(ctx, deployment.ID, location)
	if err != nil {
		return errors.Wrap(err)
	}

	return nil
}

func (dq *DeploymentQueue) handleMessage(ctx context.Context, m *queue.M) error {
	switch m.State {
	// This means it hasn't been sent to the scheduler yet
	case DeploymentStateQueued:
		return dq.scheduleDeployment(ctx, m)

	// This means it came back from the scheduler
	case DeploymentStateScheduled:
		return dq.completeDeployment(ctx, m)

	default:
		return fmt.Errorf("unrecognized state: %s", m.State)
	}
}

func (dq *DeploymentQueue) completeDeployment(ctx context.Context, m *queue.M) error {
	err := dq.worker.DeleteMessage(ctx, m)
	if err != nil {
		return errors.Wrap(err)
	}
	deployment, err := dq.repo.UpdateDeploymentS3ResultLocation(ctx, m.ID, m.Body)
	if err != nil {
		return errors.Wrap(err)
	}
	err = dq.worker.DeleteMessage(ctx, &queue.M{
		ID:            deployment.ID,
		ReqID:         queue.GetReqID(ctx),
		ReceiptHandle: deployment.ReceiptHandle.String,
	})
	if err != nil {
		return errors.Wrap(err)
	}
	return nil
}

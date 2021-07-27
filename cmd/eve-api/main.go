package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	casbinpgadapter "github.com/cychiuae/casbin-pg-adapter"
	"github.com/unanet/eve/pkg/s3"
	"github.com/unanet/go/pkg/identity"
	"go.uber.org/zap"

	"github.com/unanet/eve/internal/api"
	"github.com/unanet/eve/internal/data"
	"github.com/unanet/eve/internal/service/crud"
	"github.com/unanet/eve/internal/service/plans"
	"github.com/unanet/eve/internal/service/releases"
	"github.com/unanet/eve/pkg/artifactory"
	"github.com/unanet/eve/pkg/gitlab"
	"github.com/unanet/eve/pkg/queue"
	"github.com/unanet/go/pkg/log"
)

func GetPolicyModel() model.Model {
	m := model.NewModel()
	m.AddDef("r", "r", "sub, obj, act")
	m.AddDef("p", "p", "sub, obj, act")
	m.AddDef("e", "e", "some(where (p.eft == allow))")
	m.AddDef("m", "m", "r.sub == p.sub && (keyMatch(r.obj, p.obj) || keyMatch2(r.obj, p.obj)) && (r.act == p.act || p.act == \"*\")")
	return m
}

func main() {
	dbConfig := api.GetDBConfig()
	// Try to get a DB Connection
	db, err := data.GetDBWithTimeout(dbConfig.DbConnectionString(), dbConfig.DBConnectionTimeout)
	if err != nil {
		log.Logger.Panic("Failed to open Connection to DB.", zap.Error(err))
	}

	flags := api.GetFlagsConfig()

	if flags.MigrateFlag {
		err = data.MigrateDB(dbConfig.MigrationConnectionString(), dbConfig.LogLevel)
		if err != nil {
			log.Logger.Panic("Failed to load the Database Migration Tool.", zap.Error(err))
		}
	}

	if !flags.ServerFlag {
		return
	}

	config := api.GetConfig()

	awsSession, err := session.NewSession(&aws.Config{
		Region: aws.String(config.AWSRegion)},
	)
	if err != nil {
		log.Logger.Panic("Failed to create AWS Session", zap.Error(err))
	}
	apiQueue := queue.NewQ(awsSession, queue.Config{
		MaxNumberOfMessage: config.ApiQMaxNumberOfMessage,
		QueueURL:           config.ApiQUrl,
		WaitTimeSecond:     config.ApiQWaitTimeSecond,
		VisibilityTimeout:  config.ApiQVisibilityTimeout,
	})

	repo := data.NewRepo(db)
	artifactoryClient := artifactory.NewClient(config.ArtifactoryConfig)
	deploymentPlanGenerator := plans.NewPlanGenerator(repo, artifactoryClient, apiQueue)
	crudManager := crud.NewManager(repo)
	gitlabClient := gitlab.NewClient(config.GitlabConfig)
	releaseSvc := releases.NewReleaseSvc(repo, artifactoryClient, gitlabClient)

	controllers, err := api.InitializeControllers(deploymentPlanGenerator, crudManager, releaseSvc)
	if err != nil {
		log.Logger.Panic("Unable to Initialize the Controllers")
	}

	identitySvc, err := identity.NewService(config.Identity)
	if err != nil {
		log.Logger.Panic("Unable to Initialize the Identity Service Manager", zap.Error(err))
	}

	adapter, err := casbinpgadapter.NewAdapter(db.DB, "policies")
	if err != nil {
		log.Logger.Panic("failed to create casbin adaptor", zap.Error(err))
	}
	enforcer, err := casbin.NewEnforcer(GetPolicyModel(), adapter)
	if err != nil {
		log.Logger.Panic("failed to create casbin enforcer", zap.Error(err))
	}

	err = enforcer.LoadPolicy()
	if err != nil {
		log.Logger.Panic("failed to load casbin policy", zap.Error(err))
	}

	apiServer, err := api.NewApi(controllers, identitySvc, enforcer, config)
	if err != nil {
		log.Logger.Panic("Failed to Create Api App", zap.Error(err))
	}

	deploymentQueue := plans.NewQueue(
		queue.NewWorker("eve-api", apiQueue, config.ApiQWorkerTimeout),
		repo,
		crudManager,
		s3.NewUploader(awsSession, s3.Config{Bucket: config.S3Bucket}),
		s3.NewDownloader(awsSession),
		plans.NewCallback(config.HttpCallbackTimeout),
	)

	cron := plans.NewDeploymentCron(repo, deploymentPlanGenerator, config.CronTimeout)
	if !config.LocalDev {
		cron.Start()
		deploymentQueue.Start()
	}

	apiServer.Start(func() {
		cron.Stop()
		deploymentQueue.Stop()
	})
}

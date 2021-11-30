package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	casbinpgadapter "github.com/cychiuae/casbin-pg-adapter"
	"github.com/unanet/eve/internal/config"
	"github.com/unanet/eve/pkg/s3"
	"github.com/unanet/eve/pkg/scm"
	"go.uber.org/zap"

	"github.com/unanet/eve/internal/api"
	"github.com/unanet/eve/internal/data"
	"github.com/unanet/eve/internal/service/crud"
	"github.com/unanet/eve/internal/service/plans"
	"github.com/unanet/eve/internal/service/releases"
	"github.com/unanet/eve/pkg/artifactory"
	"github.com/unanet/eve/pkg/queue"
	"github.com/unanet/go/pkg/identity"
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
	dbConfig := config.GetDBConfig()
	// Try to get a DB Connection
	db, err := data.GetDBWithTimeout(dbConfig.DbConnectionString(), dbConfig.DBConnectionTimeout)
	if err != nil {
		log.Logger.Panic("Failed to open Connection to DB.", zap.Error(err))
	}

	flags := config.GetFlagsConfig()

	if flags.MigrateFlag {
		err = data.MigrateDB(dbConfig.MigrationConnectionString(), dbConfig.LogLevel)
		if err != nil {
			log.Logger.Panic("Failed to load the Database Migration Tool.", zap.Error(err))
		}
	}

	if !flags.ServerFlag {
		return
	}

	cfg := config.GetConfig()

	awsSession, err := session.NewSession(&aws.Config{
		Region: aws.String(cfg.AWSRegion),
	},
	)
	if err != nil {
		log.Logger.Panic("Failed to create AWS Session", zap.Error(err))
	}
	apiQueue := queue.NewQ(awsSession, queue.Config{
		MaxNumberOfMessage: cfg.ApiQMaxNumberOfMessage,
		QueueURL:           cfg.ApiQUrl,
		WaitTimeSecond:     cfg.ApiQWaitTimeSecond,
		VisibilityTimeout:  cfg.ApiQVisibilityTimeout,
	})

	repo := data.NewRepo(db)
	artifactoryClient := artifactory.NewClient(cfg.ArtifactoryConfig)
	deploymentPlanGenerator := plans.NewPlanGenerator(repo, artifactoryClient, apiQueue)
	crudManager := crud.NewManager(repo)
	scmClient := scm.New()
	releaseSvc := releases.NewReleaseSvc(repo, artifactoryClient, scmClient, crudManager)

	controllers, err := api.InitializeControllers(deploymentPlanGenerator, crudManager, releaseSvc)
	if err != nil {
		log.Logger.Panic("Unable to Initialize the Controllers")
	}

	var validatorOptions []identity.ValidatorOption

	if cfg.SigningKey != "" {
		validatorOptions = append(validatorOptions, identity.JWTClientValidatorOpt(cfg.SigningKey))
	}

	identityValidatorSvc, err := identity.NewValidator(cfg.Identity, validatorOptions...)
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

	apiServer, err := api.NewApi(controllers, identityValidatorSvc, enforcer, cfg)
	if err != nil {
		log.Logger.Panic("Failed to Create Api App", zap.Error(err))
	}

	deploymentQueue := plans.NewQueue(
		queue.NewWorker("eve-api", apiQueue, cfg.ApiQWorkerTimeout),
		repo,
		crudManager,
		s3.NewUploader(awsSession, s3.Config{Bucket: cfg.S3Bucket}),
		s3.NewDownloader(awsSession),
		plans.NewCallback(cfg.HttpCallbackTimeout),
	)

	cron := plans.NewDeploymentCron(repo, deploymentPlanGenerator, cfg.CronTimeout)
	if !cfg.LocalDev {
		cron.Start()
		deploymentQueue.Start()
	}

	apiServer.Start(func() {
		cron.Stop()
		deploymentQueue.Stop()
	})
}

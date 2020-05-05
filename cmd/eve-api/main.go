package main

import (
	"flag"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"go.uber.org/zap"

	"gitlab.unanet.io/devops/eve/internal/api"
	"gitlab.unanet.io/devops/eve/internal/data"
	"gitlab.unanet.io/devops/eve/internal/service"
	"gitlab.unanet.io/devops/eve/pkg/log"
	"gitlab.unanet.io/devops/eve/pkg/mux"
	"gitlab.unanet.io/devops/eve/pkg/queue"
	"gitlab.unanet.io/devops/eve/pkg/s3"
)

func main() {
	serverFlag := flag.Bool("server", false, "start api server")
	migrateFlag := flag.Bool("migrate-db", false, "run migration")
	dropDBFlag := flag.Bool("drop-db", false, "drop db")
	flag.Parse()
	dbConfig := api.GetDBConfig()
	// Try to get a DB Connection
	db, err := data.GetDBWithTimeout(dbConfig.DbConnectionString(), dbConfig.DBConnectionTimeout)
	if err != nil {
		log.Logger.Panic("Failed to open Connection to DB.", zap.Error(err))
	}

	if *migrateFlag || *dropDBFlag {
		err = data.MigrateDB(dbConfig.MigrationConnectionString(), dbConfig.LogLevel, *dropDBFlag)
		if err != nil {
			log.Logger.Panic("Failed to load the Database Migration Tool.", zap.Error(err))
		}
	}

	if *dropDBFlag {
		return
	}

	if !*serverFlag {
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
	controllers, err := api.InitializeControllers(config, repo, apiQueue)
	if err != nil {
		log.Logger.Panic("Unable to Initialize the Controllers")
	}
	api, err := mux.NewApi(controllers, config.MuxConfig)
	if err != nil {
		log.Logger.Panic("Failed to Create Api App", zap.Error(err))
	}

	s3Uploader := s3.NewUploader(awsSession, s3.Config{
		Bucket: config.S3Bucket,
	})

	httpCallBack := service.NewCallback(config.HttpCallbackTimeout)

	deploymentQueue := service.NewDeploymentQueue(queue.NewWorker("eve-api", apiQueue, config.ApiQWorkerTimeout), repo, s3Uploader, httpCallBack)
	deploymentQueue.Start()

	api.Start(func() {
		deploymentQueue.Stop()
	})
}

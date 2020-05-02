package main

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"go.uber.org/zap"

	"gitlab.unanet.io/devops/eve/internal/api"
	"gitlab.unanet.io/devops/eve/internal/cloud/queue"
	"gitlab.unanet.io/devops/eve/internal/cloud/s3"
	"gitlab.unanet.io/devops/eve/internal/data"
	"gitlab.unanet.io/devops/eve/internal/service"
	"gitlab.unanet.io/devops/eve/pkg/log"
	"gitlab.unanet.io/devops/eve/pkg/mux"
)

func main() {
	config := api.GetConfig()
	// Try to get a DB Connection

	db, err := data.GetDBWithTimeout(config.DbConnectionString(), 10*time.Minute)
	if err != nil {
		log.Logger.Panic("Failed to open Connection to DB.", zap.Error(err))
	}
	err = data.MigrateDB(config.MigrationConnectionString(), config.LogLevel)
	if err != nil {
		log.Logger.Panic("Failed to load the Database Migration Tool.", zap.Error(err))
	}

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

	schQueue := queue.NewQ(awsSession, queue.Config{
		QueueURL: config.SchQUrl,
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

	deploymentQueue := service.NewDeploymentQueue(queue.NewWorker("eve-api", apiQueue, config.ApiQWorkerTimeout), repo, schQueue, s3Uploader)
	deploymentQueue.Start()

	api.Start(func() {
		deploymentQueue.Stop()
	})
}

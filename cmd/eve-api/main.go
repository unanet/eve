package main

import (
	"time"

	"go.uber.org/zap"

	"gitlab.unanet.io/devops/eve/internal/api"
	"gitlab.unanet.io/devops/eve/internal/config"
	"gitlab.unanet.io/devops/eve/internal/data"
	"gitlab.unanet.io/devops/eve/pkg/log"
	"gitlab.unanet.io/devops/eve/pkg/mux"
)

func main() {
	// Try to get a DB Connection
	db, err := data.GetDBWithTimeout(10 * time.Minute)
	if err != nil {
		log.Logger.Panic("Failed to open Connection to DB.", zap.Error(err))
	}
	err = data.MigrateDB(config.Values().MigrationConnectionString())
	if err != nil {
		log.Logger.Panic("Failed to load the Database Migration Tool.", zap.Error(err))
	}
	api, err := mux.NewApi(api.InitializeControllers(db), config.Values().MuxConfig)
	if err != nil {
		log.Logger.Panic("Failed to Create Api App", zap.Error(err))
	}
	api.Start()
}

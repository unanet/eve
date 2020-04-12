package main

import (
	"go.uber.org/zap"

	"gitlab.unanet.io/devops/eve/internal/api"
	"gitlab.unanet.io/devops/eve/pkg/log"
)

func main() {
	app, err := api.NewApp()
	if err != nil {
		log.Logger.Panic("Failed to Load Api App", zap.Error(err))
	}
	app.Start()
}

package main

import (
	"go.uber.org/zap"

	"gitlab.unanet.io/devops/eve/internal/api"
	"gitlab.unanet.io/devops/eve/internal/config"
	"gitlab.unanet.io/devops/eve/pkg/log"
	"gitlab.unanet.io/devops/eve/pkg/mux"
)

func main() {
	api, err := mux.NewApi(api.Controllers, config.Values().MuxConfig)
	if err != nil {
		log.Logger.Panic("Failed to Create Api App", zap.Error(err))
	}
	api.Start()
}

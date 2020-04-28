package api

import (
	"github.com/jmoiron/sqlx"

	"gitlab.unanet.io/devops/eve/internal/config"
	"gitlab.unanet.io/devops/eve/internal/data"
	"gitlab.unanet.io/devops/eve/internal/service"
	"gitlab.unanet.io/devops/eve/pkg/artifactory"
	"gitlab.unanet.io/devops/eve/pkg/mux"
)

func InitializeControllers(db *sqlx.DB) []mux.EveController {
	repo := data.NewRepo(db)
	artifactoryClient := artifactory.NewClient(config.Values().ArtifactoryConfig)
	deploymentPlanGenerator := service.NewDeploymentPlanGenerator(repo, artifactoryClient, nil)

	return []mux.EveController{
		NewPingController(),
		NewDeploymentPlanController(deploymentPlanGenerator),
	}
}

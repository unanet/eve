package api

import (
	"gitlab.unanet.io/devops/eve/internal/cloud/queue"
	"gitlab.unanet.io/devops/eve/internal/data"
	"gitlab.unanet.io/devops/eve/internal/service"
	"gitlab.unanet.io/devops/eve/pkg/artifactory"
	"gitlab.unanet.io/devops/eve/pkg/mux"
)

func InitializeControllers(c Config, repo *data.Repo, apiQueue *queue.Q) ([]mux.EveController, error) {

	artifactoryClient := artifactory.NewClient(c.ArtifactoryConfig)
	deploymentPlanGenerator := service.NewDeploymentPlanGenerator(repo, artifactoryClient, apiQueue)

	return []mux.EveController{
		NewPingController(),
		NewDeploymentPlanController(deploymentPlanGenerator),
	}, nil
}

package api

import (
	"gitlab.unanet.io/devops/eve/internal/service"
	"gitlab.unanet.io/devops/eve/pkg/mux"
)

func InitializeControllers(deploymentPlanGenerator *service.DeploymentPlanGenerator) ([]mux.EveController, error) {
	return []mux.EveController{
		NewPingController(),
		NewDeploymentPlanController(deploymentPlanGenerator),
	}, nil
}

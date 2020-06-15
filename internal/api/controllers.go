package api

import (
	"gitlab.unanet.io/devops/eve/internal/service/crud"
	"gitlab.unanet.io/devops/eve/internal/service/deployments"
	"gitlab.unanet.io/devops/eve/pkg/mux"
)

func InitializeControllers(deploymentPlanGenerator *deployments.PlanGenerator, manager *crud.Manager) ([]mux.EveController, error) {
	return []mux.EveController{
		NewPingController(),
		NewDeploymentPlanController(deploymentPlanGenerator),
		NewCrudController(manager),
	}, nil
}

package api

import (
	"gitlab.unanet.io/devops/eve/internal/service/crud"
	"gitlab.unanet.io/devops/eve/internal/service/deployments"
	"gitlab.unanet.io/devops/eve/internal/service/releases"
	"gitlab.unanet.io/devops/eve/pkg/mux"
)

func InitializeControllers(
	deploymentPlanGenerator *deployments.PlanGenerator,
	manager *crud.Manager,
	releaseSvc *releases.ReleaseSvc,
) ([]mux.EveController, error) {
	return []mux.EveController{
		NewPingController(),
		NewDeploymentPlanController(deploymentPlanGenerator),
		NewCrudController(manager),
		NewReleaseController(releaseSvc),
	}, nil
}

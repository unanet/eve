package api

import (
	"gitlab.unanet.io/devops/eve/internal/service/crud"
	"gitlab.unanet.io/devops/eve/internal/service/plans"
	"gitlab.unanet.io/devops/eve/internal/service/releases"
	"gitlab.unanet.io/devops/eve/pkg/mux"
)

func InitializeControllers(
	deploymentPlanGenerator *plans.PlanGenerator,
	manager *crud.Manager,
	releaseSvc *releases.ReleaseSvc,
) ([]mux.EveController, error) {
	return []mux.EveController{
		NewPingController(),
		NewDeploymentPlanController(deploymentPlanGenerator),
		NewPodController(manager),
		NewReleaseController(releaseSvc),
		NewMetadataController(manager),
		NewEnvironmentController(manager),
		NewNamespaceController(manager),
		NewServiceController(manager),
	}, nil
}

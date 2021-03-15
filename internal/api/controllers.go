package api

import (
	"gitlab.unanet.io/devops/eve/internal/service/crud"
	"gitlab.unanet.io/devops/eve/internal/service/plans"
	"gitlab.unanet.io/devops/eve/internal/service/releases"
)

func InitializeControllers(
	deploymentPlanGenerator *plans.PlanGenerator,
	manager *crud.Manager,
	releaseSvc *releases.ReleaseSvc,
) ([]Controller, error) {
	return []Controller{
		NewPingController(),
		NewDeploymentsController(manager),
		NewDeploymentPlansController(deploymentPlanGenerator),
		NewPodController(manager),
		NewReleaseController(releaseSvc),
		NewMetadataController(manager),
		NewEnvironmentController(manager),
		NewNamespaceController(manager),
		NewServiceController(manager),
		NewJobController(manager),
		NewDefinitionController(manager),
	}, nil
}

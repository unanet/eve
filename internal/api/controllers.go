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
		NewArtifactController(manager),
		NewClusterController(manager),
		NewDefinitionsController(manager),
		NewDeploymentPlansController(deploymentPlanGenerator),
		NewDeploymentsController(manager),
		NewDeploymentsCronController(manager),
		NewEnvironmentController(manager),
		NewReleaseController(releaseSvc),
		NewFeedController(manager),
		NewJobController(manager),
		NewEnvironmentFeedMapController(manager),
		NewMetadataController(manager),
		NewNamespaceController(manager),
		NewServiceController(manager),
	}, nil
}

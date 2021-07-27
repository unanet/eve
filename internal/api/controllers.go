package api

import (
	"github.com/unanet/eve/internal/service/crud"
	"github.com/unanet/eve/internal/service/plans"
	"github.com/unanet/eve/internal/service/releases"
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

package service_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"gitlab.unanet.io/devops/eve/internal/config"
	"gitlab.unanet.io/devops/eve/internal/data"
	"gitlab.unanet.io/devops/eve/internal/service"
	"gitlab.unanet.io/devops/eve/pkg/artifactory"
)

func TestDeploymentPlanGenerator_Generate(t *testing.T) {
	dpg := service.NewDeploymentPlanGenerator(data.NewRepo(), artifactory.NewClient(config.Values().ArtifactoryConfig))
	err := dpg.Generate(context.TODO(), service.DeploymentPlanOptions{
		Environment: "int",
		Namespaces:  nil,
		Services:    nil,
	})
	require.NoError(t, err)
	fmt.Println(dpg.Plan)
}

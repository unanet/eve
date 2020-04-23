package service_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"gitlab.unanet.io/devops/eve/internal/config"
	"gitlab.unanet.io/devops/eve/internal/data"
	"gitlab.unanet.io/devops/eve/internal/service"
	"gitlab.unanet.io/devops/eve/pkg/artifactory"
)

func TestDeploymentPlanGenerator_Generate(t *testing.T) {
	dpg := service.NewDeploymentPlanGenerator(data.NewRepo(nil), artifactory.NewClient(config.Values().ArtifactoryConfig))
	result, err := dpg.GenerateDeploymentPlan(context.TODO(), service.PlanOptions{
		Environment:      "int",
		NamespaceAliases: nil,
		Services:         nil,
	})
	require.NoError(t, err)
	jsonData, err := json.Marshal(result)
	require.NoError(t, err)
	fmt.Println(string(jsonData))
}

// +build local

package service_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"gitlab.unanet.io/devops/eve/internal/api"
	"gitlab.unanet.io/devops/eve/internal/data"
	"gitlab.unanet.io/devops/eve/internal/service"
	"gitlab.unanet.io/devops/eve/pkg/artifactory"
)

func TestPlanGenerator_GenerateMigrationPlan(t *testing.T) {
	db, err := data.GetDBWithTimeout(time.Second * 10)
	require.NoError(t, err)
	repo := data.NewRepo(db)
	jfrog := artifactory.NewClient(api.GetConfig().ArtifactoryConfig)
	pg := service.NewDeploymentPlanGenerator(repo, jfrog)
	plan, err := pg.GenerateApplicationPlan(context.TODO(), service.DeploymentPlanOptions{
		Environment:      "int",
		NamespaceAliases: service.StringList{"cvs"},
		Artifacts:        nil,
		ForceDeploy:      false,
		DryRun:           false,
	})
	require.NoError(t, err)
	json, err := json.Marshal(&plan)
	require.NoError(t, err)
	fmt.Println(string(json))
}

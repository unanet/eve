// +build local

package data_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"gitlab.unanet.io/devops/eve/internal/api"
	"gitlab.unanet.io/devops/eve/internal/data"
)

var (
	cachedRepo *data.Repo
)

func getRepo(t *testing.T) *data.Repo {
	if cachedRepo != nil {
		return cachedRepo
	}
	db, err := data.GetDBWithTimeout(api.GetDBConfig().DbConnectionString(), 10*time.Second)
	require.NoError(t, err)
	cachedRepo = data.NewRepo(db)
	return cachedRepo
}

func TestRepo_CreateDeployment(t *testing.T) {
	repo := getRepo(t)
	result, err := repo.DeployedServicesByNamespaceID(context.TODO(), 1)
	require.NoError(t, err)
	json, err := json.Marshal(result)
	require.NoError(t, err)
	fmt.Println(string(json))

}

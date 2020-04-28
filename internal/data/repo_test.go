// +build local

package data_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"gitlab.unanet.io/devops/eve/internal/data"
)

var (
	cachedRepo *data.Repo
)

func getRepo(t *testing.T) *data.Repo {
	if cachedRepo != nil {
		return cachedRepo
	}
	db, err := data.GetDBWithTimeout(time.Second * 10)
	require.NoError(t, err)
	cachedRepo = data.NewRepo(db)
	return cachedRepo
}

func TestRepo_DatabaseInstancesByNamespaceIDs(t *testing.T) {
	repo := getRepo(t)
	var ids = []interface{}{1, 2, 3}

	result, err := repo.DeployedDatabaseInstancesByNamespaceIDs(context.TODO(), ids)
	require.NoError(t, err)
	fmt.Println(result)
}

func TestRepo_CreateDeployment(t *testing.T) {
	repo := getRepo(t)
	blah := data.RequestArtifact{
		ArtifactID:       1,
		ArtifactName:     "lbah",
		ProviderGroup:    "blah",
		FeedName:         "",
		ArtifactMetadata: nil,
		ServerMetadata:   nil,
		RequestedVersion: "",
	}
	jsonText, err := data.StructToJSONText(&blah)
	require.NoError(t, err)
	d := data.Deployment{
		EnvironmentID: 1,
		NamespaceID:   1,
		ReqID:         "testing",
		PlanOptions:   jsonText,
	}
	err = repo.CreateDeployment(context.TODO(), &d)
	require.NoError(t, err)
	fmt.Println(d)

}

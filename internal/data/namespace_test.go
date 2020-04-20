package data_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"gitlab.unanet.io/devops/eve/internal/data"
)

func getRepo() *data.Repo {
	return &data.Repo{}
}

func TestRepo_GetNamespaces(t *testing.T) {
	namespaces, err := getRepo().GetNamespaces(context.TODO())
	require.NoError(t, err)
	fmt.Println(namespaces)
}

func TestRepo_GetNamespaceByID(t *testing.T) {
	//namespaces, err := getRepo().GetNamespaces(context.TODO(), data.WhereEnvironmentID(1))
	namespaces, err := getRepo().GetNamespaces(context.TODO(), data.Where("environment_id", 1), data.Where("cluster_id", 1))
	require.NoError(t, err)
	fmt.Println(namespaces)
}

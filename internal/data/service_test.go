package data_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"gitlab.unanet.io/devops/eve/internal/data"
)

func TestRepo_getServices(t *testing.T) {
	repo := data.NewRepo()
	services, err := repo.RequestedArtifacts(context.TODO(), []int{1, 2, 3})
	require.NoError(t, err)
	fmt.Println(services)
}

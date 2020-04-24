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

func TestRepo_DatabaseInstancesByNamespaceIDs(t *testing.T) {
	db, err := data.GetDBWithTimeout(time.Second * 10)
	require.NoError(t, err)
	repo := data.NewRepo(db)
	var ids = []interface{}{1, 2, 3}

	result, err := repo.DeployedDatabaseInstancesByNamespaceIDs(context.TODO(), ids)
	require.NoError(t, err)
	fmt.Println(result)
}

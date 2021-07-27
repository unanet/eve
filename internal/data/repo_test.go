// +build local

package data_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/unanet/eve/internal/api"
	"github.com/unanet/eve/internal/data"
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

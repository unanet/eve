// +build local

package secrets_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"gitlab.unanet.io/devops/eve/internal/config"
	"gitlab.unanet.io/devops/eve/internal/secrets"
)

var (
	c *secrets.Client
)

func client(t *testing.T) *secrets.Client {
	if c != nil {
		return c
	}

	cl, err := secrets.NewClient(config.Values().VaultConfig)
	require.NoError(t, err)
	c = cl
	require.NotNil(t, c)
	return c
}

func TestClient_GetKVSecret(t *testing.T) {
	resp, err := client(t).GetKVSecret("devops/artifactory", "ci_readonly_username")
	require.NoError(t, err)
	require.Equal(t, "unanet-ci-r", resp)
}

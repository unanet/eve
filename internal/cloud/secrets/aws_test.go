// +build local

package secrets_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"gitlab.unanet.io/devops/eve/internal/cloud/secrets"
)

func TestGetSecret(t *testing.T) {
	result, err := secrets.GetSecret("artifactory")
	require.NoError(t, err)
	require.Equal(t, "unanet-ci-r", result["ci_readonly_username"])
}

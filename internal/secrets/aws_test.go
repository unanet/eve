// +build local

package secrets_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"gitlab.unanet.io/devops/eve/internal/secrets"
)

func TestAWSGetSecret(t *testing.T) {
	result, err := secrets.GetAWSSecret("artifactory")
	require.NoError(t, err)
	require.Equal(t, "unanet-ci-r", result["ci_readonly_username"])
}

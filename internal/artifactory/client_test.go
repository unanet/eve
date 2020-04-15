// +build local

package artifactory_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"gitlab.unanet.io/devops/eve/internal/artifactory"
	"gitlab.unanet.io/devops/eve/internal/config"
)

var (
	c *artifactory.Client
)

func client(t *testing.T) *artifactory.Client {
	if c != nil {
		return c
	}
	c = artifactory.NewClient(config.Values().ArtifactoryConfig)
	require.NotNil(t, c)
	return c
}

func TestClient_GetLatestVersion_Success(t *testing.T) {
	resp, err := client(t).GetLatestVersion(context.TODO(), "docker", "eve-api", "0.1.0.*")
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Contains(t, resp.Version, "0.1.0")
}

func TestClient_GetLatestVersion_NotFound(t *testing.T) {
	_, err := client(t).GetLatestVersion(context.TODO(), "docker", "eve-api", "99.0.0.*")
	ae, ok := err.(artifactory.ErrorResponse)
	require.True(t, ok)
	require.Equal(t, 404, ae.Errors[0].Status)
}

func TestClient_CopyArtifact_Success(t *testing.T) {
	resp, err := client(t).CopyArtifact(context.TODO(), "docker-int-local",
		"unanet/platform/0.1.0.254", "docker-qa-local", "unanet/platform/0.1.0.254", false)
	require.NoError(t, err)
	require.NotNil(t, resp)
}

func TestClient_CopyArtifact_Failure(t *testing.T) {
	_, err := client(t).CopyArtifact(context.TODO(), "", "", "", "", true)
	_, ok := err.(artifactory.ErrorResponse)
	require.True(t, ok)
}

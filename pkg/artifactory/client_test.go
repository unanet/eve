// +build local

package artifactory_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"gitlab.unanet.io/devops/eve/internal/api"
	"gitlab.unanet.io/devops/eve/pkg/artifactory"
)

var (
	c *artifactory.Client
)

func client(t *testing.T) *artifactory.Client {
	if c != nil {
		return c
	}
	c = artifactory.NewClient(api.GetConfig().ArtifactoryConfig)
	require.NotNil(t, c)
	return c
}

func TestClient_GetLatestVersion_UnanetDockerSuccess(t *testing.T) {
	resp, err := client(t).GetLatestVersion(context.TODO(), "docker", "ops", "2020.*")
	require.NoError(t, err)
	require.NotNil(t, resp)
	fmt.Println(resp)
}

func TestClient_GetLatestVersion_UnanetGenericSuccess(t *testing.T) {
	resp, err := client(t).GetLatestVersion(context.TODO(), "generic-int", "clearview/infocus-reports", "2020.2.*")
	require.NoError(t, err)
	require.NotNil(t, resp)
	fmt.Println(resp)
}

func TestClient_GetLatestVersion_Success(t *testing.T) {
	resp, err := client(t).GetLatestVersion(context.TODO(), "docker", "eve-api", "0.1.0.*")
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Contains(t, resp, "0.1.0")
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

func TestClient_GetLatestVersionLessThan(t *testing.T) {
	resp, err := client(t).GetLatestVersionLessThan(context.TODO(), "generic-int", "unanet/unanet", "0.3")
	fmt.Println(resp)
	require.NoError(t, err)
}

func TestClient_GetArtifactProperties_Success(t *testing.T) {
	resp, err := client(t).GetArtifactProperties(context.TODO(), "docker-int", "unanet/unanet/20.2.0.1992")
	require.NoError(t, err)
	require.Equal(t, "20.2.0.1992", resp.Property("version"))
}

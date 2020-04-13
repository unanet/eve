// +build local

package artifactory_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

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

	var err error
	c, err = artifactory.NewClient(config.Values().ArtifactoryConfig)
	assert.NoError(t, err)
	return c
}

func TestClient_GetLatestVersion(t *testing.T) {
	resp, err := client(t).GetLatestVersion(context.TODO(), "docker", "eve-api", "0.1.0.*")
	assert.NoError(t, err)
	assert.Contains(t, resp.Version, "0.1.0")
}

//func TestClient_CopyArtifactVersion(t *testing.T) {
//	resp, err := client(t).CopyArtifact(context.TODO(), "repository", "src_path", "dest_repository", "dest_path")
//}

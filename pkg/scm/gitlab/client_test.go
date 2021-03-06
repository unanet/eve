// +build local

package gitlab_test

import (
	"context"
	"github.com/unanet/eve/internal/config"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/unanet/eve/pkg/scm/gitlab"
)

var (
	c *gitlab.Client
)

func client(t *testing.T) *gitlab.Client {
	if c != nil {
		return c
	}
	c = gitlab.NewClient(config.GetConfig().GitlabConfig)
	require.NotNil(t, c)
	return c
}

func TestClient_TagCommit_Failure(t *testing.T) {
	_, err := client(t).TagCommit(context.TODO(), gitlab.TagOptions{
		ProjectID: 0,
		TagName:   "",
		GitHash:   "",
	})
	_, ok := err.(gitlab.ErrorResponse)
	require.True(t, ok)
}

func TestClient_TagCommit_Success(t *testing.T) {
	resp, err := client(t).TagCommit(context.TODO(), gitlab.TagOptions{
		ProjectID: 201,
		TagName:   "Testing",
		GitHash:   "b3e203c5857accf29196ea7c626aa8cbc9c072cb",
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, "Testing", resp.Name)
}

// +build local

package fn_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"gitlab.unanet.io/devops/eve/internal/cloud/fn"
)

var (
	c *fn.Client
)

func client(t *testing.T) *fn.Client {
	if c != nil {
		return c
	}
	c = fn.NewClient(60 * time.Second)
	require.NotNil(t, c)
	return c
}

func TestClient_Azure_Success(t *testing.T) {
	resp, err := client(t).ExecuteFn(context.TODO(), fn.Azure("unanet-cloudops",
		"ping", "u2RaEqmUcjLEQkxQfOVPFZ62Kz9HhvlhGMlJbSS8PSu3wk/oetDFSw=="),
		fn.MapBody(map[string]interface{}{"test": "test"}))
	require.NoError(t, err)
	require.Equal(t, "test", resp["test"])
}

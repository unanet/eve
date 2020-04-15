package fn

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/dghubble/sling"

	"gitlab.unanet.io/devops/eve/pkg/slinge"
)

const (
	userAgent = "eve-function"
)

type Client struct {
	sling *sling.Sling
}

func NewClient(timeout time.Duration) *Client {
	var httpClient = &http.Client{
		Timeout: timeout,
	}

	sling := sling.New().Client(httpClient).
		Add("User-Agent", userAgent).
		ResponseDecoder(slinge.NewJsonDecoder())
	return &Client{sling: sling}
}

func (c *Client) ExecuteFn(ctx context.Context, arguments ...Argument) (map[string]interface{}, error) {
	var success map[string]interface{}
	var failure string
	for _, s := range arguments {
		c.sling = s(c.sling)
	}
	r, err := c.sling.Post("").Request()
	if err != nil {
		return nil, err
	}
	resp, err := c.sling.Do(r.WithContext(ctx), &success, &failure)
	if err != nil {
		return nil, err
	}
	if http.StatusOK == resp.StatusCode {
		return success, nil
	} else {
		return nil, fmt.Errorf(failure)
	}
}

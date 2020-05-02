package scheduler

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/dghubble/sling"

	"gitlab.unanet.io/devops/eve/pkg/errors"
	"gitlab.unanet.io/devops/eve/pkg/json"
)

const (
	userAgent = "eve-scheduler"
)

type HttpCallback struct {
	sling *sling.Sling
}

func NewHttpCallback(timeout time.Duration) *HttpCallback {
	var httpClient = &http.Client{
		Timeout: timeout,
	}

	sling := sling.New().Client(httpClient).
		Add("User-Agent", userAgent).
		ResponseDecoder(json.NewJsonDecoder())
	return &HttpCallback{sling: sling}
}

func (c *HttpCallback) Post(ctx context.Context, url string) error {
	var failure string
	r, err := c.sling.New().Post(url).Request()
	if err != nil {
		return errors.Wrap(err)
	}
	resp, err := c.sling.Do(r.WithContext(ctx), nil, &failure)
	if err != nil {
		return errors.Wrap(err)
	}
	if http.StatusOK == resp.StatusCode {
		return nil
	} else {
		return errors.Wrap(fmt.Errorf(failure))
	}
}

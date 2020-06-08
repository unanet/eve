package deployments

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/dghubble/sling"

	"gitlab.unanet.io/devops/eve/pkg/errors"
	ehttp "gitlab.unanet.io/devops/eve/pkg/http"
	"gitlab.unanet.io/devops/eve/pkg/json"
)

const (
	userAgent = "eve"
)

type Callback struct {
	sling *sling.Sling
}

func NewCallback(timeout time.Duration) *Callback {
	var httpClient = &http.Client{
		Timeout:   timeout,
		Transport: ehttp.LoggingTransport,
	}

	sling := sling.New().Client(httpClient).
		Add("User-Agent", userAgent).
		ResponseDecoder(json.NewJsonDecoder())
	return &Callback{sling: sling}
}

func (c *Callback) Post(ctx context.Context, url string, body interface{}) error {
	var failure string
	r, err := c.sling.New().Post(url).BodyJSON(body).Request()
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

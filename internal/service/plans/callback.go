package plans

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/dghubble/sling"

	"github.com/unanet/go/pkg/errors"
	ehttp "github.com/unanet/go/pkg/http"
	"github.com/unanet/go/pkg/json"
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

	slingClient := sling.New().Client(httpClient).
		Add("User-Agent", userAgent).
		ResponseDecoder(json.NewJsonDecoder())
	return &Callback{sling: slingClient}
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

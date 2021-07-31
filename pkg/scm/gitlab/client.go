package gitlab

import (
	"context"
	"fmt"
	"github.com/unanet/eve/pkg/scm/types"
	"net/http"
	"strings"
	"time"

	"github.com/dghubble/sling"

	ehttp "github.com/unanet/go/pkg/http"
	"github.com/unanet/go/pkg/json"
)

const (
	userAgent = "eve-gitlab"
)

type Config struct {
	GitlabApiKey  string        `envconfig:"GITLAB_API_KEY"`
	GitlabBaseUrl string        `envconfig:"GITLAB_BASE_URL"`
	GitlabTimeout time.Duration `envconfig:"GITLAB_TIMEOUT" default:"20s"`
}

type Client struct {
	sling *sling.Sling
}

func NewClient(config Config) *Client {
	var httpClient = &http.Client{
		Timeout:   config.GitlabTimeout,
		Transport: ehttp.LoggingTransport,
	}

	if !strings.HasSuffix(config.GitlabBaseUrl, "/") {
		config.GitlabBaseUrl += "/"
	}

	s := sling.New().Base(config.GitlabBaseUrl).Client(httpClient).
		Add("PRIVATE-TOKEN", config.GitlabApiKey).
		Add("User-Agent", userAgent).
		ResponseDecoder(json.NewJsonDecoder())
	return &Client{sling: s}
}

func (c *Client) TagCommit(ctx context.Context, options types.TagOptions) (*types.Tag, error) {
	var success types.Tag
	var failure types.ErrorResponse
	r, err := c.sling.New().Post(fmt.Sprintf("v4/projects/%d/repository/tags", options.ProjectID)).QueryStruct(options).Request()
	if err != nil {
		return nil, err
	}
	resp, err := c.sling.Do(r.WithContext(ctx), &success, &failure)
	if err != nil {
		return nil, err
	}
	if http.StatusCreated == resp.StatusCode {
		return &success, nil
	} else {
		return nil, failure
	}
}

func (c *Client) GetTag(ctx context.Context, options types.TagOptions) (*types.Tag, error) {
	var success types.Tag
	var failure types.ErrorResponse
	r, err := c.sling.New().Get(fmt.Sprintf("v4/projects/%d/repository/tags/%s", options.ProjectID, options.TagName)).Request()
	if err != nil {
		return nil, err
	}
	resp, err := c.sling.Do(r.WithContext(ctx), &success, &failure)
	if err != nil {
		return nil, err
	}

	switch {
	case resp.StatusCode < 300:
		return &success, nil
	}

	return nil, failure
}

package gitlab

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dghubble/sling"

	ehttp "gitlab.unanet.io/devops/go/pkg/http"
	"gitlab.unanet.io/devops/go/pkg/json"
)

const (
	userAgent = "eve-gitlab"
)

type Config struct {
	GitlabApiKey  string        `split_words:"true" required:"true"`
	GitlabBaseUrl string        `split_words:"true" required:"true"`
	GitlabTimeout time.Duration `split_words:"true" default:"20s"`
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

func (c *Client) TagCommit(ctx context.Context, options TagOptions) (*Tag, error) {
	var success Tag
	var failure ErrorResponse
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

func (c *Client) GetTag(ctx context.Context, options TagOptions) (*Tag, error) {
	var success Tag
	var failure ErrorResponse
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

func (c *Client) GetRelease(ctx context.Context, options TagOptions) (*Release, error) {
	var success Release
	var failure ErrorResponse
	r, err := c.sling.New().Get(fmt.Sprintf("v4/projects/%d/releases/%s", options.ProjectID, options.TagName)).Request()
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

func (c *Client) CreateRelease(ctx context.Context, options TagOptions) (*Release, error) {
	var success Release
	var failure ErrorResponse

	r, err := c.sling.New().
		Post(fmt.Sprintf("v4/projects/%d/releases/%s", options.ProjectID, options.TagName)).
		QueryStruct(Release{Name: options.TagName, TagName: options.TagName}).
		Request()

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

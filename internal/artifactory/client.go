package artifactory

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dghubble/sling"

	"gitlab.unanet.io/devops/eve/internal/common"
	"gitlab.unanet.io/devops/eve/pkg/slinge"
)

const (
	userAgent = "eve-artifactory"
)

type Config struct {
	ArtifactoryApiKey  string        `split_words:"true" required:"true"`
	ArtifactoryBaseUrl string        `split_words:"true" required:"true"`
	ArtifactoryTimeout time.Duration `split_words:"true" default:"20s"`
}

type Client struct {
	sling *sling.Sling
}

func NewClient(config Config) *Client {
	var httpClient = &http.Client{
		Timeout: config.ArtifactoryTimeout,
	}

	if !strings.HasSuffix(config.ArtifactoryBaseUrl, "/") {
		config.ArtifactoryBaseUrl += "/"
	}

	sling := sling.New().Base(config.ArtifactoryBaseUrl).Client(httpClient).
		Add("X-JFrog-Art-Api", config.ArtifactoryApiKey).
		Add("User-Agent", userAgent).
		ResponseDecoder(slinge.NewJsonDecoder())
	return &Client{sling: sling}
}

func (c *Client) GetLatestVersion(ctx context.Context, repository string, path string, version string) (*VersionResponse, error) {
	var success VersionResponse
	var failure ErrorResponse
	r, err := c.sling.Get(fmt.Sprintf("versions/%s/%s", repository, path)).Request()
	if err != nil {
		return nil, err
	}
	// set this here since version could have an asterisk and sling will encode the asterisk which Artifactory doesn't like
	r.URL.RawQuery = fmt.Sprintf("version=%s", version)
	resp, err := c.sling.Do(r.WithContext(ctx), &success, &failure)
	if err != nil {
		return nil, err
	}
	if http.StatusOK == resp.StatusCode {
		return &success, nil
	} else {
		return nil, &failure
	}
}

func (c *Client) MoveArtifact(ctx context.Context, srcRepo, srcPath, destRepo, destPath string, dryRun bool) (*MessagesResponse, error) {
	var success MessagesResponse
	var failure ErrorResponse
	r, err := c.sling.Post(fmt.Sprintf("move/%s/%s", srcRepo, srcPath)).Request()
	if err != nil {
		return nil, err
	}
	r.URL.RawQuery = fmt.Sprintf("to=/%s/%s&dry=%d", destRepo, destPath, common.Bool2int(dryRun))
	resp, err := c.sling.Do(r.WithContext(ctx), &success, &failure)
	if err != nil {
		return nil, err
	}

	if http.StatusOK == resp.StatusCode {
		return &success, nil
	} else {
		return nil, &failure
	}
}

func (c *Client) CopyArtifact(ctx context.Context, srcRepo, srcPath, destRepo, destPath string, dryRun bool) (*MessagesResponse, error) {
	var success MessagesResponse
	var failure ErrorResponse
	r, err := c.sling.Post(fmt.Sprintf("copy/%s/%s", srcRepo, srcPath)).Request()
	if err != nil {
		return nil, err
	}
	r.URL.RawQuery = fmt.Sprintf("to=/%s/%s&dry=%d", destRepo, destPath, common.Bool2int(dryRun))
	resp, err := c.sling.Do(r.WithContext(ctx), &success, &failure)
	if err != nil {
		return nil, err
	}

	if http.StatusOK == resp.StatusCode {
		return &success, nil
	} else {
		return nil, &failure
	}
}

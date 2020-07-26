package artifactory

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dghubble/sling"
	"go.uber.org/zap"

	"gitlab.unanet.io/devops/eve/pkg/errors"
	ehttp "gitlab.unanet.io/devops/eve/pkg/http"
	"gitlab.unanet.io/devops/eve/pkg/json"
	"gitlab.unanet.io/devops/eve/pkg/log"
)

const (
	userAgent = "eve-artifactory"
)

type Config struct {
	ArtifactoryApiKey      string        `split_words:"true" required:"true"`
	ArtifactoryApiKeyWrite string        `split_words:"true" required:"true"`
	ArtifactoryBaseUrl     string        `split_words:"true" required:"true"`
	ArtifactoryTimeout     time.Duration `split_words:"true" default:"20s"`
}

type Client struct {
	sling *sling.Sling
}

func NewClient(config Config, write bool) *Client {
	var httpClient = &http.Client{
		Timeout:   config.ArtifactoryTimeout,
		Transport: ehttp.LoggingTransport,
	}

	accessKey := ""

	if write == true {
		accessKey = config.ArtifactoryApiKeyWrite
	} else {
		accessKey = config.ArtifactoryApiKey
	}

	if !strings.HasSuffix(config.ArtifactoryBaseUrl, "/") {
		config.ArtifactoryBaseUrl += "/"
	}

	sling := sling.New().Base(config.ArtifactoryBaseUrl).Client(httpClient).
		Add("X-JFrog-Art-Api", accessKey).
		Add("User-Agent", userAgent).
		ResponseDecoder(json.NewJsonDecoder())
	return &Client{sling: sling}
}

func (c *Client) GetLatestVersion(ctx context.Context, repository string, path string, version string) (string, error) {
	var success VersionResponse
	var failure ErrorResponse
	r, err := c.sling.New().Get(fmt.Sprintf("versions/%s/%s", repository, path)).Request()
	if err != nil {
		return "", errors.Wrap(err)
	}
	// set this here since version could have an asterisk and sling will encode the asterisk which Artifactory doesn't like
	r.URL.RawQuery = fmt.Sprintf("version=%s", version)
	resp, err := c.sling.Do(r.WithContext(ctx), &success, &failure)
	if err != nil {
		return "", errors.Wrap(err)
	}

	switch resp.StatusCode {
	case http.StatusOK:
		return success.Version, nil
	case http.StatusNotFound:
		return "", NotFoundErrorf("the following Version: %s, was not found", version)
	case http.StatusServiceUnavailable:
		return "", ServiceUnavailableErrorf("Artifactory returned a 503 and appears to be unavailable")
	default:
		return "", failure
	}
}

type failRespMsg struct {
	Level   string `json:"level"`
	Message string `json:"message"`
}

type failResp struct {
	Messages []failRespMsg `json:"messages"`
}

func (f failResp) Error() string {
	msg := ""
	for _, s := range f.Messages {
		msg = msg + s.Level + ":" + s.Message + "|"
	}
	return msg
}

func (c *Client) MoveArtifact(ctx context.Context, srcRepo, srcPath, destRepo, destPath string, dryRun bool) (*MessagesResponse, error) {
	var success MessagesResponse

	var failure failResp
	r, err := c.sling.New().Post(fmt.Sprintf("move/%s/%s", srcRepo, srcPath)).Request()
	if err != nil {
		log.Logger.Debug("move artifact client req", zap.Error(err))
		return nil, err
	}
	r.URL.RawQuery = fmt.Sprintf("to=/%s/%s&dry=%d", destRepo, destPath, Bool2int(dryRun))
	log.Logger.Debug("move artifact request obj", zap.Any("req", r))
	resp, err := c.sling.Do(r.WithContext(ctx), &success, &failure)
	if err != nil {
		log.Logger.Debug("move artifact client req do", zap.Error(err))
		return nil, err
	}

	log.Logger.Debug("move artifact client req status", zap.Any("Status", resp.Status), zap.Any("StatusCode", resp.StatusCode))

	if http.StatusOK == resp.StatusCode {
		return &success, nil
	} else {
		log.Logger.Debug("move artifact client req failure", zap.Any("failure", failure))
		return nil, failure
	}
}

func (c *Client) CopyArtifact(ctx context.Context, srcRepo, srcPath, destRepo, destPath string, dryRun bool) (*MessagesResponse, error) {
	var success MessagesResponse
	var failure ErrorResponse
	r, err := c.sling.New().Post(fmt.Sprintf("copy/%s/%s", srcRepo, srcPath)).Request()
	if err != nil {
		return nil, err
	}
	r.URL.RawQuery = fmt.Sprintf("to=/%s/%s&dry=%d", destRepo, destPath, Bool2int(dryRun))
	resp, err := c.sling.Do(r.WithContext(ctx), &success, &failure)
	if err != nil {
		return nil, err
	}

	if http.StatusOK == resp.StatusCode {
		return &success, nil
	} else {
		return nil, failure
	}
}

// GetArtifactProperties for an Artifact.
func (c *Client) GetArtifactProperties(ctx context.Context, repository, path string) (*Properties, error) {
	var success Properties
	var failure string
	r, err := c.sling.New().Get(fmt.Sprintf("storage/%s/%s", repository, path)).Request() //generic-int-local/unanet/unanet/unanet-%UNANET_VERSION%.tar.gz)
	if err != nil {
		return nil, err
	}
	r.URL.RawQuery = "properties"
	resp, err := c.sling.Do(r.WithContext(ctx), &success, &failure)
	if err != nil {
		return nil, err
	}

	switch resp.StatusCode {
	case http.StatusOK:
		return &success, nil
	case http.StatusNotFound:
		return nil, NotFoundErrorf("the following artifact: %s/%s, was not found", repository, path)
	default:
		return nil, fmt.Errorf("an error occurred while trying to retrieve the artifact properites: %s", failure)
	}
}

// GetLatestVersionLessThan Retrieves the latest version of an Artifact that is is less than the one specified
// TODO: pull logic out from below and document up here
func (c *Client) GetLatestVersionLessThan(ctx context.Context, repository string, path string, lessThanVersion string) (string, error) {
	var success AQLResult
	var failure string
	var sort string
	// TODO: We should really move this logic out of here, the artifactory client shouldn't be responsible for logic like this.
	// This occurs because the path for docker includes the version due to how docker repositories work in artifactory
	// The path only includes the folder structure for a normal repository
	// Also, since it's not one file with docker but instead a tag (the version) represents a folder, there's actually many files, none of which have the version on them
	// So instead we need to sort by path descending which has the version on it.
	if strings.Contains(repository, "docker") {
		path = fmt.Sprintf("%s/*", path)
		sort = "{\"$desc\": [\"path\"]}"
	} else {
		sort = "{\"$desc\": [\"name\"]}"
	}
	aqlQuery := fmt.Sprintf("{\"$and\":[{\"repo\":{\"$eq\":\"%s\"}},{\"@version\":{\"$lt\":\"%s\"}},{\"path\":{\"$match\":\"%s\"}}]}", repository, lessThanVersion, path)
	body := strings.NewReader(fmt.Sprintf("items.find(%s).include(\"name\",\"@version\", \"path\").sort(%s).limit(1)", aqlQuery, sort))

	r, err := c.sling.New().Post("search/aql").Body(body).Request()
	resp, err := c.sling.Do(r.WithContext(ctx), &success, &failure)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", errors.Wrap(err)
	}

	if len(success.Results) == 0 {
		return "", NotFoundErrorf("no version was found less than: %s/%s:%s", repository, path, lessThanVersion)
	}

	if len(success.Results[0].Properties) == 0 {
		return "", errors.Wrap(fmt.Errorf("there is no version property for the following path: %s", success.Results[0].Path))
	}
	return success.Results[0].Properties[0].Value, nil
}

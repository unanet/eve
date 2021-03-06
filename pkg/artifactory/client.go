package artifactory

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dghubble/sling"
	"go.uber.org/zap"

	"github.com/unanet/go/pkg/errors"
	ehttp "github.com/unanet/go/pkg/http"
	"github.com/unanet/go/pkg/json"
	"github.com/unanet/go/pkg/log"
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
	cfg   Config
}

func NewClient(config Config) *Client {

	var httpClient = &http.Client{
		Timeout:   config.ArtifactoryTimeout,
		Transport: ehttp.LoggingTransport,
	}

	if !strings.HasSuffix(config.ArtifactoryBaseUrl, "/") {
		config.ArtifactoryBaseUrl += "/"
	}

	sling := sling.New().Base(config.ArtifactoryBaseUrl).Client(httpClient).
		Add("X-JFrog-Art-Api", config.ArtifactoryApiKey).
		Add("User-Agent", userAgent).
		ResponseDecoder(json.NewJsonDecoder())
	return &Client{sling: sling, cfg: config}
}

func (c *Client) GetLatestVersion(ctx context.Context, repository string, path string, version string) (string, error) {
	var success VersionResponse
	var failure ErrorResponse

	full_path := fmt.Sprintf("versions/%s/%s", repository, path)

	log.Logger.Info("get latest artifact version", zap.String("full_path", full_path))

	r, err := c.sling.New().Get(full_path).Request()
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

func (c *Client) MoveArtifact(ctx context.Context, srcRepo, srcPath, destRepo, destPath string, dryRun bool) (*MessagesResponse, error) {
	var success MessagesResponse
	var failure MessagesResponse

	r, err := c.sling.New().Post(fmt.Sprintf("move/%s/%s", srcRepo, srcPath)).Request()
	if err != nil {
		log.Logger.Error("move artifact client req failed", zap.Error(err))
		return nil, err
	}
	r.URL.RawQuery = fmt.Sprintf("to=/%s/%s&dry=%d", destRepo, destPath, Bool2int(dryRun))

	resp, err := c.sling.Do(r.WithContext(ctx), &success, &failure)
	if err != nil {
		log.Logger.Error("move artifact client req do failed", zap.Error(err))
		return nil, err
	}

	switch resp.StatusCode {
	case http.StatusOK:
		return &success, nil
	case http.StatusNotFound:
		return nil, NotFoundErrorf("the artifact was not found; source_repo: %s source_path: %s dest_repo: %s dest_path: %s", srcRepo, srcPath, destRepo, destPath)
	case http.StatusServiceUnavailable:
		return nil, ServiceUnavailableErrorf("Artifactory returned a 503 and appears to be unavailable")
	case http.StatusBadRequest:
		return nil, InvalidRequestErrorf("invalid move artifact request: %s", failure.ToString())
	default:
		log.Logger.Error("unknown artifactory response", zap.String("status", resp.Status), zap.Int("status_code", resp.StatusCode))
		return nil, failure
	}
}

func (c *Client) DeleteArtifact(ctx context.Context, repo, path string) (*MessagesResponse, error) {
	var success MessagesResponse
	var failure MessagesResponse

	artifactoryBaseURL := c.cfg.ArtifactoryBaseUrl

	// HACK: The artifactory API for delete only works on /unanet not /unanet/api
	// trimming off the api/ here for the DELETE command below
	artifactoryBaseURL = strings.Replace(artifactoryBaseURL, "api/", "", -1)

	tmpClient := c.sling.New()

	r, err := tmpClient.Base(artifactoryBaseURL).Delete(fmt.Sprintf("%s/%s", repo, path)).Request()
	if err != nil {
		return nil, err
	}
	resp, err := c.sling.Do(r.WithContext(ctx), &success, &failure)
	if err != nil {
		return nil, err
	}

	if http.StatusNoContent == resp.StatusCode {
		return &success, nil
	}
	return nil, failure

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
	}
	return nil, failure
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

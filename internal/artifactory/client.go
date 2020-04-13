package artifactory

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"
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
	client    *http.Client
	baseURL   *url.URL
	userAgent string
	apiKey    string
}

func NewClient(config Config) (*Client, error) {
	var httpClient = &http.Client{
		Timeout: config.ArtifactoryTimeout,
	}

	baseEndpoint, err := url.Parse(config.ArtifactoryBaseUrl)

	if err != nil {
		return nil, err
	}

	if !strings.HasSuffix(baseEndpoint.Path, "/") {
		baseEndpoint.Path += "/"
	}

	c := &Client{client: httpClient, baseURL: baseEndpoint, userAgent: userAgent, apiKey: config.ArtifactoryApiKey}
	return c, nil
}

func (c *Client) GetLatestVersion(ctx context.Context, repository string, path string, version string) (VersionResponse, error) {
	var vr VersionResponse
	r, err := c.request(http.MethodGet, fmt.Sprintf("versions/%s/%s?version=%s", repository, path, version), nil)
	if err != nil {
		return vr, err
	}

	_, err = c.do(ctx, r, &vr)
	if err != nil {
		return vr, err
	}

	return vr, nil
}

func (c *Client) request(method, urlStr string, body io.Reader) (*http.Request, error) {
	u, err := c.baseURL.Parse(path.Join(c.baseURL.Path, urlStr))
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-JFrog-Art-Api", c.apiKey)

	if c.userAgent != "" {
		req.Header.Set("User-Agent", c.userAgent)
	}

	return req, nil
}

func (c *Client) do(ctx context.Context, req *http.Request, v interface{}) (*http.Response, error) {
	req = req.WithContext(ctx)
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, commonErrors(ctx, err)
	}

	defer resp.Body.Close()

	if s := resp.StatusCode; 200 > s || s > 299 {
		return resp, unmarshalError(resp)
	}

	if v == nil {
		return resp, err
	}

	switch vu := v.(type) {
	case *string:
		bodyBytes, serr := ioutil.ReadAll(resp.Body)
		if serr != nil {
			return nil, serr
		}
		*vu = string(bodyBytes)
	default:
		err = json.NewDecoder(resp.Body).Decode(v)
		if err == io.EOF {
			err = nil // ignore EOF errors caused by empty response body
		}
		return resp, err
	}
	return resp, err
}

func (c *Client) CopyArtifact(ctx context.Context, sRepo string, sPath string, dRepo string, dPath string) (interface{}, interface{}) {

}

func commonErrors(ctx context.Context, err error) error {
	// If we got an error, and the context has been canceled,
	// the context's error is probably more useful.
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	if e, ok := err.(*url.Error); ok {
		if url2, uerr := url.Parse(e.URL); uerr == nil {
			e.URL = url2.String()
			return e
		}
	}

	return err
}

func unmarshalError(r *http.Response) error {
	errorResponse := &ErrorResponse{Response: r}
	data, err := ioutil.ReadAll(r.Body)
	if err == nil && data != nil {
		err = json.Unmarshal(data, errorResponse)
		if err != nil || len(errorResponse.Errors) == 0 {
			return fmt.Errorf(string(data))
		}
	}

	return errorResponse
}

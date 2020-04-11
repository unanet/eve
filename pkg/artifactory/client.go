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

	MediaTypePlain = "text/plain"
	MediaTypeXml   = "application/xml"

	MediaTypeJson = "application/json"
	MediaTypeForm = "application/x-www-form-urlencoded"
)

type Config struct {
	ArtifactoryApiKey  string        `split_words:"true" required:"true"`
	ArtifactoryBaseUrl string        `split_words:"true" required:"true"`
	ArtifactoryTimeout time.Duration `split_words:"true" default:"20s"`
}

type Client struct {
	// HTTP Client used to communicate with the API.
	client *http.Client

	// Base URL for API requests. BaseURL should always be specified with a trailing slash.
	BaseURL *url.URL

	// User agent used when communicating with the Artifactory API.
	UserAgent string

	// Api key used for authenticating with Artifactory API.
	ApiKey string
}

// NewClient creates a Client from a provided base url for an Artifactory instance and a service Client
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

	c := &Client{client: httpClient, BaseURL: baseEndpoint, UserAgent: userAgent, ApiKey: config.ArtifactoryApiKey}
	return c, nil
}

func (c *Client) GetLatestVersion(ctx context.Context, repository string, path string, version string) error {
	r, err := c.request(http.MethodGet, fmt.Sprintf("versions/%s/%s?version=%s", repository, path, version), nil)
	if err != nil {
		return err
	}

	var respString string
	resp, err := c.do(ctx, r, &respString)
	if err != nil {
		return err
	}
	fmt.Println(resp)
	//bodyBytes, err := ioutil.ReadAll(resp.Body)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//bodyString := string(bodyBytes)
	//fmt.Printf("%v", bodyString)
	return nil
}

func (c *Client) request(method, urlStr string, body io.Reader) (*http.Request, error) {
	u, err := c.BaseURL.Parse(path.Join(c.BaseURL.Path, urlStr))
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-JFrog-Art-Api", c.ApiKey)

	if c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
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

	if s := resp.StatusCode; 200 > s && s > 299 {
		return resp, unmarshalError(resp)
	}

	if v == nil {
		return resp, err
	}

	switch v.(type) {
	case string:
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		v = string(bodyBytes)
	default:
		err = json.NewDecoder(resp.Body).Decode(v)
		if err == io.EOF {
			err = nil // ignore EOF errors caused by empty response body
		}
		return resp, err
	}
	return resp, err
}

func commonErrors(ctx context.Context, err error) error {
	// If we got an error, and the context has been canceled,
	// the context's error is probably more useful.
	select {
	case <-ctx.Done():
		ctx.Err()
	default:
	}

	if e, ok := err.(*url.Error); ok {
		if url2, err := url.Parse(e.URL); err == nil {
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

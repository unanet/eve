package github

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	gohttp "net/http"
	"time"

	"github.com/unanet/eve/pkg/scm/types"
	"github.com/unanet/go/pkg/errors"
	"github.com/unanet/go/pkg/http"
	"github.com/unanet/go/pkg/log"
	"go.uber.org/zap"
)

type (
	Config struct {
		GithubAccessToken  string        `envconfig:"GITHUB_ACCESS_TOKEN"`
		GithubBaseUrl      string        `envconfig:"GITHUB_BASE_URL"`
		GithubTimeout      time.Duration `envconfig:"GITHUB_TIMEOUT" default:"20s"`
		GithubEmailAddress string        `envconfig:"GITHUB_EMAIL_ADDRESS" default:"ghost@unknown.tld"`
	}

	Client struct {
		cfg Config
		cli *gohttp.Client
	}
	ReleaseData struct {
		TagName string `json:"tag_name"`
	}

	Tagger struct {
		Name  string `json:"name"`
		Email string `json:"email"`
		Date  string `json:"date"`
	}

	TagData struct {
		Tag     string `json:"tag"`
		Object  string `json:"object"`
		Message string `json:"message"`
		Tagger  Tagger `json:"tagger"`
		Type    string `json:"type"`
	}

	RefData struct {
		Ref string `json:"ref"`
		Sha string `json:"sha"`
	}
)

func NewClient(cfg Config) *Client {
	return &Client{
		cli: &gohttp.Client{
			Transport: http.LoggingTransport,
			Timeout:   cfg.GithubTimeout,
		},
		cfg: cfg,
	}
}

func (c *Client) createTag(options types.TagOptions) error {
	var auth = fmt.Sprintf("token %s", c.cfg.GithubAccessToken)

	b, err := json.Marshal(TagData{
		Tag:     options.TagName,
		Object:  options.GitHash,
		Message: options.TagName,
		Tagger: Tagger{
			Name:  "eve",
			Email: c.cfg.GithubEmailAddress,
			Date:  time.Now().UTC().Format(time.RFC3339),
		},
		Type: "commit",
	})
	if err != nil {
		return errors.Wrap(err, "failed to marshall tag data")
	}

	url := fmt.Sprintf("%s/repos/%s/%s/git/tags", c.cfg.GithubBaseUrl, options.Owner, options.Repo)
	log.Logger.Info("github tags url", zap.String("url", url))

	req, err := gohttp.NewRequest("POST", url, bytes.NewBuffer(b))
	if err != nil {
		return errors.Wrap(err, "failed to create tag request")
	}
	req.Header.Set("Authorization", auth)

	resp, err := c.cli.Do(req)
	if err != nil {
		return errors.Wrap(err, "failed to issue tag request")
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Logger.Error("failed to close the git tag body resp", zap.Error(err))
		}
	}()
	if resp.StatusCode > 299 {
		return fmt.Errorf("failed to tag github commit: %v", resp.Status)
	}
	return nil
}

func (c *Client) createRef(options types.TagOptions) error {
	var auth = fmt.Sprintf("token %s", c.cfg.GithubAccessToken)
	bRef, err := json.Marshal(RefData{
		Ref: fmt.Sprintf("refs/tags/%s", options.TagName),
		Sha: options.GitHash,
	})
	if err != nil {
		return errors.Wrap(err, "failed to marshall tag data")
	}
	rurl := fmt.Sprintf("%s/repos/%s/%s/git/refs", c.cfg.GithubBaseUrl, options.Owner, options.Repo)

	reqRef, err := gohttp.NewRequest("POST", rurl, bytes.NewBuffer(bRef))
	if err != nil {
		return errors.Wrap(err, "failed to create tag request")
	}
	reqRef.Header.Set("Authorization", auth)

	refResp, err := c.cli.Do(reqRef)
	if err != nil {
		return errors.Wrap(err, "failed to issue tag request")
	}
	defer func() {
		if err := refResp.Body.Close(); err != nil {
			log.Logger.Error("failed to close the git ref body resp", zap.Error(err))
		}
	}()
	if refResp.StatusCode > 299 {
		return fmt.Errorf("failed to create tag ref github commit: %v", refResp.Status)
	}
	return nil
}

func (c *Client) createRelease(options types.TagOptions) error {
	var auth = fmt.Sprintf("token %s", c.cfg.GithubAccessToken)

	rurl := fmt.Sprintf("%s/repos/%s/%s/releases", c.cfg.GithubBaseUrl, options.Owner, options.Repo)

	bRef, err := json.Marshal(ReleaseData{TagName: options.TagName})
	if err != nil {
		return errors.Wrap(err, "failed to marshall release data")
	}

	reqRef, err := gohttp.NewRequest("POST", rurl, bytes.NewBuffer(bRef))
	if err != nil {
		return errors.Wrap(err, "failed to create release request")
	}
	reqRef.Header.Set("Authorization", auth)

	refResp, err := c.cli.Do(reqRef)
	if err != nil {
		return errors.Wrap(err, "failed to issue release request")
	}
	defer func() {
		if err := refResp.Body.Close(); err != nil {
			log.Logger.Error("failed to close the git release body resp", zap.Error(err))
		}
	}()
	if refResp.StatusCode > 299 {
		return fmt.Errorf("failed to create release: %v", refResp.Status)
	}
	return nil
}

func (c *Client) TagCommit(ctx context.Context, options types.TagOptions) (*types.Tag, error) {
	log.Logger.Info("tag git commit", zap.Any("opts", options))

	if err := c.createTag(options); err != nil {
		return nil, err
	}

	if err := c.createRef(options); err != nil {
		return nil, err
	}

	if err := c.createRelease(options); err != nil {
		return nil, err
	}

	return &types.Tag{
		Name: options.TagName,
		Repo: options.Repo,
	}, nil
}

func (c *Client) GetTag(ctx context.Context, options types.TagOptions) (*types.Tag, error) {
	var auth = fmt.Sprintf("token %s", c.cfg.GithubAccessToken)
	url := fmt.Sprintf("%s/repos/%s/%s/git/ref/tags/%s", c.cfg.GithubBaseUrl, options.Owner, options.Repo, options.TagName)

	req, err := gohttp.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create get tag request")
	}
	req.Header.Set("Authorization", auth)

	resp, err := c.cli.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to issue get tag request")
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Logger.Error("failed to close the get tag ref body resp", zap.Error(err))
		}
	}()
	if resp.StatusCode > 299 {
		return nil, fmt.Errorf("failed to get tag github commit: %v", resp.Status)
	}

	return &types.Tag{
		Name: options.TagName,
		Repo: options.Repo,
	}, nil
}

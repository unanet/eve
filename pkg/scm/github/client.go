package github

import (
	"context"
	"fmt"
	"time"

	gogithub "github.com/google/go-github/v38/github"
	"github.com/unanet/eve/pkg/scm/types"
	"github.com/unanet/go/pkg/http"
	"golang.org/x/oauth2"
)

const userAgent = "ava-github"

type Config struct {
	GithubAccessToken string        `envconfig:"GITHUB_ACCESS_TOKEN"`
	GithubBaseUrl     string        `envconfig:"GITHUB_BASE_URL"`
	GithubTimeout     time.Duration `envconfig:"GITHUB_TIMEOUT" default:"20s"`
}

type Client struct {
	c *gogithub.Client
}

func NewClient(cfg Config) *Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: cfg.GithubAccessToken},
	)
	tc := oauth2.NewClient(context.TODO(), ts)
	tc.Transport = http.LoggingTransport

	c := gogithub.NewClient(tc)
	c.UserAgent = userAgent

	return &Client{
		c: c,
	}
}

func (c *Client) TagCommit(ctx context.Context, options types.TagOptions) (*types.Tag, error) {
	def := time.Now().UTC()
	author := "eve"
	tag, resp, err := c.c.Git.CreateTag(ctx, options.Owner, options.Repo, &gogithub.Tag{
		Tag: &options.TagName,
		Object: &gogithub.GitObject{
			Type: "commit",
			SHA:  &options.GitHash,
		},
		SHA:     &options.GitHash,
		URL:     nil,
		Message: nil,
		Tagger: &gogithub.CommitAuthor{
			Date: &def,
			Name: &author
		},
	})
	if err != nil {
		return nil, err
	}
	if tag == nil || resp.StatusCode > 299 {
		return nil, fmt.Errorf("failed to tag github commit: %v", resp.Status)
	}
	return &types.Tag{
		Name: tag.GetTag(),
		Repo: options.Repo,
	}, nil
}

func (c *Client) GetTag(ctx context.Context, options types.TagOptions) (*types.Tag, error) {
	tag, resp, err := c.c.Git.GetTag(ctx, options.Owner, options.Repo, options.GitHash)
	if err != nil {
		return nil, err
	}
	if tag == nil || resp.StatusCode > 299 {
		return nil, fmt.Errorf("failed to get github tag: %v", resp.Status)
	}
	return &types.Tag{
		Name: tag.GetTag(),
		Repo: options.Repo,
	}, nil
}

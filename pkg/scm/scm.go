package scm

import (
	"context"
	"github.com/unanet/eve/internal/config"
	"github.com/unanet/eve/pkg/scm/types"

	"github.com/unanet/eve/pkg/scm/github"
	"github.com/unanet/eve/pkg/scm/gitlab"
	"github.com/unanet/go/pkg/log"
	"go.uber.org/zap"
)

type SourceController interface {
	TagCommit(ctx context.Context, options types.TagOptions) (*types.Tag, error)
	GetTag(ctx context.Context, options types.TagOptions) (*types.Tag, error)
}

func New() SourceController {
	cfg := config.GetConfig()
	switch cfg.SourceControlProvider {
	case "github":
		return github.NewClient(cfg.GitHubConfig)
	case "gitlab":
		return gitlab.NewClient(cfg.GitLabConfig)
	}
	log.Logger.Fatal("invalid scm provider", zap.String("scm", cfg.SourceControlProvider))
	return nil
}

package releases

import (
	"context"

	"gitlab.unanet.io/devops/eve/internal/data"
)

type ReleaseRepo interface {
	ArtifactByName(ctx context.Context, name string) (*data.RequestArtifact, error)
	FeedByName(ctx context.Context, name string) (*data.Feed, error)
}

type ReleaseSvc struct {
	repo ReleaseRepo
}

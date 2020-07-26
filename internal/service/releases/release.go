package releases

import (
	"context"
	"fmt"
	"strings"

	"go.uber.org/zap"

	"gitlab.unanet.io/devops/eve/internal/data"
	"gitlab.unanet.io/devops/eve/internal/service"
	"gitlab.unanet.io/devops/eve/internal/service/crud"
	"gitlab.unanet.io/devops/eve/pkg/artifactory"
	"gitlab.unanet.io/devops/eve/pkg/errors"
	"gitlab.unanet.io/devops/eve/pkg/eve"
	"gitlab.unanet.io/devops/eve/pkg/log"
)

type ReleaseSvc struct {
	repo              crud.Repo
	artifactoryClient *artifactory.Client
}

func NewReleaseSvc(r crud.Repo, a *artifactory.Client) *ReleaseSvc {
	return &ReleaseSvc{
		repo:              r,
		artifactoryClient: a,
	}
}

func path(group, name string) string {
	return fmt.Sprintf("%s/%s", group, name)
}

func artifactRepoPath(providerGroup, artifactName, version string) string {
	return fmt.Sprintf("%s/%s/%s", providerGroup, artifactName, version)
}

func version(version string) string {
	if version == "" {
		return "*"
	} else if len(strings.Split(version, ".")) < 4 {
		return version + ".*"
	}
	return version
}

func (svc *ReleaseSvc) moveReleaseArtifact(ctx context.Context, release eve.Release, isPromotion bool) error {
	if release.FromFeed == release.ToFeed {
		return errors.BadRequest(fmt.Sprintf("source feed: %s and destination feed: %s cannot be equal", release.FromFeed, release.ToFeed))
	}

	if strings.ToLower(release.FromFeed) == "int" && strings.ToLower(release.ToFeed) == "qa" {
		return errors.BadRequest("int and qa share the same feed so nothing to promote")
	}

	artifact, err := svc.repo.ArtifactByName(ctx, release.Artifact)
	if err != nil {
		return service.CheckForNotFoundError(err)
	}

	fromFeed, err := svc.repo.FeedByAliasAndType(ctx, release.FromFeed, artifact.FeedType)
	if err != nil {
		return service.CheckForNotFoundError(err)
	}

	artifactVersion, err := svc.artifactoryClient.GetLatestVersion(ctx, fromFeed.Name, path(artifact.ProviderGroup, artifact.Name), version(release.Version))
	if err != nil {
		log.Logger.Debug("get latest version err", zap.Error(err))
		if _, ok := err.(artifactory.NotFoundError); ok {
			return errors.NotFound(fmt.Sprintf("artifact not found in artifactory: %s/%s/%s:%s", fromFeed.Name, path(artifact.ProviderGroup, artifact.Name), artifact.Name, version(release.Version)))
		}
		return err
	}

	toFeed, err := svc.toFeed(ctx, release, artifact, fromFeed, isPromotion)
	if err != nil {
		log.Logger.Debug("toFeed err", zap.Error(err))
		return errors.Wrap(err)
	}

	if isPromotion && fromFeed.PromotionOrder > toFeed.PromotionOrder {
		return errors.BadRequestf("cannot promote from: %s to: %s", fromFeed.Name, toFeed.Name)
	}

	if isPromotion == false && toFeed.PromotionOrder > fromFeed.PromotionOrder {
		return errors.BadRequestf("cannot demote from: %s to: %s", fromFeed.Name, toFeed.Name)
	}

	fromPath := artifactRepoPath(artifact.ProviderGroup, artifact.Name, artifactVersion)
	toPath := artifactRepoPath(artifact.ProviderGroup, artifact.Name, artifactVersion)

	resp, err := svc.artifactoryClient.MoveArtifact(ctx, fmt.Sprintf("%s-local", fromFeed.Name), fromPath, fmt.Sprintf("%s-local", toFeed.Name), toPath, false)
	if err != nil {
		log.Logger.Debug("MoveArtifact err", zap.Error(err))
		return errors.Wrap(err)
	}

	log.Logger.Debug("move artifact message", zap.Any("resp", resp))
	return nil
}

func (svc *ReleaseSvc) PromoteRelease(ctx context.Context, release eve.Release) error {
	return svc.moveReleaseArtifact(ctx, release, true)
}

func (svc *ReleaseSvc) DemoteRelease(ctx context.Context, release eve.Release) error {
	return svc.moveReleaseArtifact(ctx, release, false)
}

func (svc *ReleaseSvc) toFeed(ctx context.Context, release eve.Release, artifact *data.Artifact, fromFeed *data.Feed, isPromotion bool) (*data.Feed, error) {
	if release.ToFeed != "" {
		toFeed, errr := svc.repo.FeedByAliasAndType(ctx, release.ToFeed, artifact.FeedType)
		if errr != nil {
			return nil, service.CheckForNotFoundError(errr)
		}
		return toFeed, nil
	}

	if isPromotion {
		toFeed, errr := svc.repo.NextFeedByPromotionOrderType(ctx, fromFeed.PromotionOrder, artifact.FeedType)
		if errr != nil {
			return nil, service.CheckForNotFoundError(errr)
		}
		return toFeed, nil
	} else {
		toFeed, errr := svc.repo.PreviousFeedByPromotionOrderType(ctx, fromFeed.PromotionOrder, artifact.FeedType)
		if errr != nil {
			return nil, service.CheckForNotFoundError(errr)
		}
		return toFeed, nil
	}
}
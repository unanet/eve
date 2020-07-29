package releases

import (
	"context"
	"fmt"
	"strings"

	"gitlab.unanet.io/devops/eve/internal/data"
	"gitlab.unanet.io/devops/eve/internal/service"
	"gitlab.unanet.io/devops/eve/internal/service/crud"
	"gitlab.unanet.io/devops/eve/pkg/artifactory"
	"gitlab.unanet.io/devops/eve/pkg/errors"
	"gitlab.unanet.io/devops/eve/pkg/eve"
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

func evalArtifactImageTag(a *data.Artifact, availableVersion string) string {
	imageTag := a.ImageTag
	versionSplit := strings.Split(availableVersion, ".")
	replacementMap := make(map[string]string)
	replacementMap["$version"] = availableVersion
	for i, x := range versionSplit {
		replacementMap[fmt.Sprintf("$%d", i+1)] = x
	}
	for k, v := range replacementMap {
		imageTag = strings.Replace(imageTag, k, v, -1)
	}
	return imageTag
}

func (svc *ReleaseSvc) Release(ctx context.Context, release eve.Release) (eve.Release, error) {
	success := eve.Release{}
	if release.FromFeed == release.ToFeed {
		return success, errors.BadRequest(fmt.Sprintf("source feed: %s and destination feed: %s cannot be equal", release.FromFeed, release.ToFeed))
	}

	if strings.ToLower(release.FromFeed) == "int" && strings.ToLower(release.ToFeed) == "qa" {
		return success, errors.BadRequest("int and qa share the same feed so nothing to promote")
	}

	artifact, err := svc.repo.ArtifactByName(ctx, release.Artifact)
	if err != nil {
		return success, service.CheckForNotFoundError(err)
	}

	fromFeed, err := svc.repo.FeedByAliasAndType(ctx, release.FromFeed, artifact.FeedType)
	if err != nil {
		return success, service.CheckForNotFoundError(err)
	}

	artifactVersion, err := svc.artifactoryClient.GetLatestVersion(ctx, fromFeed.Name, path(artifact.ProviderGroup, artifact.Name), version(release.Version))
	if err != nil {
		if _, ok := err.(artifactory.NotFoundError); ok {
			return success, errors.NotFound(fmt.Sprintf("artifact not found in artifactory: %s/%s/%s:%s", fromFeed.Name, path(artifact.ProviderGroup, artifact.Name), artifact.Name, version(release.Version)))
		}
		return success, err
	}

	toFeed, err := svc.toFeed(ctx, release, artifact, fromFeed)
	if err != nil {
		return success, errors.Wrap(err)
	}

	fromPath := artifactRepoPath(artifact.ProviderGroup, artifact.Name, evalArtifactImageTag(artifact, artifactVersion))
	toPath := artifactRepoPath(artifact.ProviderGroup, artifact.Name, evalArtifactImageTag(artifact, artifactVersion))

	// HACK: Delete the destination first
	// Artifactory fails when copy/move an artifact to a location that already exists
	_, _ = svc.artifactoryClient.DeleteArtifact(ctx, fmt.Sprintf("%s-local", toFeed.Name), toPath)

	resp, err := svc.artifactoryClient.MoveArtifact(ctx, fmt.Sprintf("%s-local", fromFeed.Name), fromPath, fmt.Sprintf("%s-local", toFeed.Name), toPath, false)
	if err != nil {
		return success, errors.Wrap(err)
	}

	success.Artifact = artifact.Name
	success.Version = artifactVersion
	success.ToFeed = toFeed.Alias
	success.FromFeed = fromFeed.Alias
	success.Message = resp.ToString()

	return success, nil
}

func (svc *ReleaseSvc) toFeed(ctx context.Context, release eve.Release, artifact *data.Artifact, fromFeed *data.Feed) (*data.Feed, error) {
	if release.ToFeed != "" {
		toFeed, errr := svc.repo.FeedByAliasAndType(ctx, release.ToFeed, artifact.FeedType)
		if errr != nil {
			return nil, service.CheckForNotFoundError(errr)
		}
		return toFeed, nil
	}

	toFeed, errr := svc.repo.NextFeedByPromotionOrderType(ctx, fromFeed.PromotionOrder, artifact.FeedType)
	if errr != nil {
		return nil, service.CheckForNotFoundError(errr)
	}
	return toFeed, nil

}

package releases

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	goerrors "github.com/pkg/errors"
	"gitlab.unanet.io/devops/go/pkg/errors"
	"gitlab.unanet.io/devops/go/pkg/log"
	"go.uber.org/zap"

	"gitlab.unanet.io/devops/eve/internal/data"
	"gitlab.unanet.io/devops/eve/internal/service"
	"gitlab.unanet.io/devops/eve/pkg/artifactory"
	"gitlab.unanet.io/devops/eve/pkg/eve"
	"gitlab.unanet.io/devops/eve/pkg/gitlab"
)

type ReleaseSvc struct {
	repo              *data.Repo
	artifactoryClient *artifactory.Client
	gitlabClient      *gitlab.Client
}

func NewReleaseSvc(r *data.Repo, a *artifactory.Client, g *gitlab.Client) *ReleaseSvc {
	return &ReleaseSvc{
		repo:              r,
		artifactoryClient: a,
		gitlabClient:      g,
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
		return success, goerrors.Wrapf(err, "failed to get the latest artifact version")
	}

	toFeed, err := svc.toFeed(ctx, release, artifact, fromFeed)
	if err != nil {
		return success, goerrors.Wrapf(err, "failed to get the artifact destination (to) feed")
	}

	log.Logger.Info("release artifact",
		zap.String("artifact", artifact.Name),
		zap.String("version", artifactVersion),
		zap.String("from_feed", fromFeed.Name),
		zap.String("to_feed", toFeed.Name),
	)

	fromPath := artifactRepoPath(artifact.ProviderGroup, artifact.Name, evalArtifactImageTag(artifact, artifactVersion))
	toPath := artifactRepoPath(artifact.ProviderGroup, artifact.Name, evalArtifactImageTag(artifact, artifactVersion))

	// HACK: Delete the destination first
	// Artifactory fails when copy/move an artifact to a location that already exists
	_, _ = svc.artifactoryClient.DeleteArtifact(ctx, fmt.Sprintf("%s-local", toFeed.Name), toPath)

	fromRepo := fmt.Sprintf("%s-local", fromFeed.Name)
	toRepo := fmt.Sprintf("%s-local", toFeed.Name)

	resp, err := svc.artifactoryClient.CopyArtifact(ctx, fromRepo, fromPath, toRepo, toPath, false)
	if err != nil {
		if _, ok := err.(artifactory.NotFoundError); ok {
			return success, errors.NotFound(fmt.Sprintf("artifact not found: %s", err.Error()))
		}
		if _, ok := err.(artifactory.InvalidRequestError); ok {
			return success, errors.BadRequest(fmt.Sprintf("invalid artifact request: %s", err.Error()))
		}
		return success, goerrors.Wrapf(err, "failed to move the artifact from: %s to: %s", fromPath, toPath)
	}

	success.Artifact = artifact.Name
	success.Version = artifactVersion
	success.ToFeed = toFeed.Alias
	success.FromFeed = fromFeed.Alias
	success.Message = resp.ToString()

	// If we are releasing to prod we tag the commit in GitLab
	if strings.ToLower(toFeed.Alias) == "prod" {
		artifactProps, perr := svc.artifactoryClient.GetArtifactProperties(ctx, toRepo, toPath)
		if perr != nil {
			if _, ok := err.(artifactory.NotFoundError); ok {
				return success, errors.NotFound(fmt.Sprintf("artifact not found: %s", perr.Error()))
			}
			return success, errors.Wrap(perr)
		}

		gitBranch := artifactProps.Property("gitlab-build-properties.git-branch")
		gitSHA := artifactProps.Property("gitlab-build-properties.git-sha")
		gitProjectID := artifactProps.Property("gitlab-build-properties.project-id")
		fullVersion := artifactProps.Property("version")
		releaseVersion := parseVersion(fullVersion)

		if releaseVersion == "v" || releaseVersion == "" {
			return success, errors.BadRequestf("invalid version: %v", releaseVersion)
		}

		log.Logger.Info("artifact release to prod",
			zap.String("branch", gitBranch),
			zap.String("sha", gitSHA),
			zap.String("project", gitProjectID),
			zap.String("full_version", fullVersion),
			zap.String("release_version", fullVersion),
		)

		projectID, cErr := strconv.Atoi(gitProjectID)
		if cErr != nil {
			return success, errors.Wrap(cErr)
		}

		tOpts := gitlab.TagOptions{
			ProjectID: projectID,
			TagName:   releaseVersion,
			GitHash:   gitSHA,
		}

		tag, _ := svc.gitlabClient.GetTag(ctx, tOpts)
		if tag != nil && tag.Name != "" {
			return success, errors.BadRequestf("the version: %v has already been tagged", releaseVersion)
		}

		_, gErr := svc.gitlabClient.TagCommit(ctx, tOpts)
		if gErr != nil {
			return success, goerrors.Wrapf(gErr, "failed to tag the gitlab commit")
		}

		rel, _ := svc.gitlabClient.GetRelease(ctx, tOpts)
		if rel != nil && rel.Name != "" {
			return success, errors.BadRequestf("the version: %v has already been released", releaseVersion)
		}

		_, rErr := svc.gitlabClient.CreateRelease(ctx, tOpts)
		if rErr != nil {
			return success, goerrors.Wrapf(rErr, "failed to create gitlab release")
		}
		log.Logger.Info("artifact released",
			zap.String("branch", gitBranch),
			zap.String("sha", gitSHA),
			zap.String("project", gitProjectID),
			zap.String("full_version", fullVersion),
			zap.String("release_version", fullVersion),
		)
	}

	return success, nil
}

func parseVersion(fullVersion string) string {
	v := ""
	vParts := strings.Split(fullVersion, ".")
	switch len(vParts) {
	case 1, 2, 3:
		v = fullVersion
	default:
		v = strings.Join(vParts[0:3], ".")
	}
	if strings.HasPrefix(fullVersion, "v") {
		return v
	}
	return fmt.Sprintf("v%s", v)
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

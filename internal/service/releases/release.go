package releases

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	goerrors "github.com/pkg/errors"
	"github.com/unanet/go/pkg/errors"
	"github.com/unanet/go/pkg/log"
	"go.uber.org/zap"

	"github.com/unanet/eve/internal/config"
	"github.com/unanet/eve/internal/data"
	"github.com/unanet/eve/internal/service"
	"github.com/unanet/eve/internal/service/crud"
	"github.com/unanet/eve/pkg/artifactory"
	"github.com/unanet/eve/pkg/eve"
	"github.com/unanet/eve/pkg/scm"
	"github.com/unanet/eve/pkg/scm/types"
)

type ReleaseSvc struct {
	repo              *data.Repo
	artifactoryClient *artifactory.Client
	scm               scm.SourceController
	crud              *crud.Manager
}

func NewReleaseSvc(
	r *data.Repo,
	a *artifactory.Client,
	g scm.SourceController,
	crud *crud.Manager) *ReleaseSvc {
	return &ReleaseSvc{
		repo:              r,
		artifactoryClient: a,
		scm:               g,
		crud:              crud,
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

type artifactReleaseInfo struct {
	GitBranch, GitSHA, BuildVersion, ReleaseVersion string
	FromPath, ToPath                                string
	FromRepo, ToRepo                                string
	ProjectName                                     string
	FromFeed, ToFeed                                *data.Feed
	Artifact                                        *data.Artifact
	ProjectID                                       int
}

func (svc *ReleaseSvc) releaseInfo(ctx context.Context, release eve.Release) (*artifactReleaseInfo, error) {
	artifact, err := svc.repo.ArtifactByName(ctx, release.Artifact)
	if err != nil {
		return nil, service.CheckForNotFoundError(err)
	}

	fromFeed, err := svc.repo.FeedByAliasAndType(ctx, release.FromFeed, artifact.FeedType)
	if err != nil {
		return nil, service.CheckForNotFoundError(err)
	}

	toFeed, err := svc.toFeed(ctx, release, artifact, fromFeed)
	if err != nil {
		return nil, goerrors.Wrapf(err, "failed to get the artifact destination (to) feed")
	}

	artifactVersion, err := svc.artifactoryClient.GetLatestVersion(ctx, fromFeed.Name, path(artifact.ProviderGroup, artifact.Name), version(release.Version))
	if err != nil {
		if _, ok := err.(artifactory.NotFoundError); ok {
			return nil, errors.NotFound(fmt.Sprintf("artifact not found in artifactory: %s/%s/%s:%s", fromFeed.Name, path(artifact.ProviderGroup, artifact.Name), artifact.Name, version(release.Version)))
		}
		return nil, goerrors.Wrapf(err, "failed to get the latest artifact version")
	}

	fromPath := artifactRepoPath(artifact.ProviderGroup, artifact.Name, evalArtifactImageTag(artifact, artifactVersion))
	toPath := artifactRepoPath(artifact.ProviderGroup, artifact.Name, evalArtifactImageTag(artifact, artifactVersion))

	fromRepo := fmt.Sprintf("%s-local", fromFeed.Name)
	toRepo := fmt.Sprintf("%s-local", toFeed.Name)

	artifactProps, perr := svc.artifactoryClient.GetArtifactProperties(ctx, fromRepo, fromPath)
	if perr != nil {
		if _, ok := err.(artifactory.NotFoundError); ok {
			return nil, errors.NotFound(fmt.Sprintf("artifact not found: %s", perr.Error()))
		}
		return nil, errors.Wrap(perr)
	}

	var (
		projectID   int
		projectName string
	)

	var (
		scmId              = config.BuildPropertyID()
		projectIDBuildProp = fmt.Sprintf("%s-build-properties.project-id", scmId)
		gitBranchBuildProp = fmt.Sprintf("%s-build-properties.git-branch", scmId)
		gitShaBuildProp    = fmt.Sprintf("%s-build-properties.git-sha", scmId)
	)

	if projectID, err = strconv.Atoi(artifactProps.Property(projectIDBuildProp)); err != nil {
		projectName = artifactProps.Property(projectIDBuildProp)
	}

	relInfo := artifactReleaseInfo{
		GitBranch:      artifactProps.Property(gitBranchBuildProp),
		GitSHA:         artifactProps.Property(gitShaBuildProp),
		BuildVersion:   artifactProps.Property("version"),
		ReleaseVersion: parseVersion(artifactProps.Property("version")),
		FromPath:       fromPath,
		ToPath:         toPath,
		FromRepo:       fromRepo,
		ToRepo:         toRepo,
		ToFeed:         toFeed,
		FromFeed:       fromFeed,
		Artifact:       artifact,
		ProjectID:      projectID,
		ProjectName:    projectName,
	}

	log.Logger.Info("release artifact info", zap.Any("release_info", relInfo))

	return &relInfo, nil

}

func (svc *ReleaseSvc) Release(ctx context.Context, release eve.Release) ([]eve.Release, error) {
	var success []eve.Release

	if release.FromFeed == release.ToFeed {
		return success, errors.BadRequest(fmt.Sprintf("source feed: %s and destination feed: %s cannot be equal", release.FromFeed, release.ToFeed))
	}

	if strings.ToLower(release.FromFeed) == "int" && strings.ToLower(release.ToFeed) == "qa" {
		return success, errors.BadRequest("int and qa share the same feed so nothing to release")
	}

	var rel eve.Release
	var err error

	var isArtifactRelease bool
	switch release.Type {
		case eve.ReleaseTypeArtifact:
			rel, err = svc.releaseArtifact(ctx, release)
			isArtifactRelease = true
		case eve.ReleaseTypeNamespace:
			success, err = svc.releaseNamespace(ctx, release)
	}

	if err != nil {
		return success, errors.BadRequestf("unable to release %s", release.Type)
	}

	if isArtifactRelease {
		return []eve.Release{rel}, nil
	}

	return success, nil
}

func (svc *ReleaseSvc) releaseArtifact(ctx context.Context, release eve.Release) (eve.Release, error) {
	success := eve.Release{}

	relInfo, err := svc.releaseInfo(ctx, release)
	if err != nil {
		return success, goerrors.Wrapf(err, "failed to get the release info")
	}

	if relInfo.ReleaseVersion == "v" || relInfo.ReleaseVersion == "" {
		return success, errors.BadRequestf("invalid version: %v", relInfo.ReleaseVersion)
	}

	gitTagOpts := types.TagOptions{
		ProjectID: relInfo.ProjectID,
		TagName:   fmt.Sprintf("%v/%v", relInfo.Artifact.Name, relInfo.ReleaseVersion),
		GitHash:   relInfo.GitSHA,
	}

	if len(relInfo.ProjectName) > 0 && len(strings.Split(relInfo.ProjectName, "/")) == 2 {
		gitTagOpts.Owner = strings.Split(relInfo.ProjectName, "/")[0]
		gitTagOpts.Repo = strings.Split(relInfo.ProjectName, "/")[1]
	}

	// Check if tag already exists
	tag, _ := svc.scm.GetTag(ctx, gitTagOpts)
	if tag != nil && tag.Name != "" {
		return success, errors.BadRequestf("the version: %v has already been tagged", tag.Name)
	}

	// Delete the destination first
	// Cant move/copy to a location that already exists
	_, _ = svc.artifactoryClient.DeleteArtifact(ctx, relInfo.ToRepo, relInfo.ToPath)

	resp, err := svc.artifactoryClient.CopyArtifact(ctx, relInfo.FromRepo, relInfo.FromPath, relInfo.ToRepo, relInfo.ToPath, false)
	if err != nil {
		if _, ok := err.(artifactory.NotFoundError); ok {
			return success, errors.NotFound(fmt.Sprintf("artifact not found: %s", err.Error()))
		}
		if _, ok := err.(artifactory.InvalidRequestError); ok {
			return success, errors.BadRequest(fmt.Sprintf("invalid artifact request: %s", err.Error()))
		}
		return success, goerrors.Wrapf(err, "failed to move the artifact from: %s to: %s", relInfo.FromPath, relInfo.ToPath)
	}

	// If we are releasing to prod we tag the commit in GitLab
	if strings.ToLower(relInfo.ToFeed.Alias) == "prod" {
		_, gErr := svc.scm.TagCommit(ctx, gitTagOpts)
		if gErr != nil {
			return success, goerrors.Wrapf(gErr, "failed to tag the commit")
		}
	}

	success.Artifact = relInfo.Artifact.Name
	success.Version = relInfo.ReleaseVersion
	success.ToFeed = relInfo.ToFeed.Alias
	success.FromFeed = relInfo.FromFeed.Alias
	success.Message = resp.ToString()

	log.Logger.Info("artifact released", zap.Any("result", success))

	return success, nil
}

func (svc *ReleaseSvc) releaseNamespace(ctx context.Context, release eve.Release) ([]eve.Release, error) {
	var releases []eve.Release

	eveServices, err := svc.crud.ServicesByNamespace(ctx, fmt.Sprintf("%s-%s", release.Environment, release.Namespace))
	if err != nil {
		return releases, err
	}

	var errs []error
	for _, eveService := range eveServices {
		release.Artifact = eveService.ArtifactName
		release.Version = eveService.DeployedVersion

		rel, err := svc.releaseArtifact(ctx, release)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		releases = append(releases, rel)
	}

	if len(errs) != 0 {
		return releases, fmt.Errorf("unable to release all artifacts. %v", errs)
	}

	return releases, nil
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

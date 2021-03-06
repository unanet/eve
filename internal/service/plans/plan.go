package plans

import (
	"context"
	"sort"
	"strconv"
	"strings"

	uuid "github.com/satori/go.uuid"
	"github.com/unanet/go/pkg/errors"
	"github.com/unanet/go/pkg/json"
	"github.com/unanet/go/pkg/log"
	"go.uber.org/zap"

	"github.com/unanet/eve/internal/data"
	"github.com/unanet/eve/pkg/artifactory"
	"github.com/unanet/eve/pkg/eve"
	"github.com/unanet/eve/pkg/queue"
)

type VersionQuery interface {
	GetLatestVersion(ctx context.Context, repository string, path string, version string) (string, error)
}

type PlanGenerator struct {
	repo *data.Repo
	vq   VersionQuery
	q    QWriter
}

func NewPlanGenerator(r *data.Repo, v VersionQuery, q QWriter) *PlanGenerator {
	return &PlanGenerator{
		repo: r,
		vq:   v,
		q:    q,
	}
}

func (d *PlanGenerator) QueuePlan(ctx context.Context, options *eve.DeploymentPlanOptions) error {
	// make sure the environment name is valid
	env, err := d.repo.EnvironmentByName(ctx, options.Environment)
	if err != nil {
		if _, ok := err.(data.NotFoundError); ok {
			return errors.NotFoundf("environment: %s, not found", options.Environment)
		}
		return errors.Wrap(err)
	}

	options.DeploymentIDs = []uuid.UUID{}

	// whether they explicitly supplied artifacts or whether they were generated
	artifactsSupplied := len(options.Artifacts) > 0

	namespaceRequests, err := d.validateNamespaces(ctx, env, options)
	if err != nil {
		return errors.Wrap(err)
	}

	err = d.validateArtifactDefinitions(ctx, env, options, namespaceRequests)
	if err != nil {
		return errors.Wrap(err)
	}

	err = d.setArtifactoryVersions(ctx, options)
	if err != nil {
		return errors.Wrap(err)
	}

	// nothing to do, should exit
	if len(options.Artifacts) == 0 {
		return errors.NewRestError(400, "no artifacts would be deployed: %v", options.Messages)
	}

	for _, ns := range namespaceRequests {
		nsPlanOptions, marshalErr := json.StructToJsonObject(&eve.NamespacePlanOptions{
			NamespaceRequest:  ns,
			ArtifactsSupplied: artifactsSupplied,
			Artifacts:         options.Artifacts,
			ForceDeploy:       options.ForceDeploy,
			DryRun:            options.DryRun,
			CallbackURL:       options.CallbackURL,
			EnvironmentID:     env.ID,
			EnvironmentName:   env.Name,
			EnvironmentAlias:  env.Alias,
			Type:              options.Type,
			Metadata:          options.Metadata,
		})
		if marshalErr != nil {
			return errors.Wrap(marshalErr)
		}
		dataDeployment := data.Deployment{
			EnvironmentID: env.ID,
			NamespaceID:   ns.ID,
			ReqID:         log.GetReqID(ctx),
			PlanOptions:   nsPlanOptions,
			User:          options.User,
		}
		repoErr := d.repo.CreateDeployment(ctx, &dataDeployment)
		if repoErr != nil {
			return errors.Wrap(repoErr)
		}
		options.DeploymentIDs = append(options.DeploymentIDs, dataDeployment.ID)
		queueM := queue.M{
			ID:      dataDeployment.ID,
			GroupID: ns.GetQueueGroupID(),
			Command: queue.CommandScheduleDeployment,
		}
		if qErr := d.q.Message(ctx, &queueM); qErr != nil {
			return errors.Wrap(qErr)
		}
		repoErr = d.repo.UpdateDeploymentMessageID(ctx, queueM.ID, queueM.MessageID)
		if repoErr != nil {
			return errors.Wrap(repoErr)
		}
	}
	return nil
}

func (d *PlanGenerator) validateArtifactDefinitions(ctx context.Context, env *data.Environment, options *eve.DeploymentPlanOptions, ns eve.NamespaceRequests) error {
	// If services were supplied, we check those against the database to make sure they are valid and pull
	// required info needed to lookup in Artifactory
	// It's important to note here that we're matching on the service/database name that's configured in the database which can be different than the artifact name
	if len(options.Artifacts) > 0 {
		var artifacts eve.ArtifactDefinitions

		for _, x := range options.Artifacts {
			var ras data.RequestArtifacts
			var err error
			switch options.Type {
			case eve.DeploymentPlanTypeApplication, eve.DeploymentPlanTypeRestart:
				ras, err = d.repo.RequestServiceArtifactByEnvironment(ctx, x.Name, x.ArtifactName, env.ID, ns.ToIDs())
			case eve.DeploymentPlanTypeJob:
				ras, err = d.repo.RequestJobArtifactByEnvironment(ctx, x.Name, x.ArtifactName, env.ID, ns.ToIDs())
			}
			if err != nil {
				if _, ok := err.(data.NotFoundError); ok {
					return errors.NotFoundf("service/job not found: %s", x.Name)
				}
				return errors.Wrap(err)
			}

			for _, y := range ras {
				artifacts = append(artifacts, x.Clone(func(d *eve.ArtifactDefinition) {
					d.ID = y.ArtifactID
					d.ArtifactName = y.ArtifactName
					d.ArtifactoryFeed = y.FeedName
					d.ArtifactoryPath = y.Path()
					d.FeedType = y.FeedType
					// we're defaulting to the namespace/service version that's configured if it's not specified and
					// the data query returns it, it should be noted, that the data query only returns the requested version
					// when the namespace count is 1.
					if len(y.RequestedVersion) > 0 && len(d.RequestedVersion) == 0 {
						d.RequestedVersion = y.RequestedVersion
					}
				}))
			}

			options.Artifacts = artifacts
		}
	} else {
		// If no services were supplied, we get all services for the supplied namespaces
		var dataArtifacts data.RequestArtifacts
		var err error
		switch options.Type {
		case eve.DeploymentPlanTypeApplication, eve.DeploymentPlanTypeRestart:
			dataArtifacts, err = d.repo.ServiceArtifacts(ctx, ns.ToIDs())
		case eve.DeploymentPlanTypeJob:
			dataArtifacts, err = d.repo.JobArtifacts(ctx, ns.ToIDs())
		}

		if err != nil {
			return errors.Wrap(err)
		}
		for _, x := range dataArtifacts {
			options.Artifacts = append(options.Artifacts, &eve.ArtifactDefinition{
				ID:               x.ArtifactID,
				ArtifactName:     x.ArtifactName,
				RequestedVersion: x.RequestedVersion,
				ArtifactoryFeed:  x.FeedName,
				ArtifactoryPath:  x.Path(),
				FeedType:         x.FeedType,
			})
		}
	}

	return nil
}

func (d *PlanGenerator) validateNamespaces(ctx context.Context, env *data.Environment, options *eve.DeploymentPlanOptions) (eve.NamespaceRequests, error) {
	// lets start with all the namespaces in the Env and filter it down based on additional information passed in.
	namespacesToDeploy, err := d.repo.NamespacesByEnvironmentID(ctx, env.ID)
	if err != nil {
		return nil, err
	}
	if len(namespacesToDeploy) == 0 {
		return nil, errors.NewRestError(400, "no associated namespaces in %s", env.Name)
	}
	if len(options.NamespaceAliases) > 0 {
		// Make sure that the namespaces that are specified are also available in the environment
		included, _ := namespacesToDeploy.FilterNamespaces(func(namespace data.Namespace) bool {
			return options.NamespaceAliases.Contains(namespace.Alias)
		})
		for _, x := range options.NamespaceAliases {
			if !included.Contains(x) {
				return nil, errors.NewRestError(400, "invalid namespace: %s", x)
			}
		}
		namespacesToDeploy = included
	} else {
		// If we didn't specify any namespaces, we need to make sure were not deploying to namespaces that require you to explicitly specify them
		included, excluded := namespacesToDeploy.FilterNamespaces(func(namespace data.Namespace) bool {
			return !namespace.ExplicitDeploy
		})
		namespacesToDeploy = included
		for _, x := range excluded {
			options.Message("explicit namespace excluded: %s", x.Alias)
		}
	}

	options.NamespaceAliases = namespacesToDeploy.ToAliases()
	var namespaceRequests eve.NamespaceRequests
	for _, x := range namespacesToDeploy {
		namespaceRequests = append(namespaceRequests, &eve.NamespaceRequest{
			ID:        x.ID,
			Name:      x.Name,
			Alias:     x.Alias,
			ClusterID: x.ClusterID,
			Version:   x.RequestedVersion,
		})
	}

	return namespaceRequests, nil
}

func (d *PlanGenerator) setArtifactoryVersions(ctx context.Context, options *eve.DeploymentPlanOptions) error {
	// now we query artifactory for the actual version
	var artifacts eve.ArtifactDefinitions
	for _, a := range options.Artifacts {
		// if you didn't pass a full version, we need to add a wildcard so it work correctly to query artifactory
		log.Logger.Info("get artifact",
			zap.String("feed", a.ArtifactoryFeed),
			zap.String("path", a.ArtifactoryPath),
			zap.String("version", a.ArtifactoryRequestedVersion()),
		)
		version, err := d.vq.GetLatestVersion(ctx, a.ArtifactoryFeed, a.ArtifactoryPath, a.ArtifactoryRequestedVersion())
		if err != nil {
			if _, ok := err.(artifactory.NotFoundError); ok {
				options.Message("artifact not found in artifactory: %s/%s/%s:%s", a.ArtifactoryFeed, a.ArtifactoryPath, a.Name, a.ArtifactoryRequestedVersion())
				continue
			}
			return errors.Wrap(err)
		}

		a.RequestedVersion = ""
		a.AvailableVersion = version
		artifacts = append(artifacts, a)
	}

	// we need to sort the higher versions first so that when we match, it tries to match the highest version possible first
	sort.Slice(artifacts, func(i, j int) bool {
		jSplit := strings.Split(artifacts[j].AvailableVersion, ".")
		iSplit := strings.Split(artifacts[i].AvailableVersion, ".")
		minLength := min(len(jSplit), len(iSplit))
		for x := 0; x < minLength; x++ {
			jv, err := strconv.Atoi(jSplit[x])
			if err != nil {
				return false
			}
			iv, err := strconv.Atoi(iSplit[x])
			if err != nil {
				return false
			}
			if iv == jv {
				continue
			}
			return jv < iv
		}
		return true
	})
	options.Artifacts = artifacts
	return nil
}

func min(x, y int) int {
	if x > y {
		return y
	}
	return x
}

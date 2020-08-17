package deployments

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	uuid "github.com/satori/go.uuid"

	"gitlab.unanet.io/devops/eve/internal/data"
	"gitlab.unanet.io/devops/eve/pkg/artifactory"
	"gitlab.unanet.io/devops/eve/pkg/errors"
	"gitlab.unanet.io/devops/eve/pkg/eve"
	"gitlab.unanet.io/devops/eve/pkg/json"
	"gitlab.unanet.io/devops/eve/pkg/middleware"
	"gitlab.unanet.io/devops/eve/pkg/queue"
)

type PlanRepo interface {
	EnvironmentByName(ctx context.Context, name string) (*data.Environment, error)
	NamespacesByEnvironmentID(ctx context.Context, environmentID int) (data.Namespaces, error)
	ServiceArtifacts(ctx context.Context, namespaceIDs []int) (data.RequestArtifacts, error)
	DatabaseInstanceArtifacts(ctx context.Context, namespaceIDs []int) (data.RequestArtifacts, error)
	RequestServiceArtifactByEnvironment(ctx context.Context, artifactName string, environmentID int) (*data.RequestArtifact, error)
	RequestDatabaseArtifactByEnvironment(ctx context.Context, databaseName string, environmentID int) (*data.RequestArtifact, error)
	CreateDeployment(ctx context.Context, d *data.Deployment) error
	UpdateDeploymentMessageID(ctx context.Context, id uuid.UUID, messageID string) error
}

type StringList []string

func (s StringList) Contains(value string) bool {
	for _, a := range s {
		if a == value {
			return true
		}
	}
	return false
}

type VersionQuery interface {
	GetLatestVersion(ctx context.Context, repository string, path string, version string) (string, error)
}

type PlanType string

const (
	DeploymentPlanTypeApplication PlanType = "application"
	DeploymentPlanTypeMigration   PlanType = "migration"
)

type ArtifactDefinition struct {
	ID               int    `json:"id"`
	Name             string `json:"name"`
	ArtifactName     string `json:"artifact_name"`
	RequestedVersion string `json:"requested_version,omitempty"`
	AvailableVersion string `json:"available_version"`
	ArtifactoryFeed  string `json:"artifactory_feed"`
	ArtifactoryPath  string `json:"artifactory_path"`
	FunctionPointer  string `json:"function_pointer"`
	FeedType         string `json:"feed_type"`
	Matched          bool   `json:"-"`
}

func (ad ArtifactDefinition) ArtifactoryRequestedVersion() string {
	if ad.RequestedVersion == "" {
		return "*"
	} else if len(strings.Split(ad.RequestedVersion, ".")) < 4 {
		return ad.RequestedVersion + ".*"
	}
	return ad.RequestedVersion
}

type ArtifactDefinitions []*ArtifactDefinition

func (ad ArtifactDefinitions) ContainsVersion(name string, version string) bool {
	for _, x := range ad {
		if x.AvailableVersion == version && x.ArtifactName == name {
			return true
		}
	}
	return false
}

func (ad ArtifactDefinitions) Match(artifactID int, optName string, requestedVersion string) *ArtifactDefinition {
	for _, x := range ad {
		if x.Name != "" {
			if x.Name == optName && strings.HasPrefix(x.AvailableVersion, requestedVersion) {
				return x
			}
		} else if x.ID == artifactID && strings.HasPrefix(x.AvailableVersion, requestedVersion) {
			return x
		}
	}
	return nil
}

func (ad ArtifactDefinitions) UnMatched() ArtifactDefinitions {
	var unmatched ArtifactDefinitions
	for _, x := range ad {
		if !x.Matched {
			unmatched = append(unmatched, x)
		}
	}
	return unmatched
}

type PlanOptions struct {
	Artifacts        ArtifactDefinitions `json:"artifacts"`
	ForceDeploy      bool                `json:"force_deploy"`
	User             string              `json:"user"`
	DryRun           bool                `json:"dry_run"`
	CallbackURL      string              `json:"callback_url"`
	Environment      string              `json:"environment"`
	NamespaceAliases StringList          `json:"namespaces,omitempty"`
	Messages         []string            `json:"messages,omitempty"`
	Type             PlanType            `json:"type"`
	DeploymentIDs    []uuid.UUID         `json:"deployment_ids,omitempty"`
}

type NamespacePlanOptions struct {
	NamespaceRequest  *eve.NamespaceRequest `json:"namespace"`
	Artifacts         ArtifactDefinitions   `json:"artifacts"`
	ArtifactsSupplied bool                  `json:"artifacts_supplied"`
	ForceDeploy       bool                  `json:"force_deploy"`
	DryRun            bool                  `json:"dry_run"`
	CallbackURL       string                `json:"callback_url"`
	EnvironmentID     int                   `json:"environment_id"`
	EnvironmentName   string                `json:"environment_name"`
	EnvironmentAlias  string                `json:"environment_alias"`
	Type              PlanType              `json:"type"`
}

func (po *PlanOptions) Message(format string, a ...interface{}) {
	po.Messages = append(po.Messages, fmt.Sprintf(format, a...))
}

func (po PlanOptions) HasArtifacts() bool {
	return len(po.Artifacts) > 0
}

func (po PlanOptions) HasNamespaceAliases() bool {
	return len(po.NamespaceAliases) > 0
}

func (po PlanOptions) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, &po,
		validation.Field(&po.Environment, validation.Required),
		validation.Field(&po.Type, validation.Required, validation.In(DeploymentPlanTypeApplication, DeploymentPlanTypeMigration)),
		validation.Field(&po.User, validation.Required))
}

type PlanGenerator struct {
	repo PlanRepo
	vq   VersionQuery
	q    QWriter
}

func NewPlanGenerator(r PlanRepo, v VersionQuery, q QWriter) *PlanGenerator {
	return &PlanGenerator{
		repo: r,
		vq:   v,
		q:    q,
	}
}

func (d *PlanGenerator) QueuePlan(ctx context.Context, options *PlanOptions) error {
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

	// Business Rule: CANNOT explicitly deploy artifacts to more than 1 Namespace
	if artifactsSupplied && len(namespaceRequests) > 1 {
		return errors.BadRequestf("cannot explicitly deploy artifacts: %v to more than one namespace: %v", options.Artifacts, options.NamespaceAliases)
	}

	err = d.validateArtifactDefinitions(ctx, options, namespaceRequests)
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
		nsPlanOptions, marshalErr := json.StructToJson(&NamespacePlanOptions{
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
		})
		if marshalErr != nil {
			return errors.Wrap(marshalErr)
		}
		dataDeployment := data.Deployment{
			EnvironmentID: env.ID,
			NamespaceID:   ns.ID,
			ReqID:         middleware.GetReqID(ctx),
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
			ReqID:   middleware.GetReqID(ctx),
			Command: CommandScheduleDeployment,
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

func (d *PlanGenerator) validateArtifactDefinitions(ctx context.Context, options *PlanOptions, nsr eve.NamespaceRequests) error {
	// If services were supplied, we check those against the database to make sure they are valid and pull
	// required info needed to lookup in Artifactory
	// It's important to note here that we're matching on the service/database name that's configured in the database which can be different than the artifact name.

	if len(options.Artifacts) > 0 && len(nsr) > 1 {
		return errors.BadRequestf("cannot explicitly deploy artifacts: %v to more than one namespace: %v", options.Artifacts, nsr)
	}

	// Explicitly deploy these artifacts to this ONE namespace
	if len(options.Artifacts) > 0 && len(nsr) == 1 {
		ns := nsr[0]
		for _, x := range options.Artifacts {
			var ra *data.RequestArtifact
			var err error
			if options.Type == DeploymentPlanTypeApplication {
				ra, err = d.repo.RequestServiceArtifactByEnvironment(ctx, x.Name, ns.ID)
			} else {
				ra, err = d.repo.RequestDatabaseArtifactByEnvironment(ctx, x.Name, ns.ID)
			}
			if err != nil {
				if _, ok := err.(data.NotFoundError); ok {
					return errors.NotFoundf("service/database not found in db: %s", x.Name)
				}
				return errors.Wrap(err)
			}
			x.ID = ra.ArtifactID
			x.ArtifactName = ra.ArtifactName
			x.ArtifactoryFeed = ra.FeedName
			x.ArtifactoryPath = ra.Path()
			x.FunctionPointer = ra.FunctionPointer.String
			x.FeedType = ra.FeedType
			if len(x.RequestedVersion) == 0 {
				x.RequestedVersion = ra.RequestedVersion
			}
		}
	} else {
		// If no services were supplied, we get all services for the supplied namespaces
		var dataArtifacts data.RequestArtifacts
		var err error
		if options.Type == DeploymentPlanTypeApplication {
			dataArtifacts, err = d.repo.ServiceArtifacts(ctx, nsr.ToIDs())
		} else {
			dataArtifacts, err = d.repo.DatabaseInstanceArtifacts(ctx, nsr.ToIDs())
		}
		if err != nil {
			return errors.Wrap(err)
		}
		for _, x := range dataArtifacts {
			options.Artifacts = append(options.Artifacts, &ArtifactDefinition{
				ID:               x.ArtifactID,
				ArtifactName:     x.ArtifactName,
				RequestedVersion: x.RequestedVersion,
				ArtifactoryFeed:  x.FeedName,
				ArtifactoryPath:  x.Path(),
				FunctionPointer:  x.FunctionPointer.String,
				FeedType:         x.FeedType,
			})
		}
	}

	return nil
}

func (d *PlanGenerator) validateNamespaces(ctx context.Context, env *data.Environment, options *PlanOptions) (eve.NamespaceRequests, error) {
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
			return !namespace.ExplicitDeployOnly
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
		})
	}

	return namespaceRequests, nil
}

func (d *PlanGenerator) setArtifactoryVersions(ctx context.Context, options *PlanOptions) error {
	// now we query artifactory for the actual version
	var artifacts ArtifactDefinitions
	for _, a := range options.Artifacts {
		// if you didn't pass a full version, we need to add a wildcard so it work correctly to query artifactory
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

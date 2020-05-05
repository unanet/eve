package service

import (
	"context"
	"database/sql/driver"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/go-chi/chi/middleware"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"

	"gitlab.unanet.io/devops/eve/internal/data"
	"gitlab.unanet.io/devops/eve/pkg/artifactory"
	"gitlab.unanet.io/devops/eve/pkg/errors"
	"gitlab.unanet.io/devops/eve/pkg/eve"
	"gitlab.unanet.io/devops/eve/pkg/json"
	"gitlab.unanet.io/devops/eve/pkg/log"
	"gitlab.unanet.io/devops/eve/pkg/queue"
)

type DeploymentPlanRepo interface {
	EnvironmentByName(ctx context.Context, name string) (*data.Environment, error)
	NamespacesByEnvironmentID(ctx context.Context, environmentID int) (data.Namespaces, error)
	ServiceArtifacts(ctx context.Context, namespaceIDs []int) (data.RequestArtifacts, error)
	DatabaseInstanceArtifacts(ctx context.Context, namespaceIDs []int) (data.RequestArtifacts, error)
	RequestArtifactByEnvironment(ctx context.Context, artifactName string, environmentID int) (*data.RequestArtifact, error)
	CreateDeploymentTx(ctx context.Context, d *data.Deployment) (driver.Tx, error)
	UpdateDeploymentMessageIDTx(ctx context.Context, tx driver.Tx, id uuid.UUID, messageID string) error
}

type VersionQuery interface {
	GetLatestVersion(ctx context.Context, repository string, path string, version string) (string, error)
}

type DeploymentPlanType string

const (
	DeploymentPlanTypeApplication DeploymentPlanType = "application"
	DeploymentPlanTypeMigration   DeploymentPlanType = "migration"
)

type ArtifactDefinition struct {
	ID               int    `json:"id"`
	Name             string `json:"name"`
	RequestedVersion string `json:"requested_version,omitempty"`
	AvailableVersion string `json:"available_version"`
	ArtifactoryFeed  string `json:"artifactory_feed"`
	ArtifactoryPath  string `json:"artifactory_path"`
	FunctionPointer  string `json:"function_pointer"`
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
		if x.AvailableVersion == version && x.Name == name {
			return true
		}
	}
	return false
}

func (ad ArtifactDefinitions) Match(artifactID int, requestedVersion string) *ArtifactDefinition {
	for _, x := range ad {
		if x.ID == artifactID && strings.HasPrefix(x.AvailableVersion, requestedVersion) {
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

type DeploymentPlanOptions struct {
	Artifacts        ArtifactDefinitions `json:"artifacts"`
	ForceDeploy      bool                `json:"force_deploy"`
	User             string              `json:"user"`
	DryRun           bool                `json:"dry_run"`
	CallbackURL      string              `json:"callback_url"`
	Environment      string              `json:"environment"`
	NamespaceAliases StringList          `json:"namespaces,omitempty"`
	Messages         []string            `json:"messages,omitempty"`
	Type             DeploymentPlanType  `json:"type"`
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
	Type              DeploymentPlanType    `json:"type"`
}

func (po *DeploymentPlanOptions) Message(format string, a ...interface{}) {
	po.Messages = append(po.Messages, fmt.Sprintf(format, a...))
}

func (po DeploymentPlanOptions) HasArtifacts() bool {
	return len(po.Artifacts) > 0
}

func (po DeploymentPlanOptions) HasNamespaceAliases() bool {
	return len(po.NamespaceAliases) > 0
}

func (po DeploymentPlanOptions) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, &po,
		validation.Field(&po.Environment, validation.Required),
		validation.Field(&po.Type, validation.Required, validation.In(DeploymentPlanTypeApplication, DeploymentPlanTypeMigration)),
		validation.Field(&po.User, validation.Required))
}

type DeploymentPlanGenerator struct {
	repo DeploymentPlanRepo
	vq   VersionQuery
	q    QWriter
}

func NewDeploymentPlanGenerator(r DeploymentPlanRepo, v VersionQuery, q QWriter) *DeploymentPlanGenerator {
	return &DeploymentPlanGenerator{
		repo: r,
		vq:   v,
		q:    q,
	}
}

func (d *DeploymentPlanGenerator) QueueDeploymentPlan(ctx context.Context, options *DeploymentPlanOptions) error {
	// make sure the environment name is valid and get the metadata
	env, err := d.repo.EnvironmentByName(ctx, options.Environment)
	if err != nil {
		if _, ok := err.(data.NotFoundError); ok {
			return errors.NotFoundf("environment: %s, not found", options.Environment)
		}
		return errors.Wrap(err)
	}

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

	// nothing to do, should exit
	if len(options.Artifacts) == 0 {
		return errors.NewRestError(400, "no artifacts would be deployed: %v", options.Messages)
	}

	for _, ns := range namespaceRequests {
		nsPlanOptions, err := json.StructToJson(&NamespacePlanOptions{
			NamespaceRequest:  ns,
			ArtifactsSupplied: artifactsSupplied,
			Artifacts:         options.Artifacts,
			ForceDeploy:       options.ForceDeploy,
			DryRun:            options.DryRun,
			CallbackURL:       options.CallbackURL,
			EnvironmentID:     env.ID,
			EnvironmentName:   env.Name,
			Type:              options.Type,
		})
		if err != nil {
			return errors.Wrap(err)
		}
		dataDeployment := data.Deployment{
			EnvironmentID: env.ID,
			NamespaceID:   ns.ID,
			ReqID:         middleware.GetReqID(ctx),
			PlanOptions:   nsPlanOptions,
			User:          options.User,
		}
		tx, err := d.repo.CreateDeploymentTx(ctx, &dataDeployment)
		if err != nil {
			return errors.WrapTx(tx, err)
		}
		log.Logger.Debug("created deploymentID", zap.String("id", dataDeployment.ID.String()))
		queueM := queue.M{
			ID:      dataDeployment.ID,
			GroupID: ns.GetQueueGroupID(),
			ReqID:   middleware.GetReqID(ctx),
			Command: CommandScheduleDeployment,
		}
		if err := d.q.Message(ctx, &queueM); err != nil {
			return errors.WrapTx(tx, err)
		}
		err = d.repo.UpdateDeploymentMessageIDTx(ctx, tx, queueM.ID, queueM.MessageID)
		if err != nil {
			return errors.WrapTx(tx, err)
		}
		err = tx.Commit()
		if err != nil {
			return errors.Wrap(err)
		}
	}
	return nil
}

func (d *DeploymentPlanGenerator) validateArtifactDefinitions(ctx context.Context, env *data.Environment, options *DeploymentPlanOptions, ns eve.NamespaceRequests) error {
	// If services were supplied, we check those against the database to make sure they are valid and pull
	// required info needed to lookup in Artifactory
	if len(options.Artifacts) > 0 {
		for _, x := range options.Artifacts {
			ra, err := d.repo.RequestArtifactByEnvironment(ctx, x.Name, env.ID)
			if err != nil {
				if _, ok := err.(data.NotFoundError); ok {
					return errors.NotFoundf("artifact not found in db: %s", x.Name)
				}
				return errors.Wrap(err)
			}
			x.ArtifactoryFeed = ra.FeedName
			x.ArtifactoryPath = ra.Path()
			x.ID = ra.ArtifactID
		}
	} else {
		// If no services were supplied, we get all services for the supplied namespaces
		var dataArtifacts data.RequestArtifacts
		var err error
		if options.Type == DeploymentPlanTypeApplication {
			dataArtifacts, err = d.repo.ServiceArtifacts(ctx, ns.ToIDs())
		} else {
			dataArtifacts, err = d.repo.DatabaseInstanceArtifacts(ctx, ns.ToIDs())
		}
		if err != nil {
			return errors.Wrap(err)
		}
		for _, x := range dataArtifacts {
			options.Artifacts = append(options.Artifacts, &ArtifactDefinition{
				ID:               x.ArtifactID,
				Name:             x.ArtifactName,
				RequestedVersion: x.RequestedVersion,
				ArtifactoryFeed:  x.FeedName,
				ArtifactoryPath:  x.Path(),
				FunctionPointer:  x.FunctionPointer,
			})
		}
	}

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

		// if this version is already in the list, don't include it again
		if artifacts.ContainsVersion(a.Name, version) {
			continue
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

func (d *DeploymentPlanGenerator) validateNamespaces(ctx context.Context, env *data.Environment, options *DeploymentPlanOptions) (eve.NamespaceRequests, error) {
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

func min(x, y int) int {
	if x > y {
		return y
	}
	return x
}

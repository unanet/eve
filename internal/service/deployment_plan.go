package service

import (
	"context"
	"fmt"
	"strings"

	"gitlab.unanet.io/devops/eve/internal/data"
	"gitlab.unanet.io/devops/eve/pkg/artifactory"
	"gitlab.unanet.io/devops/eve/pkg/errors"
)

type DeploymentPlanRepo interface {
	NamespacesByEnvironmentName(context.Context, string) (data.Namespaces, error)
	ServicesByNamespaceID(ctx context.Context, namespaceID int) (data.Services, error)
	EnvironmentByName(ctx context.Context, name string) (*data.Environment, error)
	RequestedArtifacts(ctx context.Context, namespaceIDs []int) (data.RequestedArtifacts, error)
	RequestedArtifactByEnvironment(ctx context.Context, artifactName string, environmentID int) (*data.RequestedArtifact, error)
}

type VersionQuery interface {
	GetLatestVersion(ctx context.Context, repository string, path string, version string) (string, error)
}

type RequestedArtifact struct {
	ArtifactID           int    `json:"-"`
	ArtifactName         string `json:"name"`
	ProviderGroup        string `json:"provider_group"`
	FeedName             string `json:"feed_name"`
	RequestedVersion     string `json:"requested_version"`
	VersionInArtifactory string `json:"version_in_artifactory"`
}

func (ra RequestedArtifact) Path() string {
	return fmt.Sprintf("%s/%s", ra.ProviderGroup, ra.ArtifactName)
}

type RequestedArtifacts []RequestedArtifact

func fromServiceDefinition(s ServiceDefinition) RequestedArtifact {
	return RequestedArtifact{
		ArtifactName:     s.Name(),
		RequestedVersion: s.Version(),
	}
}

func fromServiceDefinitions(s ServiceDefinitions) RequestedArtifacts {
	var returnList []RequestedArtifact
	for _, x := range s {
		returnList = append(returnList, fromServiceDefinition(x))
	}
	return returnList
}

func fromDataRequestedArtifact(s data.RequestedArtifact) RequestedArtifact {
	return RequestedArtifact{
		ArtifactID:       s.ArtifactID,
		ArtifactName:     s.ArtifactName,
		ProviderGroup:    s.ProviderGroup,
		FeedName:         s.FeedName,
		RequestedVersion: s.RequestedVersion,
	}
}

func fromDataRequestedArtifacts(s data.RequestedArtifacts) RequestedArtifacts {
	var returnList []RequestedArtifact
	for _, x := range s {
		returnList = append(returnList, fromDataRequestedArtifact(x))
	}
	return returnList
}

type Namespace struct {
	ID       int                    `json:"id"`
	Name     string                 `json:"name"`
	Metadata map[string]interface{} `json:"-"`
}

type Namespaces []Namespace

func (ns Namespaces) IDs() []int {
	var ids []int
	for _, n := range ns {
		ids = append(ids, n.ID)
	}
	return ids
}

func fromDataNamespace(ns data.Namespace) Namespace {
	return Namespace{
		ID:       ns.ID,
		Name:     ns.Name,
		Metadata: ns.Metadata.AsMap(),
	}
}

func fromDataNamespaces(ns data.Namespaces) Namespaces {
	var returnList []Namespace
	for _, x := range ns {
		returnList = append(returnList, fromDataNamespace(x))
	}
	return returnList
}

type DeploymentPlan struct {
	Environment         string                 `json:"environment"`
	EnvironmentMetadata map[string]interface{} `json:"-"`
	Messages            StringList             `json:"messages"`
	Namespaces          Namespaces             `json:"namespaces"`
	RequestedArtifacts  RequestedArtifacts     `json:"requested_artifacts"`
	EnvironmentID       int                    `json:"-"`
}

func (dp *DeploymentPlan) message(format string, a ...interface{}) {
	dp.Messages = append(dp.Messages, fmt.Sprintf(format, a...))
}

type ServiceDefinition string

func (s ServiceDefinition) Name() string {
	return strings.Split(string(s), ":")[0]
}

func (s ServiceDefinition) Version() string {
	split := strings.Split(string(s), ":")
	if len(split) > 1 {
		return split[1]
	}
	return ""
}

type ServiceDefinitions []ServiceDefinition

type DeploymentPlanOptions struct {
	Environment string
	Namespaces  StringList
	Services    ServiceDefinitions
}

type DeploymentPlanGenerator struct {
	repo DeploymentPlanRepo
	vq   VersionQuery
	Plan *DeploymentPlan
}

func NewDeploymentPlanGenerator(r DeploymentPlanRepo, v VersionQuery) *DeploymentPlanGenerator {
	return &DeploymentPlanGenerator{
		repo: r,
		vq:   v,
		Plan: &DeploymentPlan{},
	}
}

func (d *DeploymentPlanGenerator) Generate(ctx context.Context, options DeploymentPlanOptions) error {
	if err := d.setEnvironment(ctx, options); err != nil {
		return err
	}
	if err := d.generateNamespaces(ctx, options); err != nil {
		return err
	}
	if err := d.generateServices(ctx, options); err != nil {
		return err
	}
	return nil
}

func (d *DeploymentPlanGenerator) generateServices(ctx context.Context, options DeploymentPlanOptions) error {
	// If services were supplied, we check those against the database to make sure they are valid and pull
	// required info needed to lookup in Artifactory
	if len(options.Services) > 0 {
		var requestedArtifacts RequestedArtifacts
		for _, x := range fromServiceDefinitions(options.Services) {
			artifact, err := d.repo.RequestedArtifactByEnvironment(ctx, x.ArtifactName, d.Plan.EnvironmentID)
			if err != nil {
				if _, ok := err.(data.NotFoundError); ok {
					return errors.NotFoundf("artifact not found in db: %s, environment: %s, ", x.ArtifactName, d.Plan.Environment)
				}
				return errors.WrapUnexpected(err)
			}
			artifact.RequestedVersion = x.RequestedVersion
			requestedArtifacts = append(requestedArtifacts, fromDataRequestedArtifact(*artifact))
			d.Plan.RequestedArtifacts = requestedArtifacts
		}
	} else {
		// If no services were supplied, we get all services for the supplied namespaces
		dataArtifacts, err := d.repo.RequestedArtifacts(ctx, d.Plan.Namespaces.IDs())
		if err != nil {
			return errors.WrapUnexpected(err)
		}
		d.Plan.RequestedArtifacts = fromDataRequestedArtifacts(dataArtifacts)
	}

	// now we query artifactory for the actual version
	for _, a := range d.Plan.RequestedArtifacts {
		// if you didn't pass a full version, we need to add a wildcard so it work correctly to query artifactory
		if a.RequestedVersion == "" {
			a.RequestedVersion = "*"
		} else if len(strings.Split(a.RequestedVersion, ".")) < 4 {
			a.RequestedVersion = a.RequestedVersion + ".*"
		}
		version, err := d.vq.GetLatestVersion(ctx, a.FeedName, a.Path(), a.RequestedVersion)
		if err != nil {
			if _, ok := err.(artifactory.NotFoundError); ok {
				return errors.NotFoundf("artifact not found in artifactory: %s, version: %s", a.ArtifactName, a.RequestedVersion)
			}
			return errors.WrapUnexpected(err)
		}
		a.VersionInArtifactory = version
	}
	return nil
}

func (d *DeploymentPlanGenerator) setEnvironment(ctx context.Context, options DeploymentPlanOptions) error {
	// make sure the environment name is valid and get the metadata
	env, err := d.repo.EnvironmentByName(ctx, options.Environment)
	if err != nil {
		if _, ok := err.(data.NotFoundError); ok {
			return errors.NotFoundf("environment: %s, not found", options.Environment)
		}
		return errors.WrapUnexpected(err)
	}
	d.Plan.EnvironmentMetadata = env.Metadata.AsMap()
	d.Plan.Environment = env.Name
	d.Plan.EnvironmentID = env.ID
	return nil
}

func (d *DeploymentPlanGenerator) generateNamespaces(ctx context.Context, options DeploymentPlanOptions) error {
	// lets start with all the namespaces in the Env and filter it down based on additional information passed in.
	namespacesToDeploy, err := d.repo.NamespacesByEnvironmentName(ctx, options.Environment)
	if err != nil {
		return err
	}
	if len(namespacesToDeploy) == 0 {
		// We have no namespaces in the env specified,so we effectively can't do anything.
		d.Plan.message("the following environment: %s, has no associated namespaces", options.Environment)
		return nil
	}
	if len(options.Namespaces) > 0 {
		// Make sure that the namespaces that are specified are also available in the environment
		included, excluded := filterNamespaces(namespacesToDeploy, func(namespace data.Namespace) bool {
			return options.Namespaces.Contains(namespace.Name)
		})
		namespacesToDeploy = included
		if len(excluded) > 0 {
			d.Plan.message("the following namespaces were supplied: %v, but are not setup in %s", excluded, options.Environment)
		}
	} else {
		// If we didn't specify any namespaces, we need to make sure were not deploying to namespaces that require you to explicitly specify them
		included, excluded := filterNamespaces(namespacesToDeploy, func(namespace data.Namespace) bool {
			return !namespace.ExplicitDeployOnly
		})
		namespacesToDeploy = included
		if len(excluded) > 0 {
			d.Plan.message("since no namespaces were supplied, the following namespaces are excluded as you must explicitly specify them: %v", excluded)
		}
	}

	d.Plan.Namespaces = fromDataNamespaces(namespacesToDeploy)
	return nil
}

func filterNamespaces(ns data.Namespaces, filter func(namespace data.Namespace) bool) (data.Namespaces, data.Namespaces) {
	var included data.Namespaces
	var excluded data.Namespaces
	for _, x := range ns {
		if filter(x) {
			included = append(included, x)
		} else {
			excluded = append(excluded, x)
		}
	}

	return included, excluded
}

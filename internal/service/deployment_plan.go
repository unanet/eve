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
	EnvironmentByName(ctx context.Context, name string) (*data.Environment, error)
	NamespacesByEnvironmentID(ctx context.Context, environmentID int) (data.Namespaces, error)
	DeployedServicesByNamespaceIDs(ctx context.Context, namespaceIDs []interface{}) (data.DeployedArtifacts, error)
	DeployedDatabaseInstancesByNamespaceIDs(ctx context.Context, namespaceIDs []interface{}) (data.DeployedArtifacts, error)
	ServiceArtifacts(ctx context.Context, namespaceIDs []interface{}) (data.RequestArtifacts, error)
	DatabaseInstanceArtifacts(ctx context.Context, namespaceIDs []interface{}) (data.RequestArtifacts, error)
	RequestArtifactByEnvironment(ctx context.Context, artifactName string, environmentID int) (*data.RequestArtifact, error)
}

type VersionQuery interface {
	GetLatestVersion(ctx context.Context, repository string, path string, version string) (string, error)
}

type RequestedArtifact struct {
	ArtifactID           int                    `json:"-"`
	ArtifactMetadata     map[string]interface{} `json:"-"`
	ArtifactName         string                 `json:"name"`
	ProviderGroup        string                 `json:"provider_group"`
	FeedName             string                 `json:"feed_name"`
	RequestedVersion     string                 `json:"requested_version"`
	VersionInArtifactory string                 `json:"version_in_artifactory"`
	matched              bool
}

func (ra RequestedArtifact) Path() string {
	return fmt.Sprintf("%s/%s", ra.ProviderGroup, ra.ArtifactName)
}

func (ra RequestedArtifact) ArtifactoryRequestedVersion() string {
	if ra.RequestedVersion == "" {
		return "*"
	} else if len(strings.Split(ra.RequestedVersion, ".")) < 4 {
		return ra.RequestedVersion + ".*"
	}
	return ra.RequestedVersion
}

type RequestedArtifacts []*RequestedArtifact

func (ra RequestedArtifacts) ID(id int) *RequestedArtifact {
	for _, x := range ra {
		if x.ArtifactID == id {
			return x
		}
	}
	empty := RequestedArtifact{}
	return &empty
}

func (ra RequestedArtifacts) Match(requestedVersion string, artifactID int) *RequestedArtifact {
	for _, x := range ra {
		if x.ArtifactID == artifactID && strings.HasPrefix(x.RequestedVersion, requestedVersion) {
			return x
		}
	}
	return nil
}

func (ra RequestedArtifacts) UnMatched() StringList {
	var unmatched []string
	for _, x := range ra {
		if !x.matched {
			unmatched = append(unmatched, fmt.Sprintf("%s:%s", x.ArtifactName, x.ArtifactoryRequestedVersion()))
		}
	}
	return unmatched
}

func fromArtifactDefinition(s ArtifactDefinition) *RequestedArtifact {
	return &RequestedArtifact{
		ArtifactName:     s.Name(),
		RequestedVersion: s.Version(),
	}
}

func fromArtifactDefinitions(s ArtifactDefinitions) RequestedArtifacts {
	var returnList []*RequestedArtifact
	for _, x := range s {
		returnList = append(returnList, fromArtifactDefinition(x))
	}
	return returnList
}

func fromDataRequestedArtifact(s data.RequestArtifact) *RequestedArtifact {
	return &RequestedArtifact{
		ArtifactID:       s.ArtifactID,
		ArtifactName:     s.ArtifactName,
		ProviderGroup:    s.ProviderGroup,
		FeedName:         s.FeedName,
		ArtifactMetadata: s.ArtifactMetadata.AsMap(),
		RequestedVersion: s.RequestedVersion,
	}
}

func fromDataRequestedArtifacts(s data.RequestArtifacts) RequestedArtifacts {
	var returnList []*RequestedArtifact
	for _, x := range s {
		returnList = append(returnList, fromDataRequestedArtifact(x))
	}
	return returnList
}

type DeployedArtifact struct {
	ID               int    `json:"-"`
	NamespaceID      int    `json:"-"`
	NamespaceName    string `json:"-"`
	ArtifactID       int    `json:"-"`
	ArtifactName     string `json:"artifact_name"`
	RequestedVersion string `json:"requested_version"`
	DeployedVersion  string `json:"deployed_version"`
	AvailableVersion string `json:"available_version"`
	CustomerName     string `json:"customer_name,omitempty"`
	ServerMetadata   M      `json:"-"`
	Metadata         M      `json:"-"`
}

type DeployedArtifacts []*DeployedArtifact

func fromDataDeployedArtifact(da data.DeployedArtifact) *DeployedArtifact {
	return &DeployedArtifact{
		ID:               da.ID,
		NamespaceID:      da.NamespaceID,
		NamespaceName:    da.NamespaceName,
		ArtifactID:       da.ArtifactID,
		ArtifactName:     da.ArtifactName,
		RequestedVersion: da.RequestedVersion,
		DeployedVersion:  da.DeployedVersion.String,
		ServerMetadata:   da.ServerMetadata.AsMap(),
		Metadata:         da.Metadata.AsMap(),
		CustomerName:     da.CustomerName.String,
	}
}

func fromDataDeployedArtifacts(services data.DeployedArtifacts) DeployedArtifacts {
	var artifactsList DeployedArtifacts
	for _, x := range services {
		artifactsList = append(artifactsList, fromDataDeployedArtifact(x))
	}
	return artifactsList
}

type Environment struct {
	ID       int                    `json:"id"`
	Name     string                 `json:"name"`
	Metadata map[string]interface{} `json:"-"`
}

func fromDataEnvironment(environment data.Environment) Environment {
	return Environment{
		ID:       environment.ID,
		Name:     environment.Name,
		Metadata: environment.Metadata.AsMap(),
	}
}

type Namespace struct {
	ID        int                    `json:"-"`
	Name      string                 `json:"-"`
	Alias     string                 `json:"namespace"`
	Artifacts DeployedArtifacts      `json:"artifacts,omitempty"`
	Metadata  map[string]interface{} `json:"-"`
}

type Namespaces []*Namespace

func (ns Namespaces) IDs() []interface{} {
	var ids []interface{}
	for _, n := range ns {
		ids = append(ids, n.ID)
	}
	return ids
}

func (ns Namespaces) ID(id int) *Namespace {
	for _, x := range ns {
		if x.ID == id {
			return x
		}
	}
	return nil
}

func fromDataNamespace(ns data.Namespace) *Namespace {
	return &Namespace{
		ID:       ns.ID,
		Name:     ns.Name,
		Alias:    ns.Alias,
		Metadata: ns.Metadata.AsMap(),
	}
}

func fromDataNamespaces(ns data.Namespaces) Namespaces {
	var returnList []*Namespace
	for _, x := range ns {
		returnList = append(returnList, fromDataNamespace(x))
	}
	return returnList
}

type DeploymentPlan struct {
	Environment        *Environment `json:"environment"`
	Messages           StringList   `json:"messages"`
	Namespaces         Namespaces   `json:"plan"`
	requestedArtifacts RequestedArtifacts
}

type messageLogger func(format string, a ...interface{})

func (dp *DeploymentPlan) message(format string, a ...interface{}) {
	dp.Messages = append(dp.Messages, fmt.Sprintf(format, a...))
}

func (dp *DeploymentPlan) getMetadata(namespaceID int, artifactID int) M {
	metadata := mergeKeys(dp.Environment.Metadata, dp.Namespaces.ID(namespaceID).Metadata)
	metadata = mergeKeys(metadata, dp.requestedArtifacts.ID(artifactID).ArtifactMetadata)
	return metadata
}

type ArtifactDefinition string

func (s ArtifactDefinition) Name() string {
	return strings.Split(string(s), ":")[0]
}

func (s ArtifactDefinition) Version() string {
	split := strings.Split(string(s), ":")
	if len(split) > 1 {
		return split[1]
	}
	return ""
}

type ArtifactDefinitions []ArtifactDefinition

type DeploymentPlanGenerator struct {
	repo DeploymentPlanRepo
	vq   VersionQuery
}

type DeploymentPlanOptions struct {
	Environment      string
	NamespaceAliases StringList
	Artifacts        ArtifactDefinitions
	ForceDeploy      bool
	DryRun           bool
}

func (po DeploymentPlanOptions) HasArtifacts() bool {
	return len(po.Artifacts) > 0
}

func (po DeploymentPlanOptions) HasNamespaceAliases() bool {
	return len(po.NamespaceAliases) > 0
}

func NewDeploymentPlanGenerator(r DeploymentPlanRepo, v VersionQuery) *DeploymentPlanGenerator {
	return &DeploymentPlanGenerator{
		repo: r,
		vq:   v,
	}
}

func (d *DeploymentPlanGenerator) setupDeploymentPlan(ctx context.Context, options DeploymentPlanOptions, artifactsFunc requestedArtifactsFunc) (*DeploymentPlan, error) {
	dp := DeploymentPlan{}
	// make sure the environment name is valid and get the metadata
	dataEnv, err := d.repo.EnvironmentByName(ctx, options.Environment)
	if err != nil {
		if _, ok := err.(data.NotFoundError); ok {
			return nil, errors.NotFoundf("environment: %s, not found", options.Environment)
		}
		return nil, errors.Wrap(err)
	}
	environment := fromDataEnvironment(*dataEnv)
	dp.Environment = &environment

	namespaces, err := d.getNamespaces(ctx, &environment, options.NamespaceAliases, dp.message)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	dp.Namespaces = namespaces

	requestedArtifacts, err := d.getRequestedArtifacts(ctx, &environment, options.Artifacts, namespaces.IDs(), dp.message, artifactsFunc)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	dp.requestedArtifacts = requestedArtifacts
	return &dp, nil
}

func (d *DeploymentPlanGenerator) GenerateMigrationPlan(ctx context.Context, options DeploymentPlanOptions) (*DeploymentPlan, error) {
	plan, err := d.setupDeploymentPlan(ctx, options, d.repo.DatabaseInstanceArtifacts)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	databaseInstances, err := d.repo.DeployedDatabaseInstancesByNamespaceIDs(ctx, plan.Namespaces.IDs())
	if err != nil {
		return nil, errors.Wrap(err)
	}

	artifacts := fromDataDeployedArtifacts(databaseInstances)

	err = d.getArtifactsToDeploy(ctx, options, artifacts, plan)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	return plan, nil
}

func (d *DeploymentPlanGenerator) GenerateApplicationPlan(ctx context.Context, options DeploymentPlanOptions) (*DeploymentPlan, error) {
	plan, err := d.setupDeploymentPlan(ctx, options, d.repo.ServiceArtifacts)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	dataServices, err := d.repo.DeployedServicesByNamespaceIDs(ctx, plan.Namespaces.IDs())
	if err != nil {
		return nil, errors.Wrap(err)
	}
	artifacts := fromDataDeployedArtifacts(dataServices)

	err = d.getArtifactsToDeploy(ctx, options, artifacts, plan)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	return plan, nil
}

type requestedArtifactsFunc func(context.Context, []interface{}) (data.RequestArtifacts, error)

func (d *DeploymentPlanGenerator) getRequestedArtifacts(ctx context.Context, environment *Environment, services ArtifactDefinitions,
	namespaceIDs []interface{}, logger messageLogger, artifactsFunc requestedArtifactsFunc) (RequestedArtifacts, error) {
	// If services were supplied, we check those against the database to make sure they are valid and pull
	// required info needed to lookup in Artifactory
	var requestedArtifacts RequestedArtifacts
	if len(services) > 0 {
		for _, x := range fromArtifactDefinitions(services) {
			artifact, err := d.repo.RequestArtifactByEnvironment(ctx, x.ArtifactName, environment.ID)
			if err != nil {
				if _, ok := err.(data.NotFoundError); ok {
					return nil, errors.NotFoundf("artifact not found in db: %s", x.ArtifactName)
				}
				return nil, errors.Wrap(err)
			}
			artifact.RequestedVersion = x.RequestedVersion
			requestedArtifacts = append(requestedArtifacts, fromDataRequestedArtifact(*artifact))
		}
	} else {
		// If no services were supplied, we get all services for the supplied namespaces
		dataArtifacts, err := artifactsFunc(ctx, namespaceIDs)
		if err != nil {
			return nil, errors.Wrap(err)
		}
		requestedArtifacts = fromDataRequestedArtifacts(dataArtifacts)
	}

	// now we query artifactory for the actual version
	for _, a := range requestedArtifacts {
		// if you didn't pass a full version, we need to add a wildcard so it work correctly to query artifactory
		version, err := d.vq.GetLatestVersion(ctx, a.FeedName, a.Path(), a.ArtifactoryRequestedVersion())
		if err != nil {
			if _, ok := err.(artifactory.NotFoundError); ok {
				logger("artifact not found in artifactory: %s/%s/%s:%s", a.FeedName, a.Path(), a.ArtifactName, a.ArtifactoryRequestedVersion())
				continue
			}
			return nil, errors.Wrap(err)
		}
		a.VersionInArtifactory = version
	}
	return requestedArtifacts, nil
}

func (d *DeploymentPlanGenerator) getArtifactsToDeploy(ctx context.Context, options DeploymentPlanOptions, artifacts DeployedArtifacts, plan *DeploymentPlan) error {
	// match services to be deployed
	for _, s := range artifacts {
		match := plan.requestedArtifacts.Match(s.RequestedVersion, s.ArtifactID)
		if match == nil {
			continue
		}
		if s.DeployedVersion == match.VersionInArtifactory && !options.ForceDeploy {
			if options.HasArtifacts() {
				plan.message("artifact: %s, version: %s, is already up to date in namespace: %s", s.ArtifactName, s.DeployedVersion, s.NamespaceName)
				match.matched = true
			}
			continue
		}
		s.AvailableVersion = match.VersionInArtifactory
		if s.AvailableVersion == "" || (s.DeployedVersion == s.AvailableVersion && !options.ForceDeploy) {
			continue
		}

		// stack environment, namespace, artifact and service in that order
		s.Metadata = mergeKeys(s.Metadata, plan.getMetadata(s.NamespaceID, s.ArtifactID))
		ns := plan.Namespaces.ID(s.NamespaceID)
		match.matched = true
		ns.Artifacts = append(ns.Artifacts, s)
	}

	// services were explicitly passed in
	if options.HasArtifacts() {
		unmatched := plan.requestedArtifacts.UnMatched()
		for _, x := range unmatched {
			plan.message("unmatched service: %s", x)
		}
	}

	return nil
}

func (d *DeploymentPlanGenerator) getNamespaces(ctx context.Context, env *Environment, requestedNamespaces StringList, logger messageLogger) (Namespaces, error) {
	// lets start with all the namespaces in the Env and filter it down based on additional information passed in.
	namespacesToDeploy, err := d.repo.NamespacesByEnvironmentID(ctx, env.ID)
	if err != nil {
		return nil, err
	}
	if len(namespacesToDeploy) == 0 {
		return nil, errors.NewRestError(400, "no associated namespaces in %s", env.Name)
	}
	if len(requestedNamespaces) > 0 {
		// Make sure that the namespaces that are specified are also available in the environment
		included, _ := filterNamespaces(namespacesToDeploy, func(namespace data.Namespace) bool {
			return requestedNamespaces.Contains(namespace.Alias)
		})
		for _, x := range requestedNamespaces {
			if !included.Contains(x) {
				return nil, errors.NewRestError(400, "invalid namespace: %s", x)
			}
		}
		namespacesToDeploy = included
	} else {
		// If we didn't specify any namespaces, we need to make sure were not deploying to namespaces that require you to explicitly specify them
		included, excluded := filterNamespaces(namespacesToDeploy, func(namespace data.Namespace) bool {
			return !namespace.ExplicitDeployOnly
		})
		namespacesToDeploy = included
		for _, x := range excluded {
			logger("explicit namespace excluded: %s", x.Alias)
		}
	}

	namespaces := fromDataNamespaces(namespacesToDeploy)
	return namespaces, nil
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

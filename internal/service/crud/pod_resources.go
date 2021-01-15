package crud

import (
	"context"
	"fmt"
	"strconv"

	"gitlab.unanet.io/devops/go/pkg/errors"
	"gitlab.unanet.io/devops/go/pkg/log"
	"go.uber.org/zap"

	"gitlab.unanet.io/devops/eve/internal/data"
	"gitlab.unanet.io/devops/eve/internal/service"
	"gitlab.unanet.io/devops/eve/pkg/eve"
)

func (m *Manager) fromDataPodResourceMaps(pams []data.PodResourcesMap) *eve.PodResources {
	stackedData, err := m.repo.PodResourcesStacked(pams)
	if err != nil {
		log.Logger.Error("failed to stack the pod resource maps")
		return nil
	}

	var sources []eve.PodResourcesMapSource
	for _, v := range pams {
		sources = append(sources, eve.PodResourcesMapSource{
			ArtifactID:                 v.ArtifactID,
			ServiceID:                  v.ServiceID,
			EnvironmentID:              v.EnvironmentID,
			NamespaceID:                v.NamespaceID,
			Data:                       v.Data,
			StackingOrder:              v.StackingOrder,
			PodResourcesDescription:    v.PodResourcesDescription,
			PodResourcesMapDescription: v.PodResourcesMapDescription,
		})
	}

	return &eve.PodResources{
		Sources: sources,
		Data:    stackedData,
	}
}

func (m *Manager) PodResourcesEnvironment(ctx context.Context, environmentID string) (*eve.PodResources, error) {
	var environment *data.Environment
	if len(environmentID) > 0 {
		if eID, err := strconv.Atoi(environmentID); err == nil {
			environment, err = m.repo.EnvironmentByID(ctx, eID)
			if err != nil {
				return nil, service.CheckForNotFoundError(err)
			}
		} else {
			environment, err = m.repo.EnvironmentByName(ctx, environmentID)
			if err != nil {
				return nil, service.CheckForNotFoundError(err)
			}
		}
	}
	if environment != nil {
		log.Logger.Info("pod resources by environment", zap.Any("environment", environment))

		prms, err := m.repo.PodResourcesMap(ctx, 0, environment.ID, 0, 0)
		if err != nil {
			return nil, errors.NotFoundf("environment pod resource maps not found: %v", environment.ID)
		}
		if len(prms) <= 0 {
			return nil, errors.NotFoundf("no pod resources config found for environment: %v", environmentID)
		}
		return m.fromDataPodResourceMaps(prms), nil
	}

	return nil, errors.NotFoundf("environment pod resources not found")
}

func (m *Manager) PodResourcesNamespace(ctx context.Context, namespaceID string) (*eve.PodResources, error) {
	var namespace *data.Namespace
	if len(namespaceID) > 0 {
		if nID, err := strconv.Atoi(namespaceID); err == nil {
			namespace, err = m.repo.NamespaceByID(ctx, nID)
			if err != nil {
				return nil, service.CheckForNotFoundError(err)
			}
		} else {
			namespace, err = m.repo.NamespaceByName(ctx, namespaceID)
			if err != nil {
				return nil, service.CheckForNotFoundError(err)
			}
		}
	}

	if namespace != nil {
		log.Logger.Info("pod resources by namespace", zap.Any("namespace", namespace))

		prms, err := m.repo.PodResourcesMap(ctx, 0, 0, namespace.ID, 0)
		if err != nil {
			return nil, errors.NotFoundf("namespace pod resource maps not found: %v", namespace.ID)
		}
		if len(prms) <= 0 {
			return nil, errors.NotFoundf("no pod resources config found for namespace: %v", namespaceID)
		}
		return m.fromDataPodResourceMaps(prms), nil
	}

	return nil, errors.NotFoundf("namespace pod resources not found")

}

func (m *Manager) PodResourcesArtifact(ctx context.Context, artifactID, environmentID, namespaceID string) (*eve.PodResources, error) {
	var artifact *data.Artifact
	if len(artifactID) > 0 {
		if aID, err := strconv.Atoi(artifactID); err == nil {
			artifact, err = m.repo.ArtifactByID(ctx, aID)
			if err != nil {
				return nil, service.CheckForNotFoundError(err)
			}
		} else {
			artifact, err = m.repo.ArtifactByName(ctx, artifactID)
			if err != nil {
				return nil, service.CheckForNotFoundError(err)
			}
		}
	}

	var environment *data.Environment
	if len(environmentID) > 0 {
		if eID, err := strconv.Atoi(environmentID); err == nil {
			environment, err = m.repo.EnvironmentByID(ctx, eID)
			if err != nil {
				return nil, service.CheckForNotFoundError(err)
			}
		} else {
			environment, err = m.repo.EnvironmentByName(ctx, environmentID)
			if err != nil {
				return nil, service.CheckForNotFoundError(err)
			}
		}
	}

	var namespace *data.Namespace
	if len(namespaceID) > 0 {
		if nID, err := strconv.Atoi(namespaceID); err == nil {
			namespace, err = m.repo.NamespaceByID(ctx, nID)
			if err != nil {
				return nil, service.CheckForNotFoundError(err)
			}
		} else {
			namespace, err = m.repo.NamespaceByName(ctx, namespaceID)
			if err != nil {
				return nil, service.CheckForNotFoundError(err)
			}
		}
	}

	if artifact != nil {
		var envID, nsID int
		log.Logger.Info("pod resources by artifact", zap.Any("artifact", artifact))

		if environment != nil {
			envID = environment.ID
		}
		if namespace != nil {
			nsID = namespace.ID
		}

		prms, err := m.repo.PodResourcesMap(ctx, 0, envID, nsID, artifact.ID)
		if err != nil {
			return nil, errors.NotFoundf("artifact pod resource maps not found: %v", artifact.ID)
		}
		if len(prms) <= 0 {
			return nil, errors.NotFoundf("no pod resources config found for artifact: %v", artifactID)
		}
		return m.fromDataPodResourceMaps(prms), nil

	}

	return nil, errors.NotFoundf("artifact pod resources not found")
}

func (m *Manager) PodResourcesService(ctx context.Context, serviceID, namespaceID string) (*eve.PodResources, error) {
	var svc *data.Service
	// If serviceID is the actual INT ID (Primary Key)
	// that's all we need
	if len(serviceID) > 0 {
		if svcID, err := strconv.Atoi(serviceID); err == nil {
			svc, err = m.repo.ServiceByID(ctx, svcID)
			if err != nil {
				return nil, service.CheckForNotFoundError(err)
			}
		} else {
			// If serviceID is the name of service
			// we also need the namespace NAME
			svc, err = m.repo.ServiceByName(ctx, serviceID, namespaceID)
			if err != nil {
				return nil, service.CheckForNotFoundError(err)
			}
		}
	}

	// Happy Path
	// A service input takes precedence over all other inputs
	// since a service contains an environment/namespaces we use the environment/namespace from the service
	if svc != nil {
		ns, err := m.repo.NamespaceByID(ctx, svc.NamespaceID)
		if err != nil {
			return nil, errors.NotFoundf("service namespace not found: %v", svc.NamespaceID)
		}
		env, err := m.repo.EnvironmentByID(ctx, ns.EnvironmentID)
		if err != nil {
			return nil, errors.NotFoundf("namespace environment not found: %v", ns.EnvironmentID)
		}
		pams, err := m.repo.PodResourcesMap(ctx, svc.ID, env.ID, svc.NamespaceID, svc.ArtifactID)
		if err != nil {
			return nil, errors.NotFoundf("service pod resources maps not found: %v", svc.ID)
		}
		if len(pams) <= 0 {
			return nil, errors.NotFoundf("no pod resources config found for service: %v", serviceID)
		}
		return m.fromDataPodResourceMaps(pams), nil
	}

	return nil, errors.NotFoundf("service pod resources not found")
}

func (m *Manager) PodResources(ctx context.Context, serviceID string, environmentID string, namespaceID string, artifactID string) (*eve.PodResources, error) {
	if validParam(serviceID, environmentID, namespaceID, artifactID) == false {
		return nil, errors.NewRestError(400, fmt.Sprintf("invalid input params serviceID: %s environmentID: %s namespaceID: %s artifact: %s", serviceID, environmentID, namespaceID, artifactID))
	}
	// Service Takes precedence since a service is made up of an Environment,Namespace, and Artifact
	// We pass the namespaceID as well, in-case the caller is passing Names and not IDs
	if len(serviceID) > 0 {
		return m.PodResourcesService(ctx, serviceID, namespaceID)
	}
	// Next we look for artifact matches
	// We pass the environment and namespace as well (and empty values are OK)
	// pod resource can be scoped to a specific artifact or and artifact in a specific environment/namespace
	if len(artifactID) > 0 {
		return m.PodResourcesArtifact(ctx, artifactID, environmentID, namespaceID)
	}
	// Then we look for namespace matches
	// namespaces take precedence over environment because a namespace belongs to an environment
	if len(namespaceID) > 0 {
		return m.PodResourcesNamespace(ctx, namespaceID)
	}
	// Lastly, we match on the environment
	if len(environmentID) > 0 {
		return m.PodResourcesEnvironment(ctx, environmentID)
	}
	// NOTE: If we ever want to match on Cluster (or maybe even provider like AKS vs EKS)...
	// ...it would be added here:
	return nil, errors.NotFoundf("pod resources not found")
}

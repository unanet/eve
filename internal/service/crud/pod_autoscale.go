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

func (m *Manager) fromDataPodAutoscaleMaps(pams []data.PodAutoscaleMap) *eve.PodAutoscale {
	stackedData, err := m.repo.PodAutoscaleStacked(pams)
	if err != nil {
		log.Logger.Error("failed to stack the pod autoscale maps")
		return nil
	}

	var sources []eve.PodAutoscaleMapSource
	for _, v := range pams {
		sources = append(sources, eve.PodAutoscaleMapSource{
			ServiceID:                  v.ServiceID,
			EnvironmentID:              v.EnvironmentID,
			NamespaceID:                v.NamespaceID,
			Data:                       v.Data,
			StackingOrder:              v.StackingOrder,
			PodAutoscaleDescription:    v.PodAutoscaleDescription,
			PodAutoscaleMapDescription: v.PodAutoscaleMapDescription,
		})
	}

	return &eve.PodAutoscale{
		Sources: sources,
		Data:    stackedData,
	}
}

func (m *Manager) PodAutoscaleEnvironment(ctx context.Context, environmentID string) (*eve.PodAutoscale, error) {
	var env *data.Environment
	if len(environmentID) > 0 {
		log.Logger.Info("PodAutoscale by environment", zap.Any("environment", environmentID))
		if envID, err := strconv.Atoi(environmentID); err == nil {
			env, err = m.repo.EnvironmentByID(ctx, envID)
			if err != nil {
				return nil, service.CheckForNotFoundError(err)
			}
		} else {
			env, err = m.repo.EnvironmentByName(ctx, environmentID)
			if err != nil {
				return nil, service.CheckForNotFoundError(err)
			}
		}
	}

	if env != nil {
		log.Logger.Info("PodAutoscale by environment", zap.Any("environment", env))
		pams, err := m.repo.EnvironmentPodAutoscaleMap(ctx, env.ID)
		if err != nil {
			return nil, service.CheckForNotFoundError(err)
		}
		if len(pams) <= 0 {
			return nil, errors.NotFoundf("no pod autoscaling config found for environment: %v", env.ID)
		}
		return m.fromDataPodAutoscaleMaps(pams), nil
	}

	return nil, errors.NotFoundf("environment pod autoscale not found")
}

func (m *Manager) PodAutoscaleNamespace(ctx context.Context, namespaceID string) (*eve.PodAutoscale, error) {
	var ns *data.Namespace
	if len(namespaceID) > 0 {
		log.Logger.Info("PodAutoscale by namespace", zap.String("namespace", namespaceID))
		if nsID, err := strconv.Atoi(namespaceID); err == nil {
			ns, err = m.repo.NamespaceByID(ctx, nsID)
			if err != nil {
				return nil, service.CheckForNotFoundError(err)
			}
			log.Logger.Info("PodAutoscale namespace by ID", zap.Any("namespace", ns))
		} else {
			ns, err = m.repo.NamespaceByName(ctx, namespaceID)
			if err != nil {
				return nil, service.CheckForNotFoundError(err)
			}
			log.Logger.Info("PodAutoscale namespace by Name", zap.Any("namespace", ns))
		}
	}

	if ns != nil {
		log.Logger.Info("PodAutoscale by namespace", zap.Any("namespace", ns))
		pams, err := m.repo.NamespacePodAutoscaleMap(ctx, ns.ID)
		if err != nil {
			return nil, service.CheckForNotFoundError(err)
		}
		if len(pams) <= 0 {
			return nil, errors.NotFoundf("no pod autoscaling config found for namespace: %v", ns.ID)
		}
		return m.fromDataPodAutoscaleMaps(pams), nil
	}
	return nil, errors.NotFoundf("namespace pod autoscale not found")
}

func (m *Manager) PodAutoscaleService(ctx context.Context, serviceID, namespaceID string) (*eve.PodAutoscale, error) {
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
		log.Logger.Info("PodAutoscale by service", zap.Any("service", svc))
		ns, err := m.repo.NamespaceByID(ctx, svc.NamespaceID)
		if err != nil {
			return nil, errors.NotFoundf("service namespace not found: %v", svc.NamespaceID)
		}
		env, err := m.repo.EnvironmentByID(ctx, ns.EnvironmentID)
		if err != nil {
			return nil, errors.NotFoundf("namespace environment not found: %v", ns.EnvironmentID)
		}
		pams, err := m.repo.PodAutoscaleMap(ctx, svc.ID, env.ID, svc.NamespaceID)
		if err != nil {
			return nil, errors.NotFoundf("service pod autoscale maps not found: %v", svc.ID)
		}
		if len(pams) <= 0 {
			return nil, errors.NotFoundf("no pod autoscaling config found for service: %v", serviceID)
		}
		return m.fromDataPodAutoscaleMaps(pams), nil
	}

	return nil, errors.NotFoundf("service pod autoscale not found")
}

// PodAutoscale finds and stacks the pod autoscale config defined in the pod_autoscale_map table
// A Service takes precedence over the Environment or Namespace
// A Namespaces takes precedence over the Environment
func (m *Manager) PodAutoscale(ctx context.Context, serviceID string, environmentID string, namespaceID string) (*eve.PodAutoscale, error) {
	log.Logger.Info("PodAutoscale",
		zap.String("service", serviceID),
		zap.String("environment", environmentID),
		zap.String("namespace", namespaceID),
	)
	if validParam(serviceID, environmentID, namespaceID) == false {
		return nil, errors.NewRestError(400, fmt.Sprintf("invalid input params serviceID: %s environmentID: %s namespaceID: %s", serviceID, environmentID, namespaceID))
	}

	if len(serviceID) > 0 {
		return m.PodAutoscaleService(ctx, serviceID, namespaceID)
	}

	if len(namespaceID) > 0 {
		return m.PodAutoscaleNamespace(ctx, namespaceID)
	}

	if len(environmentID) > 0 {
		return m.PodAutoscaleEnvironment(ctx, namespaceID)
	}

	return nil, errors.NotFoundf("pod autoscale not found")
}

// validParams just checks that at least one param was supplied
func validParam(params ...string) bool {
	var count int
	for _, p := range params {
		if len(p) >= 0 {
			count++
			break
		}
	}
	return count > 0
}

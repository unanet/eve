package crud

import (
	"gitlab.unanet.io/devops/eve/internal/data"
	"gitlab.unanet.io/devops/eve/pkg/eve"
)

func NewManager(r *data.Repo) *Manager {
	return &Manager{
		repo: r,
	}
}

type Manager struct {
	repo *data.Repo
}

// TODO: Handle this with Data default defs applied to everything (service/jobs)
// inspect the service definitions and make sure required defs are present
// if not, add default definitions
func (m *Manager) defaultServiceDefinitions(defs []eve.DefinitionResult) []eve.DefinitionResult {
	var validSvc,validDep bool
	for _, def := range defs {
		if def.Kind == "Service" {
			validSvc = true
			continue
		}
		if def.Kind == "Deployment" {
			validDep = true
			continue
		}
		if validSvc && validDep {
			break
		}
	}

	if !validSvc {
		defs = append(defs, eve.DefaultServiceResourceDef())
	}

	if !validDep {
		defs = append(defs, eve.DefaultDeploymentResourceDef())
	}

	return defs
}


// inspect the service definitions and make sure required defs are present
// if not, add default definitions
func (m *Manager) defaultJobDefinitions(defs []eve.DefinitionResult) []eve.DefinitionResult {
	var validJob bool

	for _, def := range defs {
		if def.Kind == "Job" {
			validJob = true
			break
		}
	}

	if !validJob {
		defs = append(defs, eve.DefaultJobResourceDef())
	}

	return defs
}







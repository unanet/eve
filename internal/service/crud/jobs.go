package crud

import (
	"context"
	"strconv"

	"gitlab.unanet.io/devops/eve/internal/data"
	"gitlab.unanet.io/devops/eve/internal/service"
	"gitlab.unanet.io/devops/eve/pkg/eve"
)

func fromDataJob(job data.Job) eve.Job {
	return eve.Job{
		ID:              job.ID,
		NamespaceID:     job.NamespaceID,
		NamespaceName:   job.NamespaceName,
		ArtifactID:      job.ArtifactID,
		ArtifactName:    job.ArtifactName,
		OverrideVersion: job.OverrideVersion.String,
		DeployedVersion: job.DeployedVersion.String,
		CreatedAt:       job.CreatedAt.Time,
		UpdatedAt:       job.UpdatedAt.Time,
		Name:            job.Name,
		NodeGroup:       job.NodeGroup,
	}
}

func fromDataJobs(services []data.Job) []eve.Job {
	var list []eve.Job
	for _, x := range services {
		list = append(list, fromDataJob(x))
	}
	return list
}

func toDataJob(j eve.Job) data.Job {
	d := data.Job{
		ID:            j.ID,
		NamespaceID:   j.NamespaceID,
		NamespaceName: j.NamespaceName,
		ArtifactID:    j.ArtifactID,
		ArtifactName:  j.ArtifactName,
		Name:          j.Name,
		NodeGroup:     j.NodeGroup,
	}

	if j.OverrideVersion != "" {
		d.OverrideVersion.String = j.OverrideVersion
		d.OverrideVersion.Valid = true
	}

	if j.DeployedVersion != "" {
		d.DeployedVersion.String = j.DeployedVersion
		d.DeployedVersion.Valid = true
	}

	return d
}

func (m *Manager) JobsByNamespace(ctx context.Context, namespaceID string) ([]eve.Job, error) {
	var dJobs []data.Job
	if intID, err := strconv.Atoi(namespaceID); err == nil {
		dJobs, err = m.repo.JobsByNamespaceID(ctx, intID)
		if err != nil {
			return nil, service.CheckForNotFoundError(err)
		}
	} else {
		dJobs, err = m.repo.JobsByNamespaceName(ctx, namespaceID)
		if err != nil {
			return nil, service.CheckForNotFoundError(err)
		}
	}

	return fromDataJobs(dJobs), nil
}

func (m *Manager) Job(ctx context.Context, id string, namespace string) (*eve.Job, error) {
	var d *data.Job
	if intID, err := strconv.Atoi(id); err == nil {
		d, err = m.repo.JobByID(ctx, intID)
		if err != nil {
			return nil, service.CheckForNotFoundError(err)
		}
	} else {
		d, err = m.repo.JobByName(ctx, id, namespace)
		if err != nil {
			return nil, service.CheckForNotFoundError(err)
		}
	}

	rd := fromDataJob(*d)
	return &rd, nil
}

func (m *Manager) UpdateJob(ctx context.Context, j *eve.Job) (*eve.Job, error) {
	d := toDataJob(*j)

	err := m.repo.UpdateJob(ctx, &d)
	if err != nil {
		return nil, service.CheckForNotFoundError(err)
	}

	j2 := fromDataJob(d)
	return &j2, nil
}

package crud

import (
	"context"
	"github.com/unanet/go/pkg/errors"

	"github.com/unanet/eve/internal/data"
	"github.com/unanet/eve/pkg/eve"
)

func (m *Manager) Artifacts(ctx context.Context) (models []eve.Artifact, err error) {
	dbArtifacts, err := m.repo.Artifact(ctx)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	return fromDataArtifactList(dbArtifacts), err
}

func (m *Manager) CreateArtifact(ctx context.Context, artifact *eve.Artifact) error {

	dbArtifact := toDataArtifact(*artifact)
	if err := m.repo.CreateArtifact(ctx, &dbArtifact); err != nil {
		return err
	}

	artifact.ID = dbArtifact.ID
	return nil
}

func fromDataArtifactList(artifacts []data.Artifact) []eve.Artifact {
	var list []eve.Artifact
	for _, x := range artifacts {
		list = append(list, fromDataArtifactToArtifact(x))
	}
	return list
}

func (m *Manager) UpdateArtifact(ctx context.Context, model *eve.Artifact) (err error) {
	dbModel := toDataArtifact(*model)
	if err := m.repo.UpdateArtifact(ctx, &dbModel); err != nil {
		return err
	}

	return nil
}

func (m *Manager) DeleteArtifact(ctx context.Context, id int) (err error) {
	return m.repo.DeleteArtifact(ctx, id)
}

func fromDataArtifactToArtifact(m data.Artifact) eve.Artifact {
	return eve.Artifact{
		ID:            m.ID,
		Name:          m.Name,
		FeedType:      m.FeedType,
		ProviderGroup: m.ProviderGroup,
		ImageTag:      m.ImageTag,
		ServicePort:   m.ServicePort,
		MetricsPort:   m.MetricsPort,
	}
}

func toDataArtifact(m eve.Artifact) data.Artifact {
	return data.Artifact{
		ID:            m.ID,
		Name:          m.Name,
		FeedType:      m.FeedType,
		ProviderGroup: m.ProviderGroup,
		ImageTag:      m.ImageTag,
		ServicePort:   m.ServicePort,
		MetricsPort:   m.MetricsPort,
	}
}

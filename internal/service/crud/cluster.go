package crud

import (
	"context"
	"gitlab.unanet.io/devops/eve/internal/data"
	"gitlab.unanet.io/devops/eve/pkg/eve"
	"gitlab.unanet.io/devops/go/pkg/errors"
)

func (m *Manager) Clusters(ctx context.Context) (models []eve.Cluster, err error) {
	dbClusters, err := m.repo.Clusters(ctx)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	return fromDataClusterList(dbClusters), nil
}

func (m *Manager) CreateCluster(ctx context.Context, model *eve.Cluster) error {

	dbModel := toDataCluster(*model)
	if err := m.repo.CreateCluster(ctx, &dbModel); err != nil {
		return errors.Wrap(err)
	}

	model.ID = dbModel.ID
	model.CreatedAt = dbModel.CreatedAt.Time
	model.UpdatedAt = dbModel.UpdatedAt.Time

	return nil
}

func (m *Manager) UpdateCluster(ctx context.Context, model *eve.Cluster) (err error) {
	dbModel := toDataCluster(*model)
	if err := m.repo.UpdateCluster(ctx, &dbModel); err != nil {
		return err
	}

	model.CreatedAt = dbModel.CreatedAt.Time
	model.UpdatedAt = dbModel.UpdatedAt.Time

	return nil
}

func (m *Manager) DeleteCluster(ctx context.Context, id int) (err error) {
	return m.repo.DeleteCluster(ctx, id)
}

func fromDataClusterList(clusters []data.Cluster) []eve.Cluster {
	var list []eve.Cluster
	for _, x := range clusters {
		list = append(list, fromDataClusterToCluster(x))
	}
	return list
}

func fromDataClusterToCluster(m data.Cluster) eve.Cluster {
	return eve.Cluster{
		ID:            m.ID,
		Name:          m.Name,
		ProviderGroup: m.ProviderGroup,
		SchQueueUrl:   m.SchQueueUrl,
		CreatedAt:     m.CreatedAt.Time,
		UpdatedAt:     m.UpdatedAt.Time,
	}
}

func toDataCluster(m eve.Cluster) data.Cluster {
	return data.Cluster{
		ID:            m.ID,
		Name:          m.Name,
		ProviderGroup: m.ProviderGroup,
		SchQueueUrl:   m.SchQueueUrl,
	}
}

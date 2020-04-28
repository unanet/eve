package data

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"

	"gitlab.unanet.io/devops/eve/pkg/errors"
)

type Service struct {
	ID               int            `db:"id"`
	NamespaceID      int            `db:"namespace_id"`
	NamespaceName    string         `db:"namespace_name"`
	ArtifactID       int            `db:"artifact_id"`
	ArtifactName     string         `db:"artifact_name"`
	RequestedVersion string         `db:"requested_version"`
	OverrideVersion  sql.NullString `db:"override_version"`
	DeployedVersion  sql.NullString `db:"deployed_version"`
	Metadata         JSONText       `db:"metadata"`
	CreatedAt        sql.NullTime   `db:"created_at"`
	UpdatedAt        sql.NullTime   `db:"updated_at"`
}

type Services []Service

func (r *Repo) DeployedServicesByNamespaceIDs(ctx context.Context, namespaceIDs []interface{}) (DeployedArtifacts, error) {
	sql, args, err := sqlx.In(`
		select s.id, 
		       s.namespace_id,
		       n.name as namespace_name,
		       s.artifact_id,
		       a.name as artifact_name, 
		       s.deployed_version, 
		       s.metadata,
		    COALESCE(s.override_version, n.requested_version) as requested_version 
		from service as s 
		    left join artifact as a on a.id = s.artifact_id
			left join namespace n on s.namespace_id = n.id
		where s.namespace_id in (?)
			`, namespaceIDs)
	sql = r.db.Rebind(sql)
	rows, err := r.db.QueryxContext(ctx, sql, args...)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	var deployedArtifacts []DeployedArtifact
	for rows.Next() {
		var deployedArtifact DeployedArtifact
		err = rows.StructScan(&deployedArtifact)
		if err != nil {
			return nil, errors.Wrap(err)
		}
		deployedArtifacts = append(deployedArtifacts, deployedArtifact)
	}

	return deployedArtifacts, nil
}

func (r *Repo) ServiceArtifacts(ctx context.Context, namespaceIDs []int) (RequestArtifacts, error) {
	sql, args, err := sqlx.In(`
		select distinct s.artifact_id, 
		                a.metadata as artifact_metadata,
		                a.name as artifact_name, 
		                a.provider_group as provider_group,
		                f.name as feed_name,
		                COALESCE(s.override_version, ns.requested_version) as requested_version 
		from service as s 
			left join namespace as ns on ns.id = s.namespace_id
			left join artifact as a on a.id = s.artifact_id
			left join environment e on ns.environment_id = e.id
			left join environment_feed_map efm on e.id = efm.environment_id
			left join feed f on efm.feed_id = f.id and f.feed_type = a.feed_type
		where f.name is not null and ns.id in (?)`, namespaceIDs)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	sql = r.db.Rebind(sql)
	rows, err := r.db.QueryxContext(ctx, sql, args...)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	var services []RequestArtifact
	for rows.Next() {
		var service RequestArtifact
		err = rows.StructScan(&service)
		if err != nil {
			return nil, errors.Wrap(err)
		}
		services = append(services, service)
	}
	return services, nil
}

package data

import (
	"context"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"gitlab.unanet.io/devops/eve/pkg/errors"
)

func (r *Repo) DatabaseInstanceArtifacts(ctx context.Context, namespaceIDs []interface{}) (RequestArtifacts, error) {
	sql, args, err := sqlx.In(`
			select
			 	a.id as artifact_id,
				a.name as artifact_name,
				a.provider_group as provider_group,
				f.name as feed_name,
			    a.metadata as artifact_metadata,
				ds.metadata as server_metadata,
				COALESCE(di.migration_override_version, ns.requested_version) as requested_version
			from database_instance as di
			    left join database_server ds on di.database_server_id = ds.id
			    left join namespace ns on di.namespace_id = ns.id
				left join database_type dt on di.database_type_id = dt.id
				left join artifact a on dt.migration_artifact_id = a.id
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
	var artifacts RequestArtifacts
	for rows.Next() {
		var artifact RequestArtifact
		err = rows.StructScan(&artifact)
		if err != nil {
			return nil, errors.Wrap(err)
		}
		artifacts = append(artifacts, artifact)
	}
	return artifacts, nil
}

func (r *Repo) DeployedDatabaseInstancesByNamespaceIDs(ctx context.Context, namespaceIDs []interface{}) (DeployedArtifacts, error) {
	sql, args, err := sqlx.In(`
		select 
			di.id,
			ns.id as namespace_id,
		    ns.name as namespace_name,
		    a.id as artifact_id, 
			a.name as artifact_name,
		    di.migration_deployed_version as deployed_version,
		    c.name as customer_name,
			di.metadata as metadata,
			COALESCE(di.migration_override_version, ns.requested_version) as requested_version 
		from database_instance as di 
		    left join database_server ds on di.database_server_id = ds.id
		    left join customer c on di.customer_id = c.id
			left join database_type dt on di.database_type_id = dt.id
		    left join artifact a on dt.migration_artifact_id = a.id
		    left join namespace ns on di.namespace_id = ns.id
		where a.id is not null and ns.id in (?)`, namespaceIDs)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	sql = r.db.Rebind(sql)
	rows, err := r.db.QueryxContext(ctx, sql, args...)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	var artifacts DeployedArtifacts
	for rows.Next() {
		var artifact DeployedArtifact
		err = rows.StructScan(&artifact)
		if err != nil {
			return nil, errors.Wrap(err)
		}
		artifacts = append(artifacts, artifact)
	}
	return artifacts, nil
}

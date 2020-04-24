package data

import (
	"context"
	"database/sql"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"gitlab.unanet.io/devops/eve/pkg/errors"
)

type DatabaseInstance struct {
	Id               int            `db:"id"`
	Name             string         `db:"name"`
	DatabaseTypeId   int            `db:"database_type_id"`
	DatabaseServerId int            `db:"database_server_id"`
	CustomerId       int            `db:"customer_id"`
	NamespaceId      int            `db:"namespace_id"`
	OverrideVersion  sql.NullString `db:"override_version"`
	DeployedVersion  string         `db:"deployed_version"`
	Metadata         JSONText       `db:"metadata"`
	CreatedAt        sql.NullTime   `db:"created_at"`
	UpdatedAt        sql.NullTime   `db:"updated_at"`
}

type DatabaseInstances []DatabaseInstance

func (r *Repo) DatabaseInstancesByNamespaceIDs(ctx context.Context, namespaceIDs []interface{}) (DatabaseInstances, error) {
	sql, args, err := sqlx.In(`
		select di.*,
		 	a.id as artifact_id, 
			a.metadata as artifact_metadata,
			a.name as artifact_name, 
			a.provider_group as provider_group,
			f.name as feed_name,
			COALESCE(di.override_migration_version, ns.requested_version) as requested_version 
		from database_instance as di
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
	var databaseInstances []DatabaseInstance
	for rows.Next() {
		var databaseInstance DatabaseInstance
		err = rows.StructScan(&databaseInstance)
		if err != nil {
			return nil, errors.Wrap(err)
		}
		databaseInstances = append(databaseInstances, databaseInstance)
	}
	return databaseInstances, nil
}

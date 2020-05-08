package data

import (
	"context"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"

	"gitlab.unanet.io/devops/eve/pkg/errors"
	"gitlab.unanet.io/devops/eve/pkg/json"
)

type Service struct {
	ServiceID        int            `db:"service_id"`
	ArtifactID       int            `db:"artifact_id"`
	ArtifactName     string         `db:"artifact_name"`
	RequestedVersion string         `db:"requested_version"`
	DeployedVersion  sql.NullString `db:"deployed_version"`
	Metadata         json.Text      `db:"metadata"`
	InjectVaultPaths sql.NullString `db:"inject_vault_paths"`
	CreatedAt        sql.NullTime   `db:"created_at"`
	UpdatedAt        sql.NullTime   `db:"updated_at"`
}

type Services []Service

func (r *Repo) UpdateDeployedServiceVersion(ctx context.Context, id int, version string) error {
	result, err := r.db.ExecContext(ctx, "update service set deployed_version = $1, updated_at = $2 where id = $3", version, time.Now().UTC(), id)
	if err != nil {
		return errors.Wrap(err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err)
	}

	if affected == 0 {
		return errors.Wrapf("the following id: %d was not found to update in service table", id)
	}
	return nil
}

func (r *Repo) DeployedServicesByNamespaceID(ctx context.Context, namespaceID int) (Services, error) {
	rows, err := r.db.QueryxContext(ctx, `
		select s.id as service_id, 
		   s.artifact_id,
		   a.name as artifact_name, 
		   s.deployed_version,
		   s.inject_vault_paths,
		   jsonb_merge(e.metadata, jsonb_merge(n.metadata, jsonb_merge(a.metadata, s.metadata))) as metadata,
		   COALESCE(s.override_version, n.requested_version) as requested_version,
		   s.created_at,
		   s.updated_at
		from service as s 
		    left join artifact as a on a.id = s.artifact_id
			left join namespace n on s.namespace_id = n.id
			left join environment e on n.environment_id = e.id
		where s.namespace_id = $1
	`, namespaceID)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	var services []Service
	for rows.Next() {
		var service Service
		err = rows.StructScan(&service)
		if err != nil {
			return nil, errors.Wrap(err)
		}
		services = append(services, service)
	}
	return services, nil
}

func (r *Repo) ServiceArtifacts(ctx context.Context, namespaceIDs []int) (RequestArtifacts, error) {
	esql, args, err := sqlx.In(`
		select distinct s.artifact_id, 
		                a.function_pointer as function_pointer,
		                a.name as artifact_name, 
		                a.provider_group as provider_group,
		                a.feed_type as feed_type,
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
	esql = r.db.Rebind(esql)
	rows, err := r.db.QueryxContext(ctx, esql, args...)
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

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
	ArtifactID       int            `db:"artifact_id"`
	ArtifactName     string         `db:"artifact_name"`
	RequestedVersion string         `db:"requested_version"`
	OverrideVersion  sql.NullString `db:"override_version"`
	DeployedVersion  sql.NullString `db:"deployed_version"`
	Metadata         JSONText       `db:"metadata"`
	CreatedAt        sql.NullTime   `db:"created_at"`
	UpdatedAt        sql.NullTime   `db:"updated_at"`
}

type RequestedArtifact struct {
	ArtifactID       int      `db:"artifact_id"`
	ArtifactName     string   `db:"artifact_name"`
	ProviderGroup    string   `db:"provider_group"`
	FeedName         string   `db:"feed_name"`
	ArtifactMetadata JSONText `db:"artifact_metadata"`
	RequestedVersion string   `db:"requested_version"`
}

type RequestedArtifacts []RequestedArtifact

type Services []Service

func (r *Repo) Services(ctx context.Context) (Services, error) {
	return r.services(ctx)
}

func (r *Repo) ServicesByNamespaceIDs(ctx context.Context, namespaceIDs []interface{}) (Services, error) {
	return r.services(ctx, WhereIn("namespace_id", namespaceIDs))
}

func (r *Repo) RequestedArtifactByEnvironment(ctx context.Context, artifactName string, environmentID int) (*RequestedArtifact, error) {
	db := r.getDB()
	defer db.Close()

	var requestedArtifact RequestedArtifact

	row := db.QueryRowxContext(ctx, `
		select a.id as artifact_id,
		       a.name as artifact_name,
		       a.provider_group as provider_group,
		       f.name as feed_name
		from artifact as a
		    left join environment e on e.id = $1
		    left join environment_feed_map efm on e.id = efm.environment_id
			left join feed f on efm.feed_id = f.id
		where a.name = $2
	`, environmentID, artifactName)

	err := row.StructScan(&requestedArtifact)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, NotFoundErrorf("artifact with name: %s not found", artifactName)
		}
		return nil, errors.WrapUnexpected(err)
	}

	return &requestedArtifact, nil
}

func (r *Repo) RequestedArtifacts(ctx context.Context, namespaceIDs []interface{}) (RequestedArtifacts, error) {
	db := r.getDB()
	defer db.Close()

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
		return nil, errors.WrapUnexpected(err)
	}
	sql = db.Rebind(sql)
	rows, err := db.QueryxContext(ctx, sql, args...)
	if err != nil {
		return nil, errors.WrapUnexpected(err)
	}
	var services []RequestedArtifact
	for rows.Next() {
		var service RequestedArtifact
		err = rows.StructScan(&service)
		if err != nil {
			return nil, errors.WrapUnexpected(err)
		}
		services = append(services, service)
	}
	return services, nil
}

func (r *Repo) services(ctx context.Context, whereArgs ...WhereArg) (Services, error) {
	db := r.getDB()
	defer db.Close()

	sql, args := CheckWhereArgs(`
		select s.*, a.name as artifact_name,
		    COALESCE(s.override_version, n.requested_version) as requested_version 
		from service as s 
		    left join artifact as a on a.id = s.artifact_id
			left join namespace n on s.namespace_id = n.id
`, whereArgs)
	rows, err := db.QueryxContext(ctx, sql, args...)
	if err != nil {
		return nil, errors.WrapUnexpected(err)
	}
	var services []Service
	for rows.Next() {
		var service Service
		err = rows.StructScan(&service)
		if err != nil {
			return nil, errors.WrapUnexpected(err)
		}
		services = append(services, service)
	}

	return services, nil
}

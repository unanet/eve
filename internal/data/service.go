package data

import (
	"context"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"

	"gitlab.unanet.io/devops/eve/pkg/errors"
	"gitlab.unanet.io/devops/eve/pkg/json"
)

type DeployService struct {
	ServiceID         int            `db:"service_id"`
	ServiceName       string         `db:"service_name"`
	ArtifactID        int            `db:"artifact_id"`
	ArtifactName      string         `db:"artifact_name"`
	RequestedVersion  string         `db:"requested_version"`
	DeployedVersion   sql.NullString `db:"deployed_version"`
	ServicePort       int            `db:"service_port"`
	MetricsPort       int            `db:"metrics_port"`
	ServiceAccount    string         `db:"service_account"`
	ImageTag          string         `db:"image_tag"`
	RunAs             int            `db:"run_as"`
	StickySessions    bool           `db:"sticky_sessions"`
	Count             int            `db:"count"`
	MinPod            int            `db:"min_pod"`
	MaxPod            int            `db:"max_pod"`
	Metadata          json.Text      `db:"metadata"`
	LivelinessProbe   json.Text      `db:"liveliness_probe"`
	ReadinessProbe    json.Text      `db:"readiness_probe"`
	ResourceLimits    json.Text      `db:"resource_limits"`
	ResourceRequests  json.Text      `db:"resource_requests"`
	UtilizationLimits json.Text      `db:"utilization_limits"`
	CreatedAt         sql.NullTime   `db:"created_at"`
	UpdatedAt         sql.NullTime   `db:"updated_at"`
}

type DeployServices []DeployService

type Service struct {
	ID                int            `db:"id"`
	NamespaceID       int            `db:"namespace_id"`
	NamespaceName     string         `db:"namespace_name"`
	ArtifactID        int            `db:"artifact_id"`
	ArtifactName      string         `db:"artifact_name"`
	OverrideVersion   sql.NullString `db:"override_version"`
	DeployedVersion   sql.NullString `db:"deployed_version"`
	Metadata          json.Text      `db:"metadata"`
	CreatedAt         sql.NullTime   `db:"created_at"`
	UpdatedAt         sql.NullTime   `db:"updated_at"`
	Name              string         `db:"name"`
	StickySessions    bool           `db:"sticky_sessions"`
	Count             int            `db:"count"`
	MinPod            int            `db:"min_pod"`
	MaxPod            int            `db:"max_pod"`
	ResourceLimits    json.Text      `db:"resource_limits"`
	ResourceRequests  json.Text      `db:"resource_requests"`
	UtilizationLimits json.Text      `db:"utilization_limits"`
}

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

func (r *Repo) DeployedServicesByNamespaceID(ctx context.Context, namespaceID int) (DeployServices, error) {
	rows, err := r.db.QueryxContext(ctx, `
		select s.id as service_id,
		   a.service_port,
		   a.liveliness_probe,
		   a.readiness_probe,
		   jsonb_merge(a.resource_limits,s.resource_limits) as resource_limits,
	       jsonb_merge(a.resource_requests,s.resource_requests) as resource_requests, 
		   jsonb_merge(a.utilization_limits,s.utilization_limits) as utilization_limits, 
		   a.image_tag,
		   a.metrics_port,
		   a.service_account,
		   s.sticky_sessions,
		   s.count,
		   s.min_pod,
		   s.max_pod,
           s.name as service_name,
		   a.run_as as run_as,
		   s.artifact_id,
		   a.name as artifact_name, 
		   s.deployed_version,
		   jsonb_merge(e.metadata, jsonb_merge(a.metadata, jsonb_merge(n.metadata, s.metadata))) as metadata,
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
	var services []DeployService
	for rows.Next() {
		var service DeployService
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

func (r *Repo) ServiceByName(ctx context.Context, name string, namespace string) (*Service, error) {
	var service Service

	row := r.db.QueryRowxContext(ctx, `
		select s.*, n.name as namespace_name, a.name as artifact_name
		from service s 
		    left join namespace n on s.namespace_id = n.id
			left join artifact a on s.artifact_id = a.id
		where s.name = $1 and n.name = $2
		`, name, namespace)
	err := row.StructScan(&service)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, NotFoundErrorf("service with name: %s, namespace: %s, not found", name, namespace)
		}
		return nil, errors.Wrap(err)
	}

	return &service, nil
}

func (r *Repo) ServiceByID(ctx context.Context, id int) (*Service, error) {
	var service Service

	row := r.db.QueryRowxContext(ctx, `
		select s.*, n.name as namespace_name, a.name as artifact_name
		from service s 
		    left join namespace n on s.namespace_id = n.id
			left join artifact a on s.artifact_id = a.id
		where s.id = $1
		`, id)
	err := row.StructScan(&service)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, NotFoundErrorf("service with id: %d, not found", id)
		}
		return nil, errors.Wrap(err)
	}

	return &service, nil
}

func (r *Repo) ServicesByNamespaceID(ctx context.Context, namespaceID int) ([]Service, error) {
	return r.services(ctx, Where("s.namespace_id", namespaceID))
}

func (r *Repo) ServicesByNamespaceName(ctx context.Context, namespaceName string) ([]Service, error) {
	return r.services(ctx, Where("n.name", namespaceName))
}

func (r *Repo) services(ctx context.Context, whereArgs ...WhereArg) ([]Service, error) {
	esql, args := CheckWhereArgs(`
		select s.id, 
		       s.namespace_id, 
		       s.artifact_id, 
		       s.override_version, 
		       s.deployed_version, 
		       s.created_at, 
		       s.updated_at, 
		       s.name, 
		       s.sticky_sessions, 
		       s.count,
               s.min_pod,
               s.max_pod,
		       n.name as namespace_name,
		       a.name as artifact_name
		from service s 
		    left join namespace n on s.namespace_id = n.id
			left join artifact a on s.artifact_id = a.id
		`, whereArgs)
	rows, err := r.db.QueryxContext(ctx, esql+"order by s.name", args...)
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

func (r *Repo) UpdateService(ctx context.Context, service *Service) error {
	service.UpdatedAt.Time = time.Now().UTC()
	service.UpdatedAt.Valid = true
	if service.Metadata == nil {
		service.Metadata = json.EmptyJSONText
	}

	result, err := r.db.ExecContext(ctx, `
		update service set 
		   	name = $1, 
			namespace_id = $2,
			artifact_id = $3,
		   	override_version = $4,
		    deployed_version = $5,
		   	metadata = $6,
			sticky_sessions = $7,
		    count = $8,
		    updated_at = $9
		where id = $10
	`,
		service.Name,
		service.NamespaceID,
		service.ArtifactID,
		service.OverrideVersion,
		service.DeployedVersion,
		service.Metadata,
		service.StickySessions,
		service.Count,
		service.UpdatedAt,
		service.ID)
	if err != nil {
		return errors.Wrap(err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err)
	}

	if affected == 0 {
		return errors.NotFoundf("service id: %d not found", service.ID)
	}
	return nil
}

func (r *Repo) UpdateServiceCount(ctx context.Context, serviceID int, count int) error {
	if count > 2 || count < 0 {
		return errors.BadRequest("service count must be between > -1 and less than 3")
	}
	result, err := r.db.ExecContext(ctx, `
		update service set count = $1 where id = $2
	`, count, serviceID)
	if err != nil {
		return errors.Wrap(err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err)
	}

	if affected == 0 {
		return errors.NotFoundf("service id: %d not found", serviceID)
	}
	return nil
}

func (r *Repo) UpdateServiceMetadata(ctx context.Context, serviceID int, metadata map[string]interface{}) error {
	m := json.FromMap(metadata)
	result, err := r.db.ExecContext(ctx, `
		update service set metadata = metadata || $1 where id = $2
	`, m, serviceID)
	if err != nil {
		return errors.Wrap(err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err)
	}

	if affected == 0 {
		return errors.NotFoundf("service id: %d not found", serviceID)
	}
	return nil
}

func (r *Repo) DeleteServiceMetadataKey(ctx context.Context, serviceID int, key string) error {
	result, err := r.db.ExecContext(ctx, `
		update service set metadata = metadata - $1 where id = $2
	`, key, serviceID)
	if err != nil {
		return errors.Wrap(err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err)
	}

	if affected == 0 {
		return errors.NotFoundf("service id: %d not found", serviceID)
	}
	return nil
}

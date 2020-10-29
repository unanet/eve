package data

import (
	"context"
	"database/sql"
	goErrors "errors"
	"time"

	"gitlab.unanet.io/devops/eve/pkg/errors"
	"gitlab.unanet.io/devops/eve/pkg/json"
)

type Namespace struct {
	ID                 int          `db:"id"`
	Name               string       `db:"name"`
	Alias              string       `db:"alias"`
	EnvironmentID      int          `db:"environment_id"`
	EnvironmentName    string       `db:"environment_name"`
	RequestedVersion   string       `db:"requested_version"`
	ExplicitDeployOnly bool         `db:"explicit_deploy_only"`
	ClusterID          int          `db:"cluster_id"`
	Metadata           json.Text    `db:"metadata"`
	CreatedAt          sql.NullTime `db:"created_at"`
	UpdatedAt          sql.NullTime `db:"updated_at"`
}

type Namespaces []Namespace

func (n Namespaces) ToAliases() []string {
	var aliases []string
	for _, x := range n {
		aliases = append(aliases, x.Alias)
	}
	return aliases
}

func (n Namespaces) ToIDs() []int {
	var ids []int
	for _, x := range n {
		ids = append(ids, x.ID)
	}
	return ids
}

func (n Namespaces) Contains(name string) bool {
	for _, x := range n {
		if x.Alias == name || x.Name == name {
			return true
		}
	}

	return false
}

func (n Namespaces) FilterNamespaces(filter func(namespace Namespace) bool) (Namespaces, Namespaces) {
	var included Namespaces
	var excluded Namespaces
	for _, x := range n {
		if filter(x) {
			included = append(included, x)
		} else {
			excluded = append(excluded, x)
		}
	}

	return included, excluded
}

func (r *Repo) NamespaceByName(ctx context.Context, name string) (*Namespace, error) {
	var namespace Namespace

	row := r.db.QueryRowxContext(ctx, "select * from namespace where name = $1", name)
	err := row.StructScan(&namespace)
	if err != nil {
		if goErrors.Is(err, sql.ErrNoRows) {
			return nil, NotFoundErrorf("namespace with name: %s, not found", name)
		}
		return nil, errors.Wrap(err)
	}

	return &namespace, nil
}

func (r *Repo) NamespaceByID(ctx context.Context, id int) (*Namespace, error) {
	var namespace Namespace

	row := r.db.QueryRowxContext(ctx, "select * from namespace where id = $1", id)
	err := row.StructScan(&namespace)
	if err != nil {
		if goErrors.Is(err, sql.ErrNoRows) {
			return nil, NotFoundErrorf("namespace with id: %d, not found", id)
		}
		return nil, errors.Wrap(err)
	}

	return &namespace, nil
}

func (r *Repo) NamespacesByEnvironmentID(ctx context.Context, environmentID int) (Namespaces, error) {
	return r.namespaces(ctx, Where("ns.environment_id", environmentID))
}

func (r *Repo) NamespacesByEnvironmentName(ctx context.Context, environmentName string) (Namespaces, error) {
	return r.namespaces(ctx, Where("e.name", environmentName))
}

func (r *Repo) Namespaces(ctx context.Context) (Namespaces, error) {
	return r.namespaces(ctx)
}

func (r *Repo) namespaces(ctx context.Context, whereArgs ...WhereArg) (Namespaces, error) {
	esql, args := CheckWhereArgs(`
		select ns.id, 
		       ns.alias, 
		       ns.name, 
		       ns.environment_id, 
		       ns.requested_version, 
		       ns.explicit_deploy_only, 
		       ns.cluster_id,
		       ns.created_at,
		       ns.updated_at,
		       e.name as environment_name 
		from namespace ns left join environment e on ns.environment_id = e.id
		`, whereArgs)
	rows, err := r.db.QueryxContext(ctx, esql+"order by ns.requested_version desc", args...)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	var namespaces []Namespace
	for rows.Next() {
		var namespace Namespace
		err = rows.StructScan(&namespace)
		if err != nil {
			return nil, errors.Wrap(err)
		}
		namespaces = append(namespaces, namespace)
	}

	return namespaces, nil
}

func (r *Repo) UpdateNamespace(ctx context.Context, namespace *Namespace) error {
	namespace.UpdatedAt.Time = time.Now().UTC()
	namespace.UpdatedAt.Valid = true
	result, err := r.db.ExecContext(ctx, `
		update namespace set 
			requested_version = $1,
			explicit_deploy_only = $2,
			metadata = $3,
			updated_at = $4
		where id = $5
	`,
		namespace.RequestedVersion,
		namespace.ExplicitDeployOnly,
		namespace.Metadata,
		namespace.UpdatedAt,
		namespace.ID)
	if err != nil {
		return errors.Wrap(err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err)
	}

	if affected == 0 {
		return errors.NotFoundf("namespace id: %d not found", namespace.ID)
	}
	return nil
}

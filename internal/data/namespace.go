package data

import (
	"context"
	"database/sql"

	"gitlab.unanet.io/devops/eve/pkg/errors"
	"gitlab.unanet.io/devops/eve/pkg/json"
)

type Namespace struct {
	ID                 int          `db:"id"`
	Alias              string       `db:"alias"`
	EnvironmentID      int          `db:"environment_id"`
	Domain             string       `db:"domain"`
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
		if x.Alias == name {
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
		if err.Error() == "sql: no rows in result set" {
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
		if err.Error() == "sql: no rows in result set" {
			return nil, NotFoundErrorf("namespace with id: %d, not found", id)
		}
		return nil, errors.Wrap(err)
	}

	return &namespace, nil
}

func (r *Repo) NamespacesByEnvironmentID(ctx context.Context, environmentID int) (Namespaces, error) {
	return r.namespaces(ctx, Where("environment_id", environmentID))
}

func (r *Repo) namespaces(ctx context.Context, whereArgs ...WhereArg) (Namespaces, error) {
	esql, args := CheckWhereArgs("SELECT * FROM namespace", whereArgs)
	rows, err := r.db.QueryxContext(ctx, esql, args...)
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

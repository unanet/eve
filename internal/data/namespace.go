package data

import (
	"context"
	"fmt"
	"time"

	"gitlab.unanet.io/devops/eve/internal/data/common"
)

type Namespace struct {
	ID                 int             `db:"id"`
	Name               string          `db:"name"`
	Alias              string          `db:"alias"`
	EnvironmentID      int             `db:"environment_id"`
	Domain             string          `db:"domain"`
	DefaultVersion     string          `db:"default_version"`
	ExplicitDeployOnly bool            `db:"explicit_deploy_only"`
	ClusterID          int             `db:"cluster_id"`
	Metadata           common.JSONText `db:"metadata"`
	CreatedAt          *time.Time      `db:"created_at"`
	UpdatedAt          *time.Time      `db:"updated_at"`
}

func (r *Repo) GetNamespaces(ctx context.Context, whereArgs ...common.WhereArg) ([]Namespace, error) {
	db := r.getDB()
	defer db.Close()

	sql, args := common.CheckWhereArgs("select * from namespace", whereArgs)
	rows, err := db.QueryxContext(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("Repo.GetNamespaces.QueryxContext Error: %w", err)
	}
	var namespaces []Namespace
	for rows.Next() {
		var namespace Namespace
		err = rows.StructScan(&namespace)
		if err != nil {
			return nil, fmt.Errorf("Repo.GetNamespaces.StructScan Error: %w", err)
		}
		namespaces = append(namespaces, namespace)
	}

	return namespaces, nil
}

func (r *Repo) GetNamespaceByID(ctx context.Context, id int) (*Namespace, error) {
	db := r.getDB()
	defer db.Close()

	var namespace Namespace

	row := db.QueryRowxContext(ctx, "select * from namespace where id = $1", id)
	err := row.StructScan(&namespace)
	if err != nil {
		return nil, fmt.Errorf("Repo.GetNamespaceByID.QueryRowxContext Error: %w", err)
	}

	return &namespace, nil
}

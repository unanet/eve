package data

import (
	"context"
	"time"

	"gitlab.unanet.io/devops/eve/internal/data/orm"
	"gitlab.unanet.io/devops/eve/pkg/errors"
)

type Cluster struct {
	ID            string     `db:"id"`
	Name          string     `db:"name"`
	ProviderGroup string     `db:"provider_group"`
	CreatedAt     *time.Time `db:"created_at"`
	UpdatedAt     *time.Time `db:"updated_at"`
}

type Clusters []Cluster

func (r *Repo) ClusterByID(ctx context.Context, id int) (*Cluster, error) {
	db := r.getDB()
	defer db.Close()

	var cluster Cluster

	row := db.QueryRowxContext(ctx, "select * from cluster where id = $1", id)
	err := row.StructScan(&cluster)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, NotFoundErrorf("cluster with id: %d, not found", id)
		}
		return nil, errors.WrapUnexpected(err)
	}

	return &cluster, nil
}

func (r *Repo) ClusterByName(ctx context.Context, name string) (*Cluster, error) {
	db := r.getDB()
	defer db.Close()

	var cluster Cluster

	row := db.QueryRowxContext(ctx, "select * from cluster where name = $1", name)
	err := row.StructScan(&cluster)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, NotFoundErrorf("cluster with name: %s, not found", name)
		}
		return nil, errors.WrapUnexpected(err)
	}

	return &cluster, nil
}

func (r *Repo) Clusters(ctx context.Context) (Clusters, error) {
	return r.clusters(ctx)
}

func (r *Repo) clusters(ctx context.Context, whereArgs ...orm.WhereArg) (Clusters, error) {
	db := r.getDB()
	defer db.Close()

	sql, args := orm.CheckWhereArgs("select * from cluster", whereArgs)
	rows, err := db.QueryxContext(ctx, sql, args...)
	if err != nil {
		return nil, errors.WrapUnexpected(err)
	}
	var clusters []Cluster
	for rows.Next() {
		var cluster Cluster
		err = rows.StructScan(&cluster)
		if err != nil {
			return nil, errors.WrapUnexpected(err)
		}
		clusters = append(clusters, cluster)
	}

	return clusters, nil
}

package data

import (
	"context"

	"gitlab.unanet.io/devops/eve/internal/data/orm"
	"gitlab.unanet.io/devops/eve/pkg/errors"
)

type Artifact struct {
	ID               string   `db:"id"`
	Name             string   `db:"name"`
	ArtifactType     string   `db:"artifact_type"`
	ProviderGroup    string   `db:"provider_group"`
	FunctionPointer  string   `db:"function_pointer"`
	CustomerDeployed bool     `db:"customer_deployed"`
	Metadata         JSONText `db:"metadata"`
}

type Artifacts []Artifact

func (r *Repo) ArtifactByID(ctx context.Context, id int) (*Artifact, error) {
	db := r.getDB()
	defer db.Close()

	var artifact Artifact

	row := db.QueryRowxContext(ctx, "select * from artifact where id = $1", id)
	err := row.StructScan(&artifact)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, NotFoundErrorf("artifact with id: %d, not found", id)
		}
		return nil, errors.WrapUnexpected(err)
	}

	return &artifact, nil
}

func (r *Repo) ArtifactByName(ctx context.Context, name string) (*Artifact, error) {
	db := r.getDB()
	defer db.Close()

	var artifact Artifact

	row := db.QueryRowxContext(ctx, "select * from artifact where name = $1", name)
	err := row.StructScan(&artifact)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, NotFoundErrorf("artifact with name: %s, not found", name)
		}
		return nil, errors.WrapUnexpected(err)
	}

	return &artifact, nil
}

func (r *Repo) Artifacts(ctx context.Context) (Artifacts, error) {
	return r.artifacts(ctx)
}

func (r *Repo) artifacts(ctx context.Context, whereArgs ...orm.WhereArg) (Artifacts, error) {
	db := r.getDB()
	defer db.Close()

	sql, args := orm.CheckWhereArgs("select * from artifact", whereArgs)
	rows, err := db.QueryxContext(ctx, sql, args...)
	if err != nil {
		return nil, errors.WrapUnexpected(err)
	}
	var artifacts []Artifact
	for rows.Next() {
		var artifact Artifact
		err = rows.StructScan(&artifact)
		if err != nil {
			return nil, errors.WrapUnexpected(err)
		}
		artifacts = append(artifacts, artifact)
	}

	return artifacts, nil
}

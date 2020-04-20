package data

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"gitlab.unanet.io/devops/eve/internal/data/common"
	"gitlab.unanet.io/devops/eve/pkg/log"
)

type Repo struct {
}

func (r *Repo) getDB() *sqlx.DB {
	db, err := common.GetDBWithTimeout(60 * time.Second)
	if err != nil {
		log.Logger.Panic("Error Connection to the Database", zap.Error(err))
	}

	return db
}

func WhereID(id int) common.WhereArg {
	return func(clause *common.WhereClause) {
		clause.AddClause("id=%s", common.ANDWhereCondition, id)
	}
}

func WhereEnvironmentID(environmentID int) common.WhereArg {
	return func(clause *common.WhereClause) {
		clause.AddClause("environment_id=%s", common.ANDWhereCondition, environmentID)
	}
}

func Where(key string, value interface{}) common.WhereArg {
	return func(clause *common.WhereClause) {
		beginning := fmt.Sprintf("%s=", key)
		clause.AddClause(beginning+"%s", common.ANDWhereCondition, value)
	}
}

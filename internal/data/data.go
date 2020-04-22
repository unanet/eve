package data

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"gitlab.unanet.io/devops/eve/internal/data/orm"
	"gitlab.unanet.io/devops/eve/pkg/log"
)

type Repo struct {
}

func NewRepo() *Repo {
	return &Repo{}
}

func (r *Repo) getDB() *sqlx.DB {
	db, err := GetDBWithTimeout(60 * time.Second)
	if err != nil {
		log.Logger.Panic("Error Connection to the Database", zap.Error(err))
	}

	return db
}

func Where(key string, value interface{}) orm.WhereArg {
	return func(clause *orm.WhereClause) {
		beginning := fmt.Sprintf("%s=", key)
		clause.AddClause(beginning+"%s", orm.ANDWhereCondition, value)
	}
}

func AndWhere(key string, value interface{}) orm.WhereArg {
	return func(clause *orm.WhereClause) {
		beginning := fmt.Sprintf("%s=", key)
		clause.AddClause(beginning+"%s", orm.ANDWhereCondition, value)
	}
}

func OrWhere(key string, value interface{}) orm.WhereArg {
	return func(clause *orm.WhereClause) {
		beginning := fmt.Sprintf("%s=", key)
		clause.AddClause(beginning+"%s", orm.ORWhereCondition, value)
	}
}

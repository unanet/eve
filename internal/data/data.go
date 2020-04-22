package data

import (
	"time"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

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

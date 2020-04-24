package data

import (
	"strings"

	"github.com/golang-migrate/migrate/v4"
	"github.com/jmoiron/sqlx"

	"gitlab.unanet.io/devops/eve/internal/config"
	"gitlab.unanet.io/devops/eve/pkg/log"
)

type Repo struct {
	db *sqlx.DB
}

func NewRepo(db *sqlx.DB) *Repo {
	return &Repo{
		db: db,
	}
}

func MigrateDB(DSN string) error {
	m, err := migrate.New(
		"file://migrations",
		DSN,
	)
	if err != nil {
		return err
	}

	m.Log = NewMigrationLogger(strings.ToLower(config.Values().LogLevel) == "debug")

	err = m.Up()
	if err != nil && err.Error() != "no change" {
		return err
	}

	return nil
}

func NewMigrationLogger(debug bool) MigrationLogger {
	return MigrationLogger{
		debug: debug,
	}
}

type MigrationLogger struct {
	debug bool
}

func (m MigrationLogger) Printf(format string, v ...interface{}) {
	format = strings.TrimSuffix(format, "\n")
	if m.debug {
		log.Logger.Sugar().With("migration", "postgres").Debugf(format, v...)
	} else {
		log.Logger.Sugar().With("migration", "postgres").Infof(format, v...)
	}
}

func (m MigrationLogger) Verbose() bool {
	return m.debug
}

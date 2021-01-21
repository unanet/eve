package data

import (
	"fmt"
	"strings"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"gitlab.unanet.io/devops/go/pkg/errors"
	"gitlab.unanet.io/devops/go/pkg/log"
	"go.uber.org/zap"
)

type Repo struct {
	db *sqlx.DB
}

func NewRepo(db *sqlx.DB) *Repo {
	return &Repo{
		db: db,
	}
}

// ConnectLoop tries to connect to the DB under given DSN using a give driver
// in a loop until connection succeeds. timeout specifies the timeout for the
// loop.
func GetDBWithTimeout(dsn string, timeout time.Duration) (*sqlx.DB, error) {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	timeoutExceeded := time.After(timeout)
	for {
		select {
		case <-timeoutExceeded:
			return nil, fmt.Errorf("db connection failed after %s timeout", timeout)

		case <-ticker.C:
			db, err := sqlx.Connect("postgres", dsn)
			if err == nil {
				return db, nil
			}
			// TODO: This dumps the db password to the logs, we need to scrub this.
			log.Logger.Error("Failed to Connect to DB", zap.String("DSN", dsn))
		}
	}
}

func MigrateDB(DSN, logLevel string) error {
	m, err := migrate.New(
		"file://migrations",
		DSN,
	)
	if err != nil {
		return errors.Wrap(err)
	}

	m.Log = NewMigrationLogger(strings.ToLower(logLevel) == "debug")

	err = m.Up()
	if err != nil && err.Error() != "no change" {
		return errors.Wrap(err)
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

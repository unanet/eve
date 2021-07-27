package data

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"github.com/unanet/go/pkg/errors"
	"github.com/unanet/go/pkg/log"
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

func (r *Repo) deleteByID(ctx context.Context, tableName string, id int) error {
	return r.deleteByIDWithField(ctx, tableName, "id", id)
}

func (r *Repo) deleteByIDWithField(ctx context.Context, tableName string, fieldName string, id int) error {
	query := fmt.Sprintf(`%s = %v`, fieldName, id)
	return r.deleteWithQuery(ctx, tableName, query)
}

func (r *Repo) deleteWithQuery(ctx context.Context, tableName string, query string) error {
	result, err := r.db.ExecContext(ctx, fmt.Sprintf("DELETE FROM %s WHERE %s", tableName, query))
	if err != nil {
		return errors.Wrap(err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err)
	}

	if affected == 0 {
		return errors.NotFoundf(fmt.Sprintf("unable to delete %s using query (%s)", tableName, query))
	}

	return nil
}

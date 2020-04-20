package common

import (
	"fmt"
	"time"

	// adds pg as a sql driver
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"

	"gitlab.unanet.io/devops/eve/internal/config"
	"gitlab.unanet.io/devops/eve/pkg/log"
)

// ConnectLoop tries to connect to the DB under given DSN using a give driver
// in a loop until connection succeeds. timeout specifies the timeout for the
// loop.
func GetDBWithTimeout(timeout time.Duration) (*sqlx.DB, error) {
	dsn := config.Values().DbConnectionString()
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
package main

import (
	_ "github.com/lib/pq"
	"github.com/mattes/migrate"
	_ "github.com/mattes/migrate/database/postgres"
	_ "github.com/mattes/migrate/source/file"
	"go.uber.org/zap"

	"gitlab.unanet.io/devops/eve/internal/config"
	"gitlab.unanet.io/devops/eve/pkg/log"
)

func init() {
	m, err := migrate.New(
		"file://migrations",
		config.Values().MigrationConnectionString(),
	)

	if err != nil {
		log.Logger.Panic("Failed to load the Database Migration Tool.", zap.Error(err))
	}

	err = m.Up()

	if err != nil && err.Error() != "no change" {
		log.Logger.Panic("Migration Failed", zap.Error(err))
	}
}

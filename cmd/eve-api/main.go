package main

import (
	"context"

	_ "github.com/lib/pq"
	"github.com/mattes/migrate"
	_ "github.com/mattes/migrate/database/postgres"
	_ "github.com/mattes/migrate/source/file"

	"gitlab.unanet.io/devops/eve/internal/config"
	"gitlab.unanet.io/devops/eve/pkg/artifactory"
	"gitlab.unanet.io/devops/eve/pkg/log"
)

func init() {
	m, err := migrate.New(
		"file://migrations",
		config.Values.MigrationConnectionString(),
	)

	if err != nil {
		log.Logger.WithField("error", err).Panic("Failed to load the Database Migration Tool.")
	}

	err = m.Up()

	if err != nil && err.Error() != "no change" {
		log.Logger.WithField("error", err).Panic("Migration Failed")
	}
}

func main() {
	client, err := artifactory.NewClient(config.Values.ArtifactoryConfig)
	if err != nil {
		println(err)
	}

	err = client.GetLatestVersion(context.TODO(), "docker-int", "unanet/unanet", "20.2.*")
	if err != nil {
		println(err.Error())
	}

	//internal.StartApi()
}

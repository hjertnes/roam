package commands

import (
	"context"
	"fmt"

	"github.com/hjertnes/roam/configuration"
	"github.com/hjertnes/roam/migration"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rotisserie/eris"
)

func Migrate(path string) error {
	conf, err := configuration.ReadConfigurationFile(fmt.Sprintf("%s/.config/config.yaml", path))
	if err != nil {
		return eris.Wrap(err, "failed to get config")
	}

	ctx := context.Background()
	pxp, err := pgxpool.Connect(ctx, conf.DatabaseConnectionString)
	if err != nil {
		return eris.Wrap(err, "failed to connect to database")
	}

	mig := migration.New(ctx, pxp)
	err = mig.Migrate()
	if err != nil {
		return eris.Wrap(err, "failed to migrate database")
	}

	return nil
}

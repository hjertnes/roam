package commands

import (
	"context"
	"fmt"
	"github.com/hjertnes/roam/configuration"
	"github.com/hjertnes/roam/migration"
	"github.com/hjertnes/roam/utils"
	"github.com/jackc/pgx/v4/pgxpool"
)

func Migrate(path string){
	conf, err := configuration.ReadConfigurationFile(fmt.Sprintf("%s/.config/config.yaml", path))
	utils.ErrorHandler(err)

	ctx := context.Background()
	pxp, err := pgxpool.Connect(ctx, conf.DatabaseConnectionString)
	utils.ErrorHandler(err)

	mig := migration.New(ctx, pxp)
	err = mig.Migrate()
	utils.ErrorHandler(err)

}

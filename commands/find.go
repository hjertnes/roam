package commands

import (
	"context"
	"fmt"
	"github.com/hjertnes/roam/commands/findSearch"
	"github.com/hjertnes/roam/commands/findSelect"
	"github.com/hjertnes/roam/configuration"
	dal2 "github.com/hjertnes/roam/dal"
	"github.com/hjertnes/roam/utils"
	"github.com/jackc/pgx/v4/pgxpool"
)

func Find(path string){
	search, err := findSearch.Run()
	utils.ErrorHandler(err)

	fmt.Println("Loading...")

	conf, err := configuration.ReadConfigurationFile(fmt.Sprintf("%s/.config/config.yaml", path))
	utils.ErrorHandler(err)

	ctx := context.Background()
	pxp, err := pgxpool.Connect(ctx, conf.DatabaseConnectionString)

	dal := dal2.New(ctx, pxp)

	result, err := dal.Find(search)

	utils.ErrorHandler(err)

	paths := make([]string, 0)

	for _, r := range result {
		paths = append(paths, r.Path)
	}

	findSelect.Run(paths)

}
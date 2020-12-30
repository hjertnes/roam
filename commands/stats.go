package commands

import (
	"context"
	"fmt"

	"github.com/hjertnes/roam/configuration"
	dal2 "github.com/hjertnes/roam/dal"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rotisserie/eris"
)

// Stats shows statistics.
func Stats(path string) error {
	conf, err := configuration.ReadConfigurationFile(fmt.Sprintf("%s/.config/config.yaml", path))
	if err != nil {
		return eris.Wrap(err, "failed to get config")
	}

	ctx := context.Background()

	pxp, err := pgxpool.Connect(ctx, conf.DatabaseConnectionString)
	if err != nil {
		return eris.Wrap(err, "could not connect to database")
	}

	dal := dal2.New(path, ctx, pxp)

	all, public, private, links, err := dal.Stats()
	if err != nil {
		return eris.Wrap(err, "failed to get stats")
	}

	fmt.Println("Stats")
	fmt.Printf("All: %v\n", all)
	fmt.Printf("Private: %v\n", private)
	fmt.Printf("Public: %v\n", public)
	fmt.Printf("Links: %v\n", links)

	return nil
}

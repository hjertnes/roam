package commands

import (
	"context"
	"fmt"
	"github.com/hjertnes/roam/configuration"
	dal2 "github.com/hjertnes/roam/dal"
	"github.com/hjertnes/roam/widgets/selectinput"
	"github.com/hjertnes/roam/widgets/textinput"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rotisserie/eris"
)

func getFiles(path string) (string, error){
	search, err := textinput.Run("Search for a note", "Search: ")
	if err != nil{
		return "", eris.Wrap(err, "failed to get a search string from textinput")
	}

	fmt.Println("Loading...")

	conf, err := configuration.ReadConfigurationFile(fmt.Sprintf("%s/.config/config.yaml", path))
	if err != nil {
		return "", eris.Wrap(err,"failed to get config")
	}

	ctx := context.Background()
	pxp, err := pgxpool.Connect(ctx, conf.DatabaseConnectionString)
	if err != nil{
		return "", eris.Wrap(err, "could not connect to database")
	}

	dal := dal2.New(ctx, pxp)

	result, err := dal.Find(search)

	if err != nil{
		return "", eris.Wrap(err, "failed to search for files in database")
	}

	paths := make([]selectinput.Choice, 0)

	for _, r := range result {
		fmt.Println(r.Path)
		paths = append(paths, selectinput.Choice{Title: r.Path, Value: r.Path})
	}

	choice, err := selectinput.Run("Select match", paths)
	if err != nil{
		return "", eris.Wrap(err, "failed to get selection from selectinput")
	}

	return choice.Value, nil


}

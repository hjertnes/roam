package commands

import (
	"context"
	"fmt"
	"github.com/ericaro/frontmatter"
	"github.com/hjertnes/roam/configuration"
	dal2 "github.com/hjertnes/roam/dal"
	"github.com/hjertnes/roam/models"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rotisserie/eris"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func Sync(path string) error{
	conf, err := configuration.ReadConfigurationFile(fmt.Sprintf("%s/.config/config.yaml", path))
	if err != nil{
		return eris.Wrap(err, "failed to get config")
	}


	ctx := context.Background()
	pxp, err := pgxpool.Connect(ctx, conf.DatabaseConnectionString)
	if err != nil {
		return eris.Wrap(err, "could not connect to database")
	}
	dal := dal2.New(ctx, pxp)

	err = dal.Delete()
	if err != nil{
		return eris.Wrap(err, "failed to delete files that don't exist from database")
	}

	err = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if info.Name() == ".DS_Store"{
			return nil
		}
		if strings.Contains(path, "/."){
			return nil
		}
		if info.IsDir(){
			return nil
		}

		data, err := ioutil.ReadFile(path)
		metadata := models.Fm{}
		err = frontmatter.Unmarshal(data, &metadata)
		if err != nil{
			return eris.Wrap(err, "failed to unmarshal frontmatter")
		}
		exist, err := dal.Exists(path)
		if err != nil{
			return eris.Wrap(err, "failed to check if file exists in database")
		}

		if exist{
			err = dal.Update(path, metadata.Title, metadata.Content, metadata.Private)
			if err != nil{
				return eris.Wrap(err, "failed to update in database")
			}
		} else {
			err = dal.Create(path, metadata.Title, metadata.Content, metadata.Private)
			if err != nil{
				return eris.Wrap(err, "failed to create in database")
			}
		}
		return nil
	})

	if err != nil{
		return eris.Wrap(err, "failed to sync")
	}

	return nil
}

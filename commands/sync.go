package commands

import (
	"context"
	"fmt"
	"github.com/ericaro/frontmatter"
	"github.com/hjertnes/roam/configuration"
	dal2 "github.com/hjertnes/roam/dal"
	"github.com/hjertnes/roam/models"
	"github.com/hjertnes/roam/utils"
	"github.com/jackc/pgx/v4/pgxpool"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func Sync(path string){
	conf, err := configuration.ReadConfigurationFile(fmt.Sprintf("%s/.config/config.yaml", path))
	utils.ErrorHandler(err)

	ctx := context.Background()
	pxp, err := pgxpool.Connect(ctx, conf.DatabaseConnectionString)

	dal := dal2.New(ctx, pxp)

	// TODO do partial updates based on LastModified, opened_at updated_at
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
			return err
		}
		exist, err := dal.Exists(path)
		if err != nil{
			return err
		}

		if exist{
			err = dal.Update(path, metadata.Title, metadata.Content)
			if err != nil{
				return err
			}
		} else {
			err = dal.Create(path, metadata.Title, metadata.Content)
			if err != nil{
				return err
			}
		}

		return nil
	})

	if err != nil{
		panic(err)
	}
}

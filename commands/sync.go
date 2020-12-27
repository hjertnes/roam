package commands

import (
	"context"
	"fmt"
	"github.com/ericaro/frontmatter"
	"github.com/hjertnes/roam/configuration"
	dal2 "github.com/hjertnes/roam/dal"
	"github.com/hjertnes/roam/models"
	"github.com/hjertnes/utils"
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
		metadata, err := readfile(path)
		if err != nil{
			return eris.Wrap(err, "could not read file")
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
	files, err := dal.GetFiles()

	for _, file := range files {
		if !utils.FileExist(file.Path){
			continue
		}

		metadata, err := readfile(file.Path)

		if err != nil{
			continue
		}

		links := noteLinkRegexp.FindAllString(metadata.Content, -1)
		currentInDatabaseLinks, err := dal.GetCurrentLinks(file.Id)
		currentLinks := make([]string, 0)

		for _, link := range links{
			clean := cleanLink(link)

			var path = ""
			if strings.HasPrefix(clean, "/"){
				exist1, err := dal.Exists(fmt.Sprintf("%s%s.md", path, clean))
				if err != nil {
					return eris.Wrap(err, "failed to check if link exists")
				}

				exist2, err := dal.Exists(fmt.Sprintf("%s%s/index.md", path, clean))
				if err != nil {
					return eris.Wrap(err, "failed to check if link exists")
				}

				if !exist1 && !exist2{
					continue
				}

				if exist1{
					path = fmt.Sprintf("%s%s.md", path, clean)
				}

				if exist2{
					path = fmt.Sprintf("%s%s/index.md", path, clean)
				}
			}else{
				path = clean
			}

			matches, err := dal.FindExact(path)
			if err != nil{
				return eris.Wrap(err, "failed to search for link")
			}

			if len(matches) == 0 {
				continue
			}

			if len(matches) > 1 {
				continue
			}

			err = dal.AddLink(file.Id, matches[0].Id)
			if err != nil {
				return eris.Wrap(err, "failed to add link")
			}

			currentLinks = append(currentLinks, matches[0].Id)
		}

		for _, l := range currentInDatabaseLinks{
			if !contains(l, currentLinks){
				err := dal.DeleteLink(file.Id, l)
				if err != nil {
					return eris.Wrap(err, "failed to delete link")
				}
			}
		}
	}

	if err != nil{
		return eris.Wrap(err, "failed to sync")
	}

	return nil
}

func contains(id string, files []string) bool{
	c := false

	for _, f := range files{
		if f == id{
			c = true
			break
		}
	}

	return c
}

func readfile(path string) (*models.Fm, error){
	data, err := ioutil.ReadFile(path)
	metadata := models.Fm{}
	err = frontmatter.Unmarshal(data, &metadata)
	if err != nil{
		return nil, eris.Wrap(err, "failed to unmarshal frontmatter")
	}

	return &metadata, nil
}
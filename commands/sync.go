package commands

import (
	"context"
	"fmt"
	"github.com/hjertnes/roam/errs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/ericaro/frontmatter"
	"github.com/hjertnes/roam/configuration"
	dal2 "github.com/hjertnes/roam/dal"
	"github.com/hjertnes/roam/models"
	"github.com/hjertnes/utils"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rotisserie/eris"
)

func Sync(path string) error {
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

	err = dal.Delete()
	if err != nil {
		return eris.Wrap(err, "failed to delete files that don't exist from database")
	}

	err = filepath.Walk(path, func(path string, info os.FileInfo, errr error) error {
		if errr != nil {
			return eris.Wrap(errr, "unknown problems parsing folder")
		}

		if info.Name() == ".DS_Store" {
			return nil
		}
		if strings.Contains(path, "/.") {
			return nil
		}
		if info.IsDir() {
			return nil
		}
		metadata, err := readfile(path)
		if err != nil {
			return eris.Wrap(err, "could not read file")
		}
		exist, err := dal.Exists(path)
		if err != nil {
			return eris.Wrap(err, "failed to check if file exists in database")
		}

		if exist {
			err = dal.Update(path, metadata.Title, metadata.Content, metadata.Private)
			if err != nil {
				return eris.Wrap(err, "failed to update in database")
			}
		} else {
			err = dal.Create(path, metadata.Title, metadata.Content, metadata.Private)
			if err != nil {
				return eris.Wrap(err, "failed to create in database")
			}
		}

		return nil
	})

	if err != nil {
		return eris.Wrap(err, "failed to process files")
	}

	files, err := dal.GetFiles()
	if err != nil {
		return eris.Wrap(err, "failed to get list of files")
	}

	for _, file := range files {
		if !utils.FileExist(file.Path) {
			continue
		}

		metadata, err := readfile(file.Path)
		if err != nil {
			continue
		}

		links := noteLinkRegexp.FindAllString(metadata.Content, -1)

		currentInDatabaseLinks, err := dal.GetCurrentLinks(file.ID)
		if err != nil {
			return eris.Wrap(err, "failed to get current links")
		}

		currentLinks := make([]string, 0)

		for _, link := range links {
			clean := cleanLink(link)

			filename := clean

			if strings.HasPrefix(clean, "/") {
				exist1, err := dal.Exists(fmt.Sprintf("%s%s.md", path, clean))
				if err != nil {
					return eris.Wrap(err, "failed to check if link exists")
				}

				exist2, err := dal.Exists(fmt.Sprintf("%s%s/index.md", path, clean))
				if err != nil {
					return eris.Wrap(err, "failed to check if link exists")
				}

				if exist1 {
					filename = fmt.Sprintf("%s%s.md", path, clean)
				} else if exist2 {
					filename = fmt.Sprintf("%s%s/index.md", path, clean)
				} else {
					return eris.Wrap(errs.ErrNotFound, "not found")
				}
			}

			matches, err := dal.FindExact(filename)
			if err != nil {
				return eris.Wrap(err, "failed to search for link")
			}

			if len(matches) == 0 {
				continue
			}

			if len(matches) > 1 {
				continue
			}

			err = dal.AddLink(file.ID, matches[0].ID)
			if err != nil {
				return eris.Wrap(err, "failed to add link")
			}

			currentLinks = append(currentLinks, matches[0].ID)
		}

		for _, l := range currentInDatabaseLinks {
			if !contains(l, currentLinks) {
				err := dal.DeleteLink(file.ID, l)
				if err != nil {
					return eris.Wrap(err, "failed to delete link")
				}
			}
		}
	}

	if err != nil {
		return eris.Wrap(err, "failed to sync")
	}

	return nil
}

func contains(id string, files []string) bool {
	c := false

	for _, f := range files {
		if f == id {
			c = true

			break
		}
	}

	return c
}

func readfile(path string) (*models.Frontmatter, error) {
	data, err := ioutil.ReadFile(filepath.Clean(path))
	if err != nil {
		return nil, eris.Wrap(err, "failed to read file")
	}

	metadata := models.Frontmatter{}

	err = frontmatter.Unmarshal(data, &metadata)
	if err != nil {
		return nil, eris.Wrap(err, "failed to unmarshal frontmatter")
	}

	return &metadata, nil
}

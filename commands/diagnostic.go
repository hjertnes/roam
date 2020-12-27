package commands

import (
	"context"
	"fmt"
	"github.com/hjertnes/roam/configuration"
	dal2 "github.com/hjertnes/roam/dal"
	"github.com/hjertnes/utils"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rotisserie/eris"
	"regexp"
	"strings"
)

func Diagnostic(path string) error{
	conf, err := configuration.ReadConfigurationFile(fmt.Sprintf("%s/.config/config.yaml", path))
	if err != nil {
		return eris.Wrap(err, "failed to get config")
	}

	ctx := context.Background()
	pxp, err := pgxpool.Connect(ctx, conf.DatabaseConnectionString)
	if err != nil {
		return eris.Wrap(err, "could not connect to database")
	}

	dal := dal2.New(ctx, pxp)

	files, err := dal.GetFiles()

	for _, file := range files {
		if !utils.FileExist(file.Path){
			fmt.Printf("%s: doesn't exist\n", file)
			continue
		}

		metadata, err := readfile(file.Path)

		if err != nil{
			fmt.Printf("%s could not read file\n", file)
			fmt.Println("most likely invalid front matter")
			continue
		}

		links := noteLinkRegexp.FindAllString(metadata.Content, -1)

		for _, link := range links{
			clean := cleanLink(link)
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
					fmt.Printf("%s no matches for link %s", file, clean)
					continue
				}
			}else{
				matches, err := dal.FindExact(clean)
				if err != nil{
					return eris.Wrap(err, "failed to search for link")
				}

				if len(matches) == 0 {
					fmt.Printf("%s no matches for link %s\n", file, clean)
					continue
				}

				if len(matches) > 1 {
					fmt.Printf("%s more than one match for link %s\n", file, clean)
					continue
				}
			}

		}
	}

	return nil
}

func cleanLink(input string) string{
	return strings.ReplaceAll(strings.ReplaceAll(input, "[[", ""), "]]", "")
}

var noteLinkRegexp =  regexp.MustCompile(`\[\[([\d\w\s./]+)]]`)

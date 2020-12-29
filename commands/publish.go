package commands

import (
	"bytes"
	"context"
	"fmt"
	"github.com/ericaro/frontmatter"
	"github.com/hjertnes/roam/configuration"
	dal2 "github.com/hjertnes/roam/dal"
	"github.com/hjertnes/roam/errs"
	"github.com/hjertnes/roam/models"
	utilslib "github.com/hjertnes/utils"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rotisserie/eris"
	"github.com/yuin/goldmark"
	"io/ioutil"
	"os"
	"strings"
)

func Publish(path, to string) error{
	outputDir := to
	if to == ""{
		outputDir = "./output"
	}

	if utilslib.FileExist(outputDir){
		err := os.RemoveAll(outputDir)
		if err != nil{
			return eris.Wrap(err, "failed to delete output dir")
		}
	}

	err := os.MkdirAll(outputDir, os.ModePerm)
	if err != nil{
		return eris.Wrap(err, "failed to create output dir")
	}

	conf, err := configuration.ReadConfigurationFile(fmt.Sprintf("%s/.config/config.yaml", path))
	if err != nil {
		return eris.Wrap(err, "failed to get config")
	}

	ctx := context.Background()

	pxp, err := pgxpool.Connect(ctx, conf.DatabaseConnectionString)
	if err != nil {
		return eris.Wrap(err, "failed to connect to database")
	}

	dal := dal2.New(ctx, pxp)

	files, err := dal.GetFiles()
	if err != nil{
		return eris.Wrap(err, "failed to get list of files")
	}

	index := make([]string, 0)

	template, err := ioutil.ReadFile(fmt.Sprintf("%s/.config/publish/template.html", path))
	if err != nil{
		return eris.Wrap(err, "failed to read template")
	}

	for _, file := range files{
		data, err := ioutil.ReadFile(file.Path)
		if err != nil{
			return eris.Wrap(err, "failed to read file")
		}

		metadata := models.Frontmatter{}

		err = frontmatter.Unmarshal(data, &metadata)
		if err != nil {
			return eris.Wrap(err, "could not unmarkshal frontmatter")
		}

		links := noteLinkRegexp.FindAllString(metadata.Content, -1)

		for _, link := range links {
			clean := cleanLink(link)

			if strings.HasPrefix(clean, "/") {
				exist1, err := dal.Exists(fmt.Sprintf("%s%s.md", path, clean))
				if err != nil {
					return eris.Wrap(err, "failed to check if link exists")
				}

				exist2, err := dal.Exists(fmt.Sprintf("%s%s/index.md", path, clean))
				if err != nil {
					return eris.Wrap(err, "failed to check if link exists")
				}

				if exist1{
					clean = fmt.Sprintf("%s%s.md", path, clean)
				} else if exist2{
					clean = fmt.Sprintf("%s%s.md", path, clean)
				} else {
					return eris.Wrap(errs.ErrNotFound, "no match")
				}
			} else {
				matches, err := dal.FindExact(clean)
				if err != nil {
					return eris.Wrap(err, "failed to search for link")
				}

				if len(matches) == 0 {
					return eris.Wrap(errs.ErrNotFound, "no match for link")
				}

				if len(matches) > 1 {
					return eris.Wrap(errs.ErrNotFound, "more than one match for link")
				}
				clean = strings.ReplaceAll(matches[0].Path, path, outputDir)
			}



			metadata.Content = strings.ReplaceAll(metadata.Content, link, fmt.Sprintf("[%s](%s)", cleanLink(link), fixUrl(strings.ReplaceAll(clean , outputDir, ""))))
		}

		var buf bytes.Buffer
		if err := goldmark.Convert([]byte(metadata.Content), &buf); err != nil {
			return eris.Wrap(err, "failed to build markdown")
		}

		filePath := strings.ReplaceAll(file.Path, path, outputDir)

		folderPath, filename := destructPath(filePath)

		if !utilslib.FileExist(folderPath){
			err = os.MkdirAll(folderPath, os.ModePerm)
			if err != nil{
				return eris.Wrap(err, "failed to create folder")
			}
		}

		var fullFilePath = fmt.Sprintf("%s/%s",folderPath, filename)

		out := string(template)
		out = strings.ReplaceAll(out, "$$TITLE$$", metadata.Title)
		out = strings.ReplaceAll(out, "$$TEXT$$", buf.String())

		backlinks := make([]string, 0)

		backlinks = append(backlinks, "## Backlinks")

		bl, err := dal.GetBacklinks(file.ID)
		if err != nil{
			return eris.Wrap(err, "coult not get backlinks")
		}

		for _, l := range bl{
			lt := fmt.Sprintf("- [%s](%s)", l.Title, fixUrl(strings.ReplaceAll(l.Path , outputDir, "")))
			backlinks = append(backlinks, lt)
		}

		buf = bytes.Buffer{}

		if err := goldmark.Convert([]byte(strings.Join(backlinks, "\n")), &buf); err != nil {
			return eris.Wrap(err, "failed to build markdown")
		}

		if len(backlinks) > 1{
			out = strings.ReplaceAll(out, "$$BACKLINKS$$", buf.String())
		}else {
			out = strings.ReplaceAll(out, "$$BACKLINKS$$", "")
		}


		err = ioutil.WriteFile(fullFilePath, []byte(out), os.ModePerm)
		if err != nil{
			return eris.Wrap(err, "failed to write file")
		}

		index = append(index, fmt.Sprintf("- [%s](%s)", strings.ReplaceAll(fullFilePath, outputDir, ""), fixUrl(strings.ReplaceAll(fullFilePath, outputDir, ""))))
	}

	var buf bytes.Buffer
	if err := goldmark.Convert([]byte(strings.Join(index, "\n")), &buf); err != nil {
		return eris.Wrap(err, "failed to build markdown")
	}

	out := string(template)
	out = strings.ReplaceAll(out, "$$TITLE$$", "Index")
	out = strings.ReplaceAll(out, "$$TEXT$$", buf.String())
	out = strings.ReplaceAll(out, "$$BACKLINKS$$", "")

	err = ioutil.WriteFile(fmt.Sprintf("%s/index.html", outputDir), []byte(out), os.ModePerm)
	if err != nil{
		return eris.Wrap(err, "failed to write file")
	}

	data, err := ioutil.ReadFile(fmt.Sprintf("%s/.config/publish/style.css", path))
	if err != nil{
		return eris.Wrap(err, "failed to read style")
	}
	err = ioutil.WriteFile(fmt.Sprintf("%s/style.css", outputDir), data, os.ModePerm)
	if err != nil{
		return eris.Wrap(err, "failed to write style")
	}

	return nil
}

func destructPath(path string)(string, string){
	elems := strings.Split(path, "/")

	folderPath := make([]string, 0)
	filename := ""

	lastElem := len(elems)-1

	for i, e := range elems{
		if i == lastElem{
			filename=e
		} else {
			folderPath = append(folderPath, e)
		}
	}

	return strings.Join(folderPath, "/"), strings.ReplaceAll(filename, ".md", ".html")
}

func fixUrl(input string) string{
	output := strings.ReplaceAll(input, " ", "%20")
	output = strings.ReplaceAll(output, ".md", ".html")
	return output
}
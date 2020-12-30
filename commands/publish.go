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
	"path/filepath"
	"strings"
)

func Publish(path, to string, excludePrivate bool) error{
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

	dal := dal2.New(path, ctx, pxp)

	files, err := dal.GetFiles()
	if err != nil{
		return eris.Wrap(err, "failed to get list of files")
	}

	template, err := ioutil.ReadFile(fmt.Sprintf("%s/.config/publish/template.html", path))
	if err != nil{
		return eris.Wrap(err, "failed to read template")
	}

	for _, file := range files{
		if excludePrivate && file.Private{
			continue
		}
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

	}



	root, err := dal.GetRootFolder()
	if err != nil{
		return eris.Wrap(err, "failed to get root folder")
	}

	output := make([]string, 0)
	output, err = printAndIterate(excludePrivate, path, dal, root, output)
	if err != nil{
		return eris.Wrap(err, "failed")
	}


	out := string(template)
	out = strings.ReplaceAll(out, "$$TITLE$$", "Index")
	out = strings.ReplaceAll(out, "$$TEXT$$", strings.Join(output, "\n"))
	out = strings.ReplaceAll(out, "$$BACKLINKS$$", "")

	err = ioutil.WriteFile(fmt.Sprintf("%s/index.html", outputDir), []byte(out), os.ModePerm)
	if err != nil{
		return eris.Wrap(err, "failed to write file")
	}

	p := path
	err = filepath.Walk(fmt.Sprintf("%s/.config/publish", path), func(path string, info os.FileInfo, errr error) error {
		if info.IsDir(){
			return nil
		}

		if strings.HasSuffix(path, ".css") || strings.HasSuffix(path,".js"){
			to := strings.ReplaceAll(path, fmt.Sprintf("%s/.config/publish", p), outputDir)
			data, err := ioutil.ReadFile(path)
			if err != nil{
				return eris.Wrap(err, "failed to read file")
			}

			err = ioutil.WriteFile(to, data, os.ModePerm)
			if err != nil{
				return eris.Wrap(err, "failed to write file")
			}
		}

		return nil
	})





	return nil
}

func getLast(path string) string{
	elems := strings.Split(path, "/")
	return elems[len(elems)-1]
}

func printAndIterate(excludePrivate bool, path string, dal *dal2.Dal, folder *models.Folder, o []string) ([]string, error) {
	output := o


	files, err := dal.GetFolderFiles(folder.ID)
	if err != nil{
		return output, eris.Wrap(err, "could not get files")
	}
	folders, err := dal.GetSubFolders(folder.ID)
	if err != nil{
		return output, eris.Wrap(err, "could not get folders")
	}

	for _, f := range files{
		if excludePrivate && f.Private{
			continue
		}
		if strings.HasSuffix(f.Path, "index.md"){
			if folder.Path != path {
				output = append(output, fmt.Sprintf(`<li><a href="%s">%s</a></li>`, strings.ReplaceAll(folder.Path, path, ""), getLast(strings.ReplaceAll(folder.Path, path, ""))))
			}
		}
	}

	output = append(output, "<ul>")

	for _, f := range files{
		if excludePrivate && f.Private{
			continue
		}
		if strings.HasSuffix(f.Path, "index.md"){
			continue
		}
		output = append(output, fmt.Sprintf(`<li><a href="%s">%s</a></li>`, strings.ReplaceAll(strings.ReplaceAll(f.Path, path, ""), ".md", ".html"), f.Title))
	}

	for _, f := range folders{
		output, err = printAndIterate(excludePrivate, path, dal, &f, output)
		if err != nil{
			return output, eris.Wrap(err, "failed to iterate over folder")
		}
	}
	output = append(output, "</ul>")
	return output, nil
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

package publish

import (
	"bytes"

	"fmt"
	"github.com/ericaro/frontmatter"

	dal2 "github.com/hjertnes/roam/dal"
	"github.com/hjertnes/roam/errs"
	"github.com/hjertnes/roam/models"
	"github.com/hjertnes/roam/state"
	"github.com/hjertnes/roam/utils"
	utilslib "github.com/hjertnes/utils"

	"github.com/rotisserie/eris"
	"github.com/yuin/goldmark"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func Run(path, to string, excludePrivate bool) error{
	s, err := state.New(path)
	if err != nil{
		return eris.Wrap(err, "failed to create state")
	}

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

	err = os.MkdirAll(outputDir, os.ModePerm)
	if err != nil{
		return eris.Wrap(err, "failed to create output dir")
	}



	files, err := s.Dal.GetFiles()
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



		links := utils.NoteLinkRegexp.FindAllString(metadata.Content, -1)

		for _, link := range links {
			clean := utils.CleanLink(link)

			if strings.HasPrefix(clean, "/") {
				exist1, err := s.Dal.FileExists(fmt.Sprintf("%s%s.md", path, clean))
				if err != nil {
					return eris.Wrap(err, "failed to check if link exists")
				}

				exist2, err := s.Dal.FileExists(fmt.Sprintf("%s%s/index.md", path, clean))
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
				matches, err := s.Dal.FindFileExact(clean)
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



			metadata.Content = strings.ReplaceAll(metadata.Content, link, fmt.Sprintf("[%s](%s)", utils.CleanLink(link), utils.FixUrl(strings.ReplaceAll(clean , outputDir, ""))))
		}

		var buf bytes.Buffer
		if err := goldmark.Convert([]byte(metadata.Content), &buf); err != nil {
			return eris.Wrap(err, "failed to build markdown")
		}

		filePath := strings.ReplaceAll(file.Path, path, outputDir)

		folderPath, filename := utils.DestructPath(filePath)

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

		bl, err := s.Dal.GetBacklinks(file.ID)
		if err != nil{
			return eris.Wrap(err, "coult not get backlinks")
		}

		for _, l := range bl{
			lt := fmt.Sprintf("- [%s](%s)", l.Title, utils.FixUrl(strings.ReplaceAll(l.Path , outputDir, "")))
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



	root, err := s.Dal.GetRootFolder()
	if err != nil{
		return eris.Wrap(err, "failed to get root folder")
	}

	output := make([]string, 0)
	output, err = buildIndex(excludePrivate, path, s.Dal, root, output)
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

func buildIndex(excludePrivate bool, path string, dal dal2.Dal, folder *models.Folder, o []string) ([]string, error) {
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
				output = append(output, fmt.Sprintf(`<li><a href="%s">%s</a></li>`, strings.ReplaceAll(folder.Path, path, ""), utils.GetLast(strings.ReplaceAll(folder.Path, path, ""))))
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
		output, err = buildIndex(excludePrivate, path, dal, &f, output)
		if err != nil{
			return output, eris.Wrap(err, "failed to iterate over folder")
		}
	}
	output = append(output, "</ul>")
	return output, nil
}

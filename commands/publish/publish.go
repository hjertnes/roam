// Package publish builds a website from your roam
package publish

import (
	"bytes"
	"fmt"
	"github.com/hjertnes/roam/utils/pathutils"
	"github.com/rotisserie/eris"
	"github.com/yuin/goldmark"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/ericaro/frontmatter"
	"github.com/hjertnes/roam/constants"
	"github.com/hjertnes/roam/errs"
	"github.com/hjertnes/roam/models"
	"github.com/hjertnes/roam/state"
	"github.com/hjertnes/roam/utils"
	spinner2 "github.com/hjertnes/roam/widgets/spinner"
	utilslib "github.com/hjertnes/utils"
)

type Publish struct {
	state *state.State
	to string
	excludePrivate bool
}

// New is the constructor.
func New(path string, args []string) (*Publish, error) {
	excludePrivate := true
	to := ""

	for _, a := range args {
		if a == "--include-private" {
			excludePrivate = false
		} else {
			to = a
		}
	}

	outputDir := to
	if to == "" {
		outputDir = "./output"
	}

	s, err := state.New(path, args)
	if err != nil {
		return nil, eris.Wrap(err, "could not create state")
	}


	return &Publish{
		excludePrivate: excludePrivate,
		to: outputDir,
		state: s,
	}, nil
}

func (p *Publish) processLink(link string, metadata *models.Frontmatter) (string, error) {
	clean := utils.CleanLink(link)

	if strings.HasPrefix(clean, "/") {
		exist1, err := p.state.Dal.FileExists(fmt.Sprintf("%s%s.md", p.state.Path, clean))
		if err != nil {
			return "", eris.Wrap(err, "failed to check if link exists")
		}

		exist2, err := p.state.Dal.FileExists(fmt.Sprintf("%s%s/index.md", p.state.Path, clean))
		if err != nil {
			return "", eris.Wrap(err, "failed to check if link exists")
		}

		if exist1 {
			clean = fmt.Sprintf("%s%s.md", p.state.Path, clean)
		} else if exist2 {
			clean = fmt.Sprintf("%s%s.md", p.state.Path, clean)
		} else {
			return "", eris.Wrap(errs.ErrNotFound, "no match")
		}
	} else {
		matches, err := p.state.Dal.FindFileExact(clean)
		if err != nil {
			return "", eris.Wrap(err, "failed to search for link")
		}

		if len(matches) == 0 {
			return "", eris.Wrap(errs.ErrNotFound, "no match for link")
		}

		if len(matches) > 1 {
			return "", eris.Wrap(errs.ErrNotFound, "more than one match for link")
		}

		clean = strings.ReplaceAll(matches[0].Path, p.state.Path, p.to)
	}

	return strings.ReplaceAll(metadata.Content, link, fmt.Sprintf("[%s](%s)", utils.CleanLink(link), fixURL(strings.ReplaceAll(clean, p.to, "")))), nil
}

func (p *Publish) processFile(file *models.File, template []byte) error {
	if p.excludePrivate && file.Private {
		return nil
	}

	data, err := ioutil.ReadFile(file.Path)
	if err != nil {
		return eris.Wrap(err, "failed to read file")
	}

	metadata := models.Frontmatter{}

	err = frontmatter.Unmarshal(data, &metadata)
	if err != nil {
		return eris.Wrap(err, "could not unmarkshal frontmatter")
	}

	links := constants.NoteLinkRegexp.FindAllString(metadata.Content, -1)

	for _, link := range links {
		l, err := p.processLink(link, &metadata)
		if err != nil{
			return eris.Wrap(err, "failed to process link")
		}

		metadata.Content = l
	}

	err = p.writeFile(file, &metadata, template)
	if err != nil{
		return eris.Wrap(err, "failed to write file")
	}

	return nil
}

func (p *Publish) writeFile(file *models.File, metadata *models.Frontmatter, template []byte) error{
	var buf bytes.Buffer

	if err := goldmark.Convert([]byte(metadata.Content), &buf); err != nil {
		return eris.Wrap(err, "failed to build markdown")
	}

	filePath := strings.ReplaceAll(file.Path, p.state.Path, p.to)

	folderPath, filename := pathutils.New(filePath).Destruct()

	if !utilslib.FileExist(folderPath) {
		err := os.MkdirAll(folderPath, constants.FolderPermission)
		if err != nil {
			return eris.Wrap(err, "failed to create folder")
		}
	}

	var fullFilePath = fmt.Sprintf("%s/%s", folderPath, strings.ReplaceAll(filename, ".md", ".html"))

	out := string(template)
	out = strings.ReplaceAll(out, "$$TITLE$$", metadata.Title)
	out = strings.ReplaceAll(out, "$$TEXT$$", buf.String())

	backlinks := make([]string, 0)

	backlinks = append(backlinks, "## Backlinks")

	bl, err := p.state.Dal.GetBacklinks(file.ID, !p.excludePrivate)
	if err != nil {
		return eris.Wrap(err, "coult not get backlinks")
	}

	for _, l := range bl {
		lt := fmt.Sprintf("- [%s](%s)", l.Title, fixURL(strings.ReplaceAll(l.Path, p.to, "")))
		backlinks = append(backlinks, lt)
	}

	buf = bytes.Buffer{}

	if err := goldmark.Convert([]byte(strings.Join(backlinks, "\n")), &buf); err != nil {
		return eris.Wrap(err, "failed to build markdown")
	}

	if len(backlinks) > 1 {
		out = strings.ReplaceAll(out, "$$BACKLINKS$$", buf.String())
	} else {
		out = strings.ReplaceAll(out, "$$BACKLINKS$$", "")
	}

	err = ioutil.WriteFile(fullFilePath, []byte(out), constants.FilePermission)
	if err != nil {
		return eris.Wrap(err, "failed to write file")
	}

	return nil
}

func (p *Publish) initOutputDir() error{
	if utilslib.FileExist(p.to) {
		err := os.RemoveAll(p.to)
		if err != nil {
			return eris.Wrap(err, "failed to delete output dir")
		}
	}

	err := os.MkdirAll(p.to, constants.FolderPermission)
	if err != nil {
		return eris.Wrap(err, "failed to create output dir")
	}

	return nil
}

// Run is the entrypoint.
func (p *Publish) Run() error {
	spinner, err := spinner2.Run("")
	if err != nil {
		return eris.Wrap(err, "failed to create spinner")
	}

	err = spinner.Start()
	if err != nil {
		return eris.Wrap(err, "failed to start spinner")
	}

	err = p.initOutputDir()
	if err != nil{
		return eris.Wrap(err, "failed to init output dir")
	}

	files, err := p.state.Dal.GetFiles()
	if err != nil {
		return eris.Wrap(err, "failed to get list of files")
	}

	template, err := ioutil.ReadFile(fmt.Sprintf("%s/.config/publish/template.html", p.state.Path))
	if err != nil {
		return eris.Wrap(err, "failed to read template")
	}

	for _, file := range files {
		err = p.processFile(&file, template)
		if err != nil{
			return eris.Wrap(err, "failed to process file")
		}
	}

	err = p.createIndex(template)
	if err != nil{
		return eris.Wrap(err, "failed to create index")
	}

	err = spinner.Stop()
	if err != nil {
		return eris.Wrap(err, "failed to stop spinner")
	}

	return nil
}

func (p *Publish) createIndex(template []byte) error {
	root, err := p.state.Dal.GetRootFolder()
	if err != nil {
		return eris.Wrap(err, "failed to get root folder")
	}

	output := make([]string, 0)
	output, err = p.buildIndex(root, output)
	if err != nil {
		return eris.Wrap(err, "failed")
	}
	out := string(template)
	out = strings.ReplaceAll(out, "$$TITLE$$", "Index")
	out = strings.ReplaceAll(out, "$$TEXT$$", strings.Join(output, "\n"))
	out = strings.ReplaceAll(out, "$$BACKLINKS$$", "")

	err = ioutil.WriteFile(fmt.Sprintf("%s/index.html", p.to), []byte(out), constants.FilePermission)
	if err != nil {
		return eris.Wrap(err, "failed to write file")
	}

	err = filepath.Walk(fmt.Sprintf("%s/.config/publish", p.state.Path), func(path string, info os.FileInfo, errr error) error {
		if info.IsDir() {
			return nil
		}

		if contains(p.state.Conf.Publish.FilesToCopy, path){
			to := strings.ReplaceAll(path, fmt.Sprintf("%s/.config/publish", p), p.state.Path)

			data, err := ioutil.ReadFile(path)
			if err != nil {
				return eris.Wrap(err, "failed to read file")
			}

			err = ioutil.WriteFile(to, data, constants.FilePermission)
			if err != nil {
				return eris.Wrap(err, "failed to write file")
			}
		}

		return nil
	})

	return nil
}

func (p *Publish) buildIndex(folder *models.Folder, o []string) ([]string, error) {
	output := o

	files, err := p.state.Dal.GetFolderFiles(folder.ID)
	if err != nil {
		return output, eris.Wrap(err, "could not get files")
	}

	folders, err := p.state.Dal.GetSubFolders(folder.ID)
	if err != nil {
		return output, eris.Wrap(err, "could not get folders")
	}

	for _, f := range files {
		if p.excludePrivate && f.Private {
			continue
		}

		if strings.HasSuffix(f.Path, "index.md") {
			if folder.Path != p.state.Path {

				output = append(output, fmt.Sprintf(`<li><a href="%s">%s</a></li>`, strings.ReplaceAll(folder.Path, p.state.Path, ""), pathutils.New(strings.ReplaceAll(folder.Path, p.state.Path, "")).GetLast()))
			}
		}
	}

	output = append(output, "<ul>")

	for _, f := range files {
		if p.excludePrivate && f.Private {
			continue
		}

		if strings.HasSuffix(f.Path, "index.md") {
			continue
		}

		output = append(output, fmt.Sprintf(`<li><a href="%s">%s</a></li>`, strings.ReplaceAll(strings.ReplaceAll(f.Path, p.state.Path, ""), ".md", ".html"), f.Title))
	}

	for _, f := range folders {
		output, err = p.buildIndex(&f, output)
		if err != nil {
			return output, eris.Wrap(err, "failed to iterate over folder")
		}
	}

	output = append(output, "</ul>")

	return output, nil
}

func fixURL(input string) string {
	output := strings.ReplaceAll(input, " ", "%20")

	output = strings.ReplaceAll(output, ".md", ".html")

	return output
}

func contains(patterns []string, filename string) bool{
	lowerCaseFilename := strings.ToLower(filename)
	for _, i := range patterns{
		if strings.HasSuffix(lowerCaseFilename, i){
			return true
		}
	}

	return false
}
package bulkimport

import (
	"fmt"
	"github.com/ericaro/frontmatter"
	create2 "github.com/hjertnes/roam/commands/create"
	"github.com/hjertnes/roam/commands/help"
	"github.com/hjertnes/roam/constants"
	"github.com/hjertnes/roam/errs"
	"github.com/hjertnes/roam/models"
	"github.com/hjertnes/roam/state"
	"github.com/hjertnes/roam/utils/pathutils"
	utilslib "github.com/hjertnes/utils"
	"github.com/rotisserie/eris"
	"io/ioutil"
	"os"
	"strings"
)

type Import struct {
	state *state.State
	create *create2.Create
}


// New is the constructor.
func New(path string, args []string) (*Import, error) {
	s, err := state.New(path, args)
	if err != nil {
		return nil, eris.Wrap(err, "could not create state")
	}

	c := create2.New(s)
	return &Import{
		create: c,
		state: s,
	}, nil
}

// RunImport runs a import.
func (c *Import) Run() error {
	var file string

	dryRun := false

	for i := range c.state.Arguments {
		if i <= 1 {
			continue
		}

		if c.state.Arguments[i] == "--dry" {
			dryRun = true
		} else {
			file = c.state.Arguments[i]
		}
	}

	if file == "" {
		help.Run([]string{})
		return nil
	}

	err := c.doImport(file, dryRun)
	if err != nil {
		return eris.Wrap(err, "failed to import")
	}

	return nil
}

func (c *Import) doImport(file string, dryRun bool) error {
	importFile, err := ioutil.ReadFile(file) // #nosec
	if err != nil {
		return eris.Wrap(err, "Failed to read import file")
	}

	importContent := strings.Split(string(importFile), "\n")

	fileContent := make([]string, 0)
	sepCounter := 0

	counter := 0

	for _, line := range importContent {
		if line == "---" {
			sepCounter++
		}

		if sepCounter == constants.Three {
			err = c.writeImport(fileContent, dryRun)
			if err != nil {
				return eris.Wrap(err, "failed to import")
			}

			fileContent = make([]string, 0)

			sepCounter = 1

			counter++
		}

		fileContent = append(fileContent, line)
	}

	err = c.writeImport(fileContent, dryRun)
	if err != nil {
		return eris.Wrap(err, "failed to import")
	}

	counter++

	fmt.Printf("Imported %v notes\n", counter)

	return nil
}

func readfileImport(data string) (*models.ImportFrontmatter, error) {
	metadata := models.ImportFrontmatter{}

	err := frontmatter.Unmarshal([]byte(data), &metadata)
	if err != nil {
		return nil, eris.Wrap(err, "failed to unmarshal frontmatter")
	}

	return &metadata, nil
}

func (c *Import) writeImport(fileContent []string, dryRun bool) error {
	data := strings.Join(fileContent, "\n")

	metadata, err := readfileImport(data)
	if err != nil {
		return eris.Wrap(err, "failed to read file for import")
	}

	exist, err := c.state.Dal.FileExists(metadata.Path)
	if err != nil {
		return eris.Wrap(err, "failed to check if file xist")
	}

	if exist {
		return eris.Wrapf(errs.ErrDuplicate, "file exist %s", metadata.Path)
	}

	parent := pathutils.New(fmt.Sprintf("%s/%s", c.state.Path, metadata.Path)).GetParent()

	if !dryRun {
		err = os.MkdirAll(parent, constants.FolderPermission)
		if err != nil {
			return eris.Wrap(err, "failed to create parent dir")
		}
	}

	p := "false"
	if metadata.Private {
		p = "true"
	}

	d := []string{
		"---",
		fmt.Sprintf(`title: "%s"`, metadata.Title),
		fmt.Sprintf(`private: %s`, p),
		"---",
		metadata.Content,
	}

	if !dryRun {
		err = c.create.CreateFile(metadata.Path, metadata.Title, []byte(strings.Join(d, "\n")))
		if err != nil {
			return eris.Wrap(err, "failed to write file for import")
		}
	} else {
		if utilslib.FileExist(fmt.Sprintf("%s/%s", c.state.Path, metadata.Path)) {
			fmt.Printf("Filename %s exist\n", metadata.Path)
		}

		if !strings.HasSuffix(metadata.Path, ".md"){
			fmt.Println("Path doesn't end in .md\n")
		}

		if strings.HasSuffix(metadata.Path, "/.md"){
			fmt.Println("Path ends in /.md seems like you forgot a filename\n")
		}
	}

	return nil
}
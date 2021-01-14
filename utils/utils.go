// Package utils contains various methods I don't have a better place for
package utils

import (
	"fmt"
	"github.com/charmbracelet/glamour"
	"github.com/ericaro/frontmatter"
	"github.com/hjertnes/roam/errs"
	"github.com/hjertnes/roam/models"
	"github.com/hjertnes/utils"
	"github.com/rotisserie/eris"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

// GetPath returns the value of the ROAM environment variable or a default value if not set.
func GetPath() string {
	if utils.FileExist("./.roam"){
		data, err := ioutil.ReadFile("./.roam")
		if err == nil{
			return utils.ExpandTilde(strings.ReplaceAll(string(data), "\n", ""))
		}
	}
	path, isSet := os.LookupEnv("ROAM")

	if !isSet {
		return utils.ExpandTilde("~/txt/roam2")
	}

	return utils.ExpandTilde(path)
}

// GetEditor returns the value of the EDITOR enlivenment variable or a default value if not set.
func GetEditor() string {
	editor, isSet := os.LookupEnv("EDITOR")

	if !isSet {
		return "emacs"
	}

	return editor
}

// FilesToChoices maps a []models.File to []selectinput.Choice.
func FilesToChoices(input []models.File) []models.Choice {
	paths := make([]models.Choice, 0)

	for _, r := range input {
		paths = append(paths, models.Choice{Title: r.Path, Value: r.ID})
	}

	return paths
}

// EditNote opens the specified file in EDITOR.
func EditNote(path string) error {
	editor := GetEditor()
	cmd := exec.Command(editor, path) // #nosec G204

	err := cmd.Start()
	if err != nil {
		return eris.Wrap(err, "could not open file in editor")
	}

	return nil
}

// ViewNote renders the specified note as markdown in terminal.
func ViewNote(path string) error {
	data, err := ioutil.ReadFile(filepath.Clean(path))
	if err != nil {
		return eris.Wrap(err, "could not read file")
	}

	metadata := models.Frontmatter{}

	err = frontmatter.Unmarshal(data, &metadata)
	if err != nil {
		return eris.Wrap(err, "could not unmarkshal frontmatter")
	}

	r, _ := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
	)

	out, err := r.Render(fmt.Sprintf("# %s\n%s", metadata.Title, metadata.Content))
	if err != nil {
		return eris.Wrap(err, "failed to render markdown file")
	}

	fmt.Print(out)

	return nil
}

func PrintListOfLinks(output []string, links []models.File) []string{
	if len(links) == 0 {
		output = append(output, "No links")
	}

	for _, link := range links {
		output = append(output, fmt.Sprintf("- %s\n", link.Path))
	}

	return output
}

// GetFile returns a file with a id form a list of files.
func GetFile(files []models.File, id string) (*models.File, error) {
	var f *models.File

	for i := range files {
		if files[i].ID == id {
			f = &files[i]

			break
		}
	}

	if f == nil {
		return nil, eris.Wrap(errs.ErrNotFound, "no match")
	}

	return f, nil
}

// ConvertTemplateFiles convert TemplateFiles to Choice.
func ConvertTemplateFiles(templates []models.TemplateFile) []models.Choice {
	result := make([]models.Choice, 0)

	for _, f := range templates {
		result = append(result, models.Choice{
			Title: f.Title,
			Value: f.Filename,
		})
	}

	return result
}

func RemoveFilenameFromPath(path string) string{
	res := make([]string, 0)

	elems := strings.Split(path, "/")

	for _, e := range elems{
		if strings.HasSuffix(e, ".md"){
			continue
		}

		res = append(res, e)
	}

	return strings.Join(res, "/")
}

func GetParent(path string) string{
	res := make([]string, 0)

	elems := strings.Split(path, "/")

	for i, e := range elems{
		if i == len(elems) -1{
			continue
		}

		res = append(res, e)
	}

	return strings.Join(res, "/")
}

func BuildVectorSearch(input string) string {
	if !strings.Contains(input, " ") {
		return fmt.Sprintf("%s:*", input)
	}

	output := make([]string, 0)

	for _, l := range strings.Split(input, " ") {
		output = append(output, fmt.Sprintf("%s:*", l))
	}

	return strings.Join(output, "&")
}

func CleanLink(input string) string {
	return strings.ReplaceAll(strings.ReplaceAll(input, "[[", ""), "]]", "")
}

var NoteLinkRegexp = regexp.MustCompile(`\[\[([\d\w\s./]+)]]`)

func DestructPath(path string)(string, string){
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

func FixUrl(input string) string{
	output := strings.ReplaceAll(input, " ", "%20")
	output = strings.ReplaceAll(output, ".md", ".html")
	return output
}

func GetLast(path string) string{
	elems := strings.Split(path, "/")
	return elems[len(elems)-1]
}

func Readfile(path string) (*models.Frontmatter, error) {
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

func ReadfileImport(data string) (*models.ImportFrontmatter, error) {
	metadata := models.ImportFrontmatter{}

	err := frontmatter.Unmarshal([]byte(data), &metadata)
	if err != nil {
		return nil, eris.Wrap(err, "failed to unmarshal frontmatter")
	}

	return &metadata, nil
}

func ErrorHandler(err error) {
	if err != nil {
		if eris.Is(err, errs.ErrNotFound) {
			fmt.Println("No matches to search query")
		}
		fmt.Println("Error")
		fmt.Println(eris.ToString(err, true))

		os.Exit(0)
	}
}
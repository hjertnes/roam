package commands

import (
	"context"
	"fmt"
	"github.com/hjertnes/roam/commands/findSearch"
	"github.com/hjertnes/roam/commands/findSelect"
	"github.com/hjertnes/roam/commands/noteTitle"
	"github.com/hjertnes/roam/commands/selectType"
	"github.com/hjertnes/roam/configuration"
	dal2 "github.com/hjertnes/roam/dal"
	"github.com/hjertnes/roam/utils"
	utilslib "github.com/hjertnes/utils"
	"github.com/jackc/pgx/v4/pgxpool"
	"io/ioutil"
	"os/exec"
	"strings"
	"time"
)

func FindEdit(path string){
	search, err := findSearch.Run()
	utils.ErrorHandler(err)

	fmt.Println("Loading...")

	conf, err := configuration.ReadConfigurationFile(fmt.Sprintf("%s/.config/config.yaml", path))
	utils.ErrorHandler(err)

	ctx := context.Background()
	pxp, err := pgxpool.Connect(ctx, conf.DatabaseConnectionString)

	dal := dal2.New(ctx, pxp)

	result, err := dal.Find(search)

	utils.ErrorHandler(err)

	paths := make([]string, 0)

	for _, r := range result {
		paths = append(paths, r.Path)
	}

	findSelect.Run(paths, true, dal)

}

func Create(path, filepath string) {
	conf, err := configuration.ReadConfigurationFile(fmt.Sprintf("%s/.config/config.yaml", path))
	utils.ErrorHandler(err)

	ctx := context.Background()
	pxp, err := pgxpool.Connect(ctx, conf.DatabaseConnectionString)

	dal := dal2.New(ctx, pxp)

	title, err := noteTitle.Run()
	utils.ErrorHandler(err)

	template := selectType.Run(conf.Templates)

	if utilslib.FileExist(filepath) {
		fmt.Println("File already exist")
		return
	}

	templatedata, err := ioutil.ReadFile(fmt.Sprintf("%s/.config/templates/%s", path, template.Filename))
	utils.ErrorHandler(err)

	now := time.Now()

	noteText := strings.ReplaceAll(string(templatedata), "$$TITLE$$", title)
	noteText = strings.ReplaceAll(noteText, "$$DATE$$", now.Format(conf.DateFormat))
	noteText = strings.ReplaceAll(noteText, "$$TIME$$", now.Format(conf.TimeFormat))
	noteText = strings.ReplaceAll(noteText, "$$DATETIME$$", now.Format(conf.DateTimeFormat))

	err = ioutil.WriteFile(filepath, []byte(noteText), 0700)
	utils.ErrorHandler(err)

	err = dal.Create(filepath, title, "")
	utils.ErrorHandler(err)

	err = dal.SetOpened(filepath)
	utils.ErrorHandler(err)

	editor := utils.GetEditor()

	cmd := exec.Command(editor, filepath)

	err = cmd.Start()
	utils.ErrorHandler(err)
}

func Daily(path, date string) {
	conf, err := configuration.ReadConfigurationFile(fmt.Sprintf("%s/.config/config.yaml", path))
	utils.ErrorHandler(err)

	ctx := context.Background()
	pxp, err := pgxpool.Connect(ctx, conf.DatabaseConnectionString)
	dal := dal2.New(ctx, pxp)
	title := ""

	filename := fmt.Sprintf("%s/Daily Notes/%s.md", path, date)

	if !utilslib.FileExist(filename) {

		templatedata, err := ioutil.ReadFile(fmt.Sprintf("%s/.config/templates/%s", path, "daily.txt"))
		utils.ErrorHandler(err)

		now := time.Now()

		noteText := strings.ReplaceAll(string(templatedata), "$$TITLE$$", title)
		noteText = strings.ReplaceAll(noteText, "$$DATE$$", now.Format(conf.DateFormat))
		noteText = strings.ReplaceAll(noteText, "$$TIME$$", now.Format(conf.TimeFormat))
		noteText = strings.ReplaceAll(noteText, "$$DATETIME$$", now.Format(conf.DateTimeFormat))

		err = ioutil.WriteFile(filename, []byte(noteText), 0700)
		utils.ErrorHandler(err)

		err = dal.Create(filename, title, "")
		utils.ErrorHandler(err)
	}

	err = dal.SetOpened(filename)
	utils.ErrorHandler(err)

	editor := utils.GetEditor()

	cmd := exec.Command(editor, filename)

	err = cmd.Start()
	utils.ErrorHandler(err)
}


func FindView(path string){
	search, err := findSearch.Run()
	utils.ErrorHandler(err)

	fmt.Println("Loading...")

	conf, err := configuration.ReadConfigurationFile(fmt.Sprintf("%s/.config/config.yaml", path))
	utils.ErrorHandler(err)

	ctx := context.Background()
	pxp, err := pgxpool.Connect(ctx, conf.DatabaseConnectionString)

	dal := dal2.New(ctx, pxp)

	result, err := dal.Find(search)

	utils.ErrorHandler(err)

	paths := make([]string, 0)

	for _, r := range result {
		paths = append(paths, r.Path)
	}

	findSelect.Run(paths, false, dal)

}
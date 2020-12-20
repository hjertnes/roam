package main

import (
	"context"
	"fmt"
	"github.com/ericaro/frontmatter"
	"github.com/hjertnes/roam/configuration"
	dal2 "github.com/hjertnes/roam/dal"
	"github.com/hjertnes/roam/migration"
	utils "github.com/hjertnes/utils"
	"github.com/jackc/pgx/v4/pgxpool"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)



func getPath() string{
	path, isSet := os.LookupEnv("ROAM")

	if !isSet{
		return utils.ExpandTilde("~/txt/roam2")
	}

	return utils.ExpandTilde(path)
}

func errorHandler(err error){
	if err != nil{
		panic(err)
	}
}

func setup(path string){
	configFolder := fmt.Sprintf("%s/.config", path)
	configFile := fmt.Sprintf("%s/config.yaml", configFolder)
	if !utils.FileExist(configFolder){
		err := os.Mkdir(configFolder, 0755)
		errorHandler(err)
	}

	if !utils.FileExist(configFile){
		err := configuration.CreateConfigurationFile(configFile)
		errorHandler(err)
	}
}

func getEditor() string{
	editor, isSet := os.LookupEnv("EDITOR")

	if !isSet{
		return "emacs"
	}

	return editor
}

func openConfig(path string) {
	editor := getEditor()
	configFile := fmt.Sprintf("%s/.config/config.yaml", path)
	cmd := exec.Command(editor, configFile)

	err := cmd.Start()
	errorHandler(err)
}

func migrate(path string){
	conf, err := configuration.ReadConfigurationFile(fmt.Sprintf("%s/.config/config.yaml", path))
	errorHandler(err)

	ctx := context.Background()
	pxp, err := pgxpool.Connect(ctx, conf.DatabaseConnectionString)
	errorHandler(err)

	mig := migration.New(ctx, pxp)
	err = mig.Migrate()
	errorHandler(err)

}

type fm struct {
	Title string `fm:"title"`
	Private bool `fm:"private"`
	Content string `fm:"content"`
}

func sync(path string){
	conf, err := configuration.ReadConfigurationFile(fmt.Sprintf("%s/.config/config.yaml", path))
	errorHandler(err)

	ctx := context.Background()
	pxp, err := pgxpool.Connect(ctx, conf.DatabaseConnectionString)

	dal := dal2.New(ctx, pxp)

	// TODO do partial updates based on LastModified, opened_at updated_at
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

		data, err := ioutil.ReadFile(path)
		metadata := fm{}
		err = frontmatter.Unmarshal(data, &metadata)
		if err != nil{
			fmt.Println("Here?")
			return err
		}
		exist, err := dal.Exists(path)
		if err != nil{
			return err
		}

		if exist{
			err = dal.Update(path, metadata.Title, metadata.Content)
			if err != nil{
				return err
			}
		} else {
			err = dal.Create(path, metadata.Title, metadata.Content)
			if err != nil{
				return err
			}
		}
		//fmt.Println(path)
		//fmt.Println(info.Name())
		return nil
	})

	if err != nil{
		panic(err)
	}
}

//todo plan repo
//

func edit(path string){
	if len(os.Args) == 2 {
		help()
	}

	switch os.Args[2] {
	case "config":
		openConfig(path)
		return
	default:
		help()
		return
	}
}

func help(){
	fmt.Printf("roam\n")
	fmt.Printf("A command line utility that will replace my use of org-roam\n")
	fmt.Printf("\n")
	fmt.Printf("Usage:\n")
	fmt.Printf("\troam\t[command]\n")
	fmt.Printf("\n")
	fmt.Printf("Available commands:\n")
	fmt.Printf("\thelp\tprints this text\n")
	fmt.Printf("\tinit\tcreates configuration files\n")
	fmt.Printf("\tmigrate\tsets up the database\n")
	fmt.Printf("\tedit")
	fmt.Printf("\t\tconfig\topens config file in EDITOR")
}
//"

func main(){
	path := getPath()
	if len(os.Args) == 1{
		help()
		return
	}

	switch os.Args[1] {
	case "init":
		setup(path)
		return
	case "edit":
		edit(path)
		return
	case "migrate":
		migrate(path)
		return
	case "sync":
		sync(path)
		return
	default:
		help()
		return
	}
}
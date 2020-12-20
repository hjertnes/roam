package utils

import (
	"github.com/hjertnes/utils"
	"os"
)

func GetPath() string{
	path, isSet := os.LookupEnv("ROAM")

	if !isSet{
		return utils.ExpandTilde("~/txt/roam2")
	}

	return utils.ExpandTilde(path)
}

func ErrorHandler(err error){
	if err != nil{
		panic(err)
	}
}

func GetEditor() string{
	editor, isSet := os.LookupEnv("EDITOR")

	if !isSet{
		return "emacs"
	}

	return editor
}



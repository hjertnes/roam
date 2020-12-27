package utils

import (
	"github.com/hjertnes/utils"
	"os"

)

// GetPath returns the value of the ROAM enviornment variable or a default value if not set
func GetPath() string {
	path, isSet := os.LookupEnv("ROAM")

	if !isSet {
		return utils.ExpandTilde("~/txt/roam2")
	}

	return utils.ExpandTilde(path)
}

// GetEditor returns the value of the EDITOR enviornment variable or a default value if not set
func GetEditor() string {
	editor, isSet := os.LookupEnv("EDITOR")

	if !isSet {
		return "emacs"
	}

	return editor
}

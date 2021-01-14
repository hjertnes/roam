package version

import (
	"fmt"
	"github.com/hjertnes/roam/constants"
)

func Run(){
	fmt.Printf("Version: %s\n", constants.Version)
}

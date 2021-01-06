package textinput2

import (
	"fmt"
	"github.com/rotisserie/eris"
)

func Run(label string) (string, error){
	var input string
	fmt.Printf("%s: ", label)
	_, err := fmt.Scanln(&input)
	if err != nil{
		return input, eris.Wrap(err, "could not get text input from user")
	}


	return input, nil


}

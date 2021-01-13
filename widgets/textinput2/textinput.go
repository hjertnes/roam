package textinput2

import (
	"fmt"
	"github.com/rotisserie/eris"
)

func Run(label string) (string, error){
	var input string
	fmt.Printf("%s: ", label)
	_, err := fmt.Scan(&input)
	if err != nil{
		fmt.Println(input)
		return input, eris.Wrap(err, "could not get text input from user")
	}


	return input, nil


}

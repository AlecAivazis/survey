package main

import (
	"fmt"

	"github.com/alecaivazis/survey"
)

func main() {
	prompt := &survey.Confirm{
		Message: "Are you happy?",
	}

	answer, err := survey.AskOne(prompt)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Printf("response string: %s\n", answer)
}

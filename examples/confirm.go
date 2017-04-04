package main

import (
	"fmt"

	"github.com/alecaivazis/survey"
)

func main() {
	var happy bool
	prompt := &survey.Confirm{
		Message: "Are you happy?",
		Default: true,
		Answer:  &happy,
	}

	answer, err := survey.AskOne(prompt)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Printf("response string: %s\n", answer)
	fmt.Printf("response happy: %t\n", happy)
}

package main

import (
	"fmt"

	"github.com/alecaivazis/survey"
)

func main() {
	prompt := &survey.Confirm{
		Message: "Are you happy?",
	}
	ans := false
	err := survey.AskOne(prompt, &ans, nil)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Printf("response: %v\n", ans)
}

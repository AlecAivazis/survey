package main

import (
	"fmt"

	"github.com/alecaivazis/survey"
)

// the questions to ask
var simpleQs = []*survey.Question{
	{
		Name: "color",
		Prompt: &survey.Select{
			Message: "Choose a color:",
			Options: []string{"red", "blue", "green"},
		},
		Validate: survey.Required,
	},
}

func main() {

	_, err := survey.Ask(simpleQs)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

}

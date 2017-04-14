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
			Options: []string{"red", "blue", "green", "yellow"},
			Default: "yellow",
		},
		Validate: survey.Required,
	},
}

func main() {

	answers, err := survey.Ask(simpleQs)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Printf("%s chose %s.\n", answers["name"], answers["color"])
}

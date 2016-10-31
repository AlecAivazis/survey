package main

import (
	"fmt"
	"github.com/alecaivazis/survey"
)

// the questions to ask
var simpleQs = []*survey.Question{
	{
		Name: "name",
		Prompt: &survey.Input{
			Message: "What is your name?",
			Default: "Johnny Appleseed",
		},
	},
	{
		Name: "color",
		Prompt: &survey.Choice{
			Message: "Choose a color:",
			Choices: []string{"red", "blue", "green"},
			Default: "red",
		},
	},
}

func main() {

	answers, err := survey.Ask(simpleQs)

	if err != nil {
		fmt.Println("\n", err.Error())
		return
	}

	fmt.Printf("%s chose %s.\n", answers["name"], answers["color"])
}

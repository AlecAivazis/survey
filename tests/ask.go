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
			Choices: []string{"red", "blue", "green", "yellow"},
			Default: "yellow",
		},
		Validate: survey.Required,
	},
}

func main() {

	fmt.Println("Asking many.")

	answers, err := survey.Ask(simpleQs)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("%s chose %s.\n", answers["name"], answers["color"])

	fmt.Println("Asking one.")
	answer, err := survey.AskOne(simpleQs[0])
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("Answered with %v.\n", answer)

	fmt.Println("Asking one with validation.")
	answer, err := survey.AskOneValidate(simpleQs[0], survey.Required)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("Answered with %v.\n", answer)
}

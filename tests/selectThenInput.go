package main

import (
	"fmt"

	"gopkg.in/AlecAivazis/survey.v1"
)

// the questions to ask
var simpleQs = []*survey.Question{
	{
		Name: "color",
		Prompt: survey.NewSingleSelect().SetMessage("Choose a color:").
			AddOption("red", nil, false).
			AddOption("blue", nil, false).
			AddOption("green", nil, false),
		Validate: survey.Required,
	},
	{
		Name: "name",
		Prompt: &survey.Input{
			Message: "What is your name?",
		},
		Validate: survey.Required,
	},
}

func main() {
	answers := struct {
		Color *survey.Option
		Name  string
	}{}
	// ask the question
	err := survey.Ask(simpleQs, &answers)

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	// print the answers
	fmt.Printf("%s chose %s.\n", answers.Name, answers.Color)
}

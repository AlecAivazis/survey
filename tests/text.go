package main

import (
	"fmt"

	"gopkg.in/AlecAivazis/survey.v1"
)

var simpleQs = []*survey.Question{
	{
		Name: "color",
		Prompt: &survey.Select{
			Message: "select1:",
			Options: []string{"red", "blue", "green"},
		},
		Validate: survey.Required,
	},
	{
		Name: "name",
		Prompt: &survey.Input{
			Message: "Input your name:",
		},
		Validate: survey.Required,
	},
	{
		Prompt: &survey.Text{
			Message: "Thanks for you input.",
		},
		Validate: survey.Required,
	},
	{
		Name: "color2",
		Prompt: &survey.Select{
			Message: "select2:",
			Options: []string{"red", "blue", "green"},
		},
		Validate: survey.Required,
	},
}

func main() {
	answers := struct {
		Color  string
		Color2 string
		Name   string
	}{}
	// ask the question
	err := survey.Ask(simpleQs, &answers)

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	// print the answers
	fmt.Printf("%s and %s, name is %s\n", answers.Color, answers.Color2, answers.Name)
}

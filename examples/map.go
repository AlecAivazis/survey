package main

import (
	"fmt"

	"gopkg.in/AlecAivazis/survey.v1"
)

// the questions to ask
var simpleQs = []*survey.Question{
	{
		Name: "name",
		Prompt: &survey.Input{
			Message: "What is your name?",
		},
		Validate: survey.Required,
	},
	{
		Name: "color",
		Prompt: survey.NewSingleSelect().SetMessage("select1:").
			AddOption("red", nil, false).
			AddOption("blue", nil, false).
			AddOption("green", nil, false),
		Validate: survey.Required,
	},
}

func main() {
	ansmap := make(map[string]interface{})

	// ask the question
	err := survey.Ask(simpleQs, &ansmap)

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	// print the answers
	fmt.Printf("%s chose %s.\n", ansmap["name"], ansmap["color"])
}

package main

import (
	"fmt"

	"gopkg.in/AlecAivazis/survey.v1"
)

var simpleDoubleQs = []*survey.Question{
	{
		Name: "color",
		Prompt: survey.NewSingleSelect().SetMessage("select1:").
			AddOption("red", nil, false).
			AddOption("blue", nil, false).
			AddOption("green", nil, false),
		Validate: survey.Required,
	},
	{
		Name: "color2",
		Prompt: survey.NewSingleSelect().SetMessage("select2:").
			AddOption("red", nil, false).
			AddOption("blue", nil, false).
			AddOption("green", nil, false),
		Validate: survey.Required,
	},
}

func main() {
	answers := struct {
		Color  *survey.Option
		Color2 *survey.Option
	}{}
	// ask the question
	err := survey.Ask(simpleDoubleQs, &answers)

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	// print the answers
	fmt.Printf("%s and %s.\n", answers.Color.Display, answers.Color2.Display)
}

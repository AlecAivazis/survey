package main

import (
	"fmt"

	"gopkg.in/AlecAivazis/survey.v1"
)

type colorDetail struct {
	Name string
	Hex string
	WikiLink string
}

// the questions to ask
var simpleQs = []*survey.Question{
	{
		Name: "name",
		Prompt: &survey.Input{
			Message: "What is your name?",
		},
		Validate:  survey.Required,
		Transform: survey.Title,
	},
	{
		Name: "color",
		Prompt: survey.NewSingleSelect().SetMessage("Select Color:").
			AddOption("red", &colorDetail{"Red", "FF0000", "https://en.wikipedia.org/wiki/Red"}, false).
			AddOption("blue", &colorDetail{"Blue", "0000FF", "https://en.wikipedia.org/wiki/Blue"}, false).
			AddOption("green", &colorDetail{"Green", "00FF00", "https://en.wikipedia.org/wiki/Green"}, true),
		Validate: survey.Required,
	},
}

func main() {
	answers := struct {
		Name  string
		Color *survey.Option
	}{}

	// ask the question
	err := survey.Ask(simpleQs, &answers)

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	// print the answers
	color := answers.Color.Value.(*colorDetail)
	fmt.Printf("%s chose %s which has a hex value of #%s and you can read about it @%s.\n", answers.Name, color.Name, color.Hex, color.WikiLink)
}

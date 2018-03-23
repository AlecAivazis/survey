package main

import (
	"fmt"

	"gopkg.in/AlecAivazis/survey.v1"
)

// the questions to ask
var simpleQs = []*survey.Question{
	{
		Name: "letter",
		Prompt: survey.NewSingleSelect().
			SetMessage("Choose a letter:").
			AddOption("a", nil, false).
			AddOption("b", nil, false).
			AddOption("c", nil, false).
			AddOption("d", nil, false).
			AddOption("e", nil, false).
			AddOption("f", nil, false).
			AddOption("g", nil, false).
			AddOption("h", nil, false).
			AddOption("i", nil, false).
			AddOption("j", nil, false),
		Validate: survey.Required,
	},
}

func main() {
	answers := struct {
		Letter *survey.Option
	}{}

	// ask the question
	err := survey.Ask(simpleQs, &answers)

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	// print the answers
	fmt.Printf("you chose %s.\n", answers.Letter)
}

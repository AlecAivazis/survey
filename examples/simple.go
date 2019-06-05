package main

import (
	"errors"
	"fmt"

	"github.com/AlecAivazis/survey/v2"
)

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
		Prompt: &survey.Select{
			Message: "Choose a color:",
			Options: []string{"red", "blue", "green"},
		},
		Validate: survey.Required,
	},
	{
		Name: "friends",
		Prompt: &survey.MultiInput{
			Message: "Enter the names of your friends:",
		},
		Validate: func(val interface{}) error {
			if list, ok := val.([]string); !ok || len(list) < 1 {
				return errors.New("there should be at least one response.")
			}
			return nil
		},
	},
}

func main() {
	answers := struct {
		Name    string
		Color   string
		Friends []string
	}{}

	// ask the question
	err := survey.Ask(simpleQs, &answers)

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	// print the answers
	fmt.Printf("%s chose %s and is friends with %s.\n", answers.Name, answers.Color, answers.Friends)
}

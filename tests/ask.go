//go:build ignore
// +build ignore

package main

import (
	"fmt"
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
		Prompt: &survey.Select{
			Message: "Choose a color:",
			Options: []string{"red", "blue", "green", "yellow"},
			Default: "yellow",
		},
		Validate: survey.Required,
	},
}

var singlePrompt = &survey.Input{
	Message: "What is your name?",
	Default: "Johnny Appleseed",
}

func main() {

	fmt.Println("Asking many.")
	// a place to store the answers
	ans := struct {
		Name  string
		Color string
	}{}
	err := survey.Ask(simpleQs, &ans)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("Asking one.")
	answer := ""

	err = survey.AskOne(singlePrompt, &answer)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("Answered with %v.\n", answer)

	fmt.Println("Asking one with validation.")
	vAns := ""
	err = survey.AskOne(&survey.Input{Message: "What is your name?"}, &vAns, survey.WithValidator(survey.Required))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("Answered with %v.\n", vAns)
}

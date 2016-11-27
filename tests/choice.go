package main

import (
	"fmt"

	"github.com/alecaivazis/survey"
)

type choiceTestTable struct {
	Name   string
	Prompt *survey.Choice
}

var table = []*choiceTestTable{
	{
		"standard", &survey.Choice{
			Message: "Choose a color:",
			Choices: []string{"red", "blue", "green"},
			Default: "red",
		},
	},
	{
		"short", &survey.Choice{
			Message: "Choose a color:",
			Choices: []string{"red", "blue"},
			Default: "red",
		},
	},
	{
		"default", &survey.Choice{
			Message: "Choose a color (should default blue):",
			Choices: []string{"red", "blue", "green"},
			Default: "blue",
		},
	},
}

func main() {
	// go over every entry in the table
	for _, entry := range table {
		// tell the user what we are going to ask them
		fmt.Println(entry.Name)
		// perform the ask
		_, err := survey.AskOne(entry.Prompt)
		if err != nil {
			fmt.Printf("AskOne on %v's prompt failed: %v.", entry.Name, err.Error())
			break
		}
	}
}

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
		},
		Validate: survey.Required,
	},
}

func main() {
	ansmap := make(map[string]interface{})

	// ask the question
	err := survey.Ask(simpleQs, &ansmap, survey.WithShowCursor(true))

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	// print the answers
	fmt.Printf("Your name is %s.\n", ansmap["name"])
}

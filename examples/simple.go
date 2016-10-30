package main

import (
	"fmt"
	"github.com/alecaivazis/probe"
)

// the questions to ask
var simpleQs = []*probe.Question{
	{
		Name:   "name",
		Prompt: &probe.Input{"What is your name?"},
	},
	{
		Name: "color",
		Prompt: &probe.Choice{
			Message: "Choose a color:",
			Choices: []string{"red", "blue", "green"},
		},
	},
}

func main() {

	answers, err := probe.Ask(simpleQs)

	if err != nil {
		fmt.Println("\n", err.Error())
		return
	}

	fmt.Printf("%s chose %s.\n", answers["name"], answers["color"])
}

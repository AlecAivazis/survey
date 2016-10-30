package main

import (
	"fmt"
	"github.com/alecaivazis/probe"
)

var qs = []*probe.Question{
	{
		Name:   "name",
		Prompt: &probe.Input{"What is your name?"},
	},
	{
		Name:   "birthday",
		Prompt: &probe.Input{"When is your birthday?"},
	},
	{
		Name: "gender",
		Prompt: &probe.Choice{
			Question: "Are you male or female?",
			Choices:  []string{"male", "female"},
		},
	},
}

func main() {
	answers := probe.Ask(qs)
	fmt.Println(answers)
}

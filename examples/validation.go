package main

import (
	"fmt"
	"github.com/alecaivazis/probe"
)

// the questions to ask
var validationQs = []*probe.Question{
	{
		Name:     "name",
		Prompt:   &probe.Input{"What is your name?"},
		Validate: probe.NonNull,
	},
}

func main() {

	answers, err := probe.Ask(validationQs)

	if err != nil {
		fmt.Println("\n", err.Error())
		return
	}

	fmt.Printf("%s chose %s.\n", answers["name"], answers["color"])
}

package main

import (
	"errors"
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
	{
		Name:   "valid",
		Prompt: &probe.Input{"Enter 'foo':"},
		Validate: func(str string) error {
			// if the input matches the expectation
			if str != "foo" {
				return errors.New(fmt.Sprintf("You entered %s, not 'foo'.", str))
			}
			// nothing was wrong
			return nil
		},
	},
}

func main() {

	answers, err := probe.Ask(validationQs)

	if err != nil {
		fmt.Println("\n", err.Error())
	}
}

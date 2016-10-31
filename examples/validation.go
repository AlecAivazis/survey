package main

import (
	"errors"
	"fmt"
	"github.com/alecaivazis/survey"
)

// the questions to ask
var validationQs = []*survey.Question{
	{
		Name:     "name",
		Prompt:   &survey.Input{"What is your name?"},
		Validate: survey.Required,
	},
	{
		Name:   "valid",
		Prompt: &survey.Input{"Enter 'foo':"},
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

	answers, err := survey.Ask(validationQs)

	if err != nil {
		fmt.Println("\n", err.Error())
	}
}

package main

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/alecaivazis/survey"
)

// the questions to ask
var validationQs = []*survey.Question{
	{
		Name:     "name",
		Prompt:   &survey.Input{"What is your name?", ""},
		Validate: survey.Required,
	},
	{
		Name:   "valid",
		Prompt: &survey.Input{"Enter 'foo':", "not foo"},
		Validate: func(val interface{}) error {
			// if the value passed in is the zero value of the appropriate type
			if val == reflect.Zero(reflect.TypeOf(val)).Interface() {
				return errors.New("Value is required")
			}
			return nil
		},
	},
}

func main() {
	// the place to hold the answers
	answers := struct {
		Name  string
		Valid string
	}{}
	err := survey.Ask(validationQs, &answers)

	if err != nil {
		fmt.Println("\n", err.Error())
	}
}

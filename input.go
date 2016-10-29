package main

import (
	"fmt"
)

type Input struct {
	question string
}

func (input *Input) AskQuestion() {
	fmt.Print(fmt.Sprintf("%v ", input.question))
}

func (input *Input) Prompt() (string, error) {
	// a string to hold the user's input
	var res string
	// wait for a newline or carriage return
	fmt.Scanln(&res)
	// return the value
	return res, nil
}

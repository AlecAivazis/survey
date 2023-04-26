//go:build ignore
// +build ignore

package main

import (
	"fmt"
)

// the questions to ask
var defaultPasswordCharacterPrompt = &survey.Password{
	Message: "What is your password? (Default hide character)",
}
var customPasswordCharacterPrompt = &survey.Password{
	Message: "What is your password? (Custom hide character)",
}

func main() {

	var defaultPass string
	var customPass string

	// ask the question
	err := survey.AskOne(defaultPasswordCharacterPrompt, &defaultPass)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println()
	err = survey.AskOne(customPasswordCharacterPrompt, &customPass, survey.WithHideCharacter('-'))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	// print the answers
	fmt.Printf("Password 1: %s.\n Password 2: %s\n", defaultPass, customPass)
}

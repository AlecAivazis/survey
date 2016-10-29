package main

import (
	"fmt"
)

var qs []*Question = []*Question{
	{Name: "name", Prompt: &Input{"What is your name?"}},
	{Name: "birthday", Prompt: &Input{"When is your birthday?"}},
}

func main() {
	answers := Ask(qs)
	fmt.Println(answers)
}

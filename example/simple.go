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
		Prompt: &probe.Input{"Wh en is your birthday?"},
	},
}

func main() {
	answers := probe.Ask(qs)
	fmt.Println(answers)
}

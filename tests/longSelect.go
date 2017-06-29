package main

import "github.com/AlecAivazis/survey"

func main() {
	color := ""
	prompt := &survey.Select{
		Message: "Choose a color:",
		Options: []string{
			"red",
			"blue",
			"red",
			"blue",
			"red",
			"blue",
			"red",
			"blue",
			"red",
			"blue",
			"red",
			"blue",
			"red",
			"blue",
			"red",
			"blue",
			"red",
			"blue",
			"red",
			"blue",
			"red",
			"blue",
			"red",
			"blue",
			"red",
			"blue",
			"red",
			"blue",
			"red",
			"blue",
			"red",
			"blue",
			"green",
		},
	}
	survey.AskOne(prompt, &color, nil)
}

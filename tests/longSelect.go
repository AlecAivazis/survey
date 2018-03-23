package main

import "gopkg.in/AlecAivazis/survey.v1"

func main() {
	color := ""
	prompt := survey.NewSingleSelect().
		SetMessage("Choose a color:").
		AddOption("a", nil, false).
		AddOption("b", nil, false).
		AddOption("c", nil, false).
		AddOption("d", nil, false).
		AddOption("e", nil, false).
		AddOption("f", nil, false).
		AddOption("g", nil, false).
		AddOption("h", nil, false).
		AddOption("i", nil, false).
		AddOption("j", nil, false)
	survey.AskOne(prompt, &color, nil)
}

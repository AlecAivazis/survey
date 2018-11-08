package main

import "gopkg.in/AlecAivazis/survey.v1"

func main() {
	answer := ""
	survey.AskOne(&survey.Text{Message: "This is a info message!"}, nil, nil)
	survey.AskOne(&survey.Text{Message: "This is a info message!", Level: survey.Info}, &answer, nil)
	survey.AskOne(&survey.Text{Message: "This is a warning message!", Level: survey.Warning}, &answer, nil)
	survey.AskOne(&survey.Text{Message: "This is a danger message!", Level: survey.Danger}, &answer, nil)
}

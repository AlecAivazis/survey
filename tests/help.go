package main

import (
	"gopkg.in/AlecAivazis/survey.v1"
	"gopkg.in/AlecAivazis/survey.v1/tests/util"
)

var (
	confirmAns     = false
	inputAns       = ""
	multiselectAns = make(survey.Options, 0)
	selectAns      = &survey.Option{}
	passwordAns    = ""
	goodTable = []TestUtil.TestTableEntry{
		{
			"confirm", &survey.Confirm{
				Message: "Is it raining?",
				Help:    "Go outside, if your head becomes wet the answer is probably 'yes'",
			}, &confirmAns,
		},
		{
			"input", &survey.Input{
				Message: "What is your phone number:",
				Help:    "Phone number should include the area code, parentheses optional",
			}, &inputAns,
		},
		{
			"multi-select",
			survey.NewMultiSelect().SetMessage("What days are you available:").
				SetHelp("We are closed weekends and avaibility is limited on Wednesday").
				AddOption("Monday", nil, true).
				AddOption("Tuesday", nil, true).
				AddOption("Wednesday", nil, false).
				AddOption("Thursday", nil, true).
				AddOption("Friday", nil, true),
			&multiselectAns,
		},
		{
			"select",
			survey.NewSingleSelect().SetMessage("Choose a color:").
				SetHelp("Blue is the best color, but it is your choice").
				AddOption("red", nil, false).
				AddOption("blue", nil, true).
				AddOption("green", nil, false),
			&selectAns,
		},
		{
			"password", &survey.Password{
				Message: "Enter a secret:",
				Help:    "Don't really enter a secret, this is just for testing",
			}, &passwordAns,
		},
	}
)


func main() {
	TestUtil.RunTable(goodTable)
}

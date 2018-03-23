package main

import (
	"gopkg.in/AlecAivazis/survey.v1"
	"gopkg.in/AlecAivazis/survey.v1/tests/util"
)

var table = []TestUtil.TestTableEntry{
	{
		"standard",
		survey.NewMultiSelect().SetMessage("What days do you prefer:").
			AddOption("Sunday", nil, false).
			AddOption("Monday", nil, false).
			AddOption("Tuesday", nil, false).
			AddOption("Wednesday", nil, false).
			AddOption("Thursday", nil, false).
			AddOption("Friday", nil, false).
			AddOption("Saturday", nil, false),
		&survey.Options{},
	},
	{
		"default (sunday, tuesday)",
		survey.NewMultiSelect().SetMessage("What days do you prefer:").
			AddOption("Sunday", nil, true).
			AddOption("Monday", nil, false).
			AddOption("Tuesday", nil, true).
			AddOption("Wednesday", nil, false).
			AddOption("Thursday", nil, false).
			AddOption("Friday", nil, false).
			AddOption("Saturday", nil, false),
		&survey.Options{},
	},
	{
		"no help - type ?",
		survey.NewMultiSelect().SetMessage("What days do you prefer:").
			AddOption("Sunday", nil, false).
			AddOption("Monday", nil, false).
			AddOption("Tuesday", nil, false).
			AddOption("Wednesday", nil, false).
			AddOption("Thursday", nil, false).
			AddOption("Friday", nil, false).
			AddOption("Saturday", nil, false),
		&survey.Options{},
	},
	{
		"can navigate with j/k",
		survey.NewMultiSelect().SetMessage("What days do you prefer:").
			AddOption("Sunday", nil, false).
			AddOption("Monday", nil, false).
			AddOption("Tuesday", nil, false).
			AddOption("Wednesday", nil, false).
			AddOption("Thursday", nil, false).
			AddOption("Friday", nil, false).
			AddOption("Saturday", nil, false).
			SetVimMode(true),
		&survey.Options{},
	},
}

func main() {
	TestUtil.RunTable(table)
}

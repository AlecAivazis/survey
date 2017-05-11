package main

import (
	"github.com/AlecAivazis/survey"
	"github.com/AlecAivazis/survey/tests/util"
)

var answer = false

var goodTable = []TestUtil.TestTableEntry{
	{
		"Enter 'yes'", &survey.Confirm{
			Message: "yes:",
		}, &answer,
	},
	{
		"Enter 'no'", &survey.Confirm{
			Message: "yes:",
		}, &answer,
	},
	{
		"default", &survey.Confirm{
			Message: "yes:",
			Default: true,
		}, &answer,
	},
	{
		"not recognized (enter random letter)", &survey.Confirm{
			Message: "yes:",
			Default: true,
		}, &answer,
	},
	{
		"no help - type '?'", &survey.Confirm{
			Message: "yes:",
			Default: true,
		}, &answer,
	},
}

func main() {
	TestUtil.RunTable(goodTable)
}

//go:build ignore

package main

import (
	"github.com/AlecAivazis/survey/v2"
	TestUtil "github.com/AlecAivazis/survey/v2/tests/util"
)

var answer = false

var goodTable = []TestUtil.TestTableEntry{
	{
		"Enter 'yes'", &survey.Confirm{
			Message: "yes:",
		}, &answer, nil,
	},
	{
		"Enter 'no'", &survey.Confirm{
			Message: "yes:",
		}, &answer, nil,
	},
	{
		"default", &survey.Confirm{
			Message: "yes:",
			Default: true,
		}, &answer, nil,
	},
	{
		"not recognized (enter random letter)", &survey.Confirm{
			Message: "yes:",
			Default: true,
		}, &answer, nil,
	},
	{
		"no help - type '?'", &survey.Confirm{
			Message: "yes:",
			Default: true,
		}, &answer, nil,
	},
}

func main() {
	TestUtil.RunTable(goodTable)
}

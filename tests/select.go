package main

import (
	"github.com/AlecAivazis/survey"
	"github.com/AlecAivazis/survey/tests/util"
)

var answer = ""

var goodTable = []TestUtil.TestTableEntry{
	{
		"standard", &survey.Select{
			Message: "Choose a color:",
			Options: []string{"red", "blue", "green"},
		}, &answer,
	},
	{
		"short", &survey.Select{
			Message: "Choose a color:",
			Options: []string{"red", "blue"},
		}, &answer,
	},
	{
		"default", &survey.Select{
			Message: "Choose a color (should default blue):",
			Options: []string{"red", "blue", "green"},
			Default: "blue",
		}, &answer,
	},
	{
		"one", &survey.Select{
			Message: "Choose one:",
			Options: []string{"hello"},
		}, &answer,
	},
	{
		"no help, type ?", &survey.Select{
			Message: "Choose a color:",
			Options: []string{"red", "blue"},
		}, &answer,
	},
}

var badTable = []TestUtil.TestTableEntry{
	{
		"no options", &survey.Select{
			Message: "Choose one:",
		}, &answer,
	},
}

func main() {
	TestUtil.RunTable(goodTable)
	TestUtil.RunErrorTable(badTable)
}

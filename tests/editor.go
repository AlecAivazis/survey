package main

import (
	"github.com/AlecAivazis/survey"
	"github.com/AlecAivazis/survey/tests/util"
)

var answer = ""

var goodTable = []TestUtil.TestTableEntry{
	{
		"should open in editor", &survey.Editor{
			Message: "should open",
		}, &answer,
	},
	{
		"has help", &survey.Editor{
			Message: "press ? to see message",
			Help:    "Does this work?",
		}, &answer,
	},
}

func main() {
	TestUtil.RunTable(goodTable)
}

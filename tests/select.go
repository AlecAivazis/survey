package main

import (
	"gopkg.in/AlecAivazis/survey.v1"
	"gopkg.in/AlecAivazis/survey.v1/tests/util"
)

var answer = &survey.Option{}

var goodTable = []TestUtil.TestTableEntry{
	{
		"standard",
		survey.NewSingleSelect().SetMessage("Choose a color:").
			AddOption("red", nil, false).
			AddOption("blue", nil, false).
			AddOption("green", nil, false),
		&answer,
	},
	{
		"short",
		survey.NewSingleSelect().SetMessage("Choose a color:").
			AddOption("red", nil, false).
			AddOption("blue", nil, false),
		&answer,
	},
	{
		"default",
		survey.NewSingleSelect().SetMessage("Choose a color:").
			AddOption("red", nil, false).
			AddOption("blue", nil, true).
			AddOption("green", nil, false),
		&answer,
	},
	{
		"one",
		survey.NewSingleSelect().SetMessage("Choose one:").
			AddOption("hello", nil, false),
		&answer,
	},
	{
		"no help, type ?",
		survey.NewSingleSelect().SetMessage("Choose a color:").
			AddOption("red", nil, false).
			AddOption("blue", nil, false),
		&answer,
	},
	{
		"passes through bottom",
		survey.NewSingleSelect().SetMessage("Choose a color:").
			AddOption("red", nil, false).
			AddOption("blue", nil, false),
		&answer,
	},
	{
		"passes through top",
		survey.NewSingleSelect().SetMessage("Choose a color:").
			AddOption("red", nil, false).
			AddOption("blue", nil, false),
		&answer,
	},
	{
		"can navigate with j/k",
		survey.NewSingleSelect().SetMessage("Choose a color:").
			AddOption("red", nil, false).
			AddOption("blue", nil, false).
			AddOption("green", nil, false).
			SetVimMode(true),
		&answer,
	},
}

var badTable = []TestUtil.TestTableEntry{
	{
		"no options",
		survey.NewSingleSelect().SetMessage("Choose one:"),
		&answer,
	},
}

func main() {
	TestUtil.RunTable(goodTable)
	TestUtil.RunErrorTable(badTable)
}

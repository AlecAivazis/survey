package main

import (
	"github.com/AlecAivazis/survey"
	"github.com/AlecAivazis/survey/tests/util"
)

var val = ""

var table = []TestUtil.TestTableEntry{
	{
		"no default", &survey.Input{"Hello world", ""}, &val,
	},
	{
		"default", &survey.Input{"Hello world", "default"}, &val,
	},
}

func main() {
	TestUtil.RunTable(table)
}

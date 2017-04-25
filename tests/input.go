package main

import (
	"github.com/AlecAivazis/survey"
	"github.com/AlecAivazis/survey/core"
	"github.com/AlecAivazis/survey/tests/util"
)

var val = ""

var table = []TestUtil.TestTableEntry{
	{
		"no default", &survey.Input{core.Renderer{}, "Hello world", ""}, &val,
	},
	{
		"default", &survey.Input{core.Renderer{}, "Hello world", "default"}, &val,
	},
}

func main() {
	TestUtil.RunTable(table)
}

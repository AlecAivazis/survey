package main

import (
	"github.com/alecaivazis/survey"
	"github.com/alecaivazis/survey/tests/util"
)

var table = []TestUtil.TestTableEntry{
	{
		"no default", &survey.Input{"Hello world", ""},
	},
	{
		"default", &survey.Input{"Hello world", "default"},
	},
}

func main() {
	TestUtil.RunTable(table)
}

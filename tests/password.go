package main

import (
	"github.com/alecaivazis/survey"
	"github.com/alecaivazis/survey/tests/util"
)

var table = []TestUtil.TestTableEntry{
	{
		"standard", &survey.Password{"Please enter your password:"},
	},
}

func main() {
	TestUtil.RunTable(table)
}

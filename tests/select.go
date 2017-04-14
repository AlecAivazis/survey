package main

import (
	"github.com/alecaivazis/survey"
	"github.com/alecaivazis/survey/tests/util"
)

var goodTable = []TestUtil.TestTableEntry{
	{
		"standard", &survey.Select{
			Message: "Choose a color:",
			Options: []string{"red", "blue", "green"},
		},
	},
	{
		"short", &survey.Select{
			Message: "Choose a color:",
			Options: []string{"red", "blue"},
		},
	},
	{
		"default", &survey.Select{
			Message: "Choose a color (should default blue):",
			Options: []string{"red", "blue", "green"},
			Default: "blue",
		},
	},
	{
		"one", &survey.Select{
			Message: "Choose one:",
			Options: []string{"hello"},
		},
	},
}

var badTable = []TestUtil.TestTableEntry{
	{
		"no options", &survey.Select{
			Message: "Choose one:",
		},
	},
}

func main() {
	TestUtil.RunTable(goodTable)
	TestUtil.RunErrorTable(badTable)
}

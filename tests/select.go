package main

import (
	"github.com/alecaivazis/survey"
	"github.com/alecaivazis/survey/tests/util"
)

var goodTable = []TestUtil.TestTableEntry{
	{
		"standard", &survey.Choice{
			Message: "Choose a color:",
			Choices: []string{"red", "blue", "green"},
		},
	},
	{
		"short", &survey.Choice{
			Message: "Choose a color:",
			Choices: []string{"red", "blue"},
		},
	},
	{
		"default", &survey.Choice{
			Message: "Choose a color (should default blue):",
			Choices: []string{"red", "blue", "green"},
			Default: "blue",
		},
	},
	{
		"one", &survey.Choice{
			Message: "Choose one:",
			Choices: []string{"hello"},
		},
	},
}

var badTable = []TestUtil.TestTableEntry{
	{
		"no Choices", &survey.Choice{
			Message: "Choose one:",
		},
	},
}

func main() {
	TestUtil.RunTable(goodTable)
	TestUtil.RunErrorTable(badTable)
}

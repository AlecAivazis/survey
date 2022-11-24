//go:build ignore

package main

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	TestUtil "github.com/AlecAivazis/survey/v2/tests/util"
)

var answer = []string{}

var table = []TestUtil.TestTableEntry{
	{
		"standard", &survey.MultiSelect{
			Message: "What days do you prefer:",
			Options: []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"},
		}, &answer, nil,
	},
	{
		"default (sunday, tuesday)", &survey.MultiSelect{
			Message: "What days do you prefer:",
			Options: []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"},
			Default: []string{"Sunday", "Tuesday"},
		}, &answer, nil,
	},
	{
		"default not found", &survey.MultiSelect{
			Message: "What days do you prefer:",
			Options: []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"},
			Default: []string{"Sundayaa"},
		}, &answer, nil,
	},
	{
		"no help - type ?", &survey.MultiSelect{
			Message: "What days do you prefer:",
			Options: []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"},
			Default: []string{"Sundayaa"},
		}, &answer, nil,
	},
	{
		"can navigate with j/k", &survey.MultiSelect{
			Message: "What days do you prefer:",
			Options: []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"},
			Default: []string{"Sundayaa"},
		}, &answer, nil,
	},
	{
		"descriptions", &survey.MultiSelect{
			Message: "What days do you prefer:",
			Options: []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"},
			Description: func(value string, index int) string {
				return value + fmt.Sprint(index)

			},
		}, &answer, nil,
	},
}

func main() {
	TestUtil.RunTable(table)
}

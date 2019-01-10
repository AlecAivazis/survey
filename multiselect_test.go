package survey

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/Netflix/go-expect"
	"github.com/stretchr/testify/assert"
	"gopkg.in/AlecAivazis/survey.v1/core"
	"gopkg.in/AlecAivazis/survey.v1/terminal"
)

func init() {
	// disable color output for all prompts to simplify testing
	core.DisableColor = true
}

func TestMultiSelectRender(t *testing.T) {

	prompt := MultiSelect{
		Message: "Pick your words:",
		Options: []string{"foo", "bar", "baz", "buz"},
		Default: []string{"bar", "buz"},
	}

	helpfulPrompt := prompt
	helpfulPrompt.Help = "This is helpful"

	tests := []struct {
		title    string
		prompt   MultiSelect
		data     MultiSelectTemplateData
		expected string
	}{
		{
			"Test MultiSelect question output",
			prompt,
			MultiSelectTemplateData{
				SelectedIndex: 2,
				PageEntries:   prompt.Options,
				Checked:       map[string]bool{"bar": true, "buz": true},
			},
			strings.Join(
				[]string{
					fmt.Sprintf("%s Pick your words:  [Use arrows to move, type to filter]", core.QuestionIcon),
					fmt.Sprintf("  %s  foo", core.UnmarkedOptionIcon),
					fmt.Sprintf("  %s  bar", core.MarkedOptionIcon),
					fmt.Sprintf("%s %s  baz", core.SelectFocusIcon, core.UnmarkedOptionIcon),
					fmt.Sprintf("  %s  buz\n", core.MarkedOptionIcon),
				},
				"\n",
			),
		},
		{
			"Test MultiSelect answer output",
			prompt,
			MultiSelectTemplateData{
				Answer:     "foo, buz",
				ShowAnswer: true,
			},
			fmt.Sprintf("%s Pick your words: foo, buz\n", core.QuestionIcon),
		},
		{
			"Test MultiSelect question output with help hidden",
			helpfulPrompt,
			MultiSelectTemplateData{
				SelectedIndex: 2,
				PageEntries:   prompt.Options,
				Checked:       map[string]bool{"bar": true, "buz": true},
			},
			strings.Join(
				[]string{
					fmt.Sprintf("%s Pick your words:  [Use arrows to move, type to filter, %s for more help]", core.QuestionIcon, string(core.HelpInputRune)),
					fmt.Sprintf("  %s  foo", core.UnmarkedOptionIcon),
					fmt.Sprintf("  %s  bar", core.MarkedOptionIcon),
					fmt.Sprintf("%s %s  baz", core.SelectFocusIcon, core.UnmarkedOptionIcon),
					fmt.Sprintf("  %s  buz\n", core.MarkedOptionIcon),
				},
				"\n",
			),
		},
		{
			"Test MultiSelect question output with help shown",
			helpfulPrompt,
			MultiSelectTemplateData{
				SelectedIndex: 2,
				PageEntries:   prompt.Options,
				Checked:       map[string]bool{"bar": true, "buz": true},
				ShowHelp:      true,
			},
			strings.Join(
				[]string{
					fmt.Sprintf("%s This is helpful", core.HelpIcon),
					fmt.Sprintf("%s Pick your words:  [Use arrows to move, type to filter]", core.QuestionIcon),
					fmt.Sprintf("  %s  foo", core.UnmarkedOptionIcon),
					fmt.Sprintf("  %s  bar", core.MarkedOptionIcon),
					fmt.Sprintf("%s %s  baz", core.SelectFocusIcon, core.UnmarkedOptionIcon),
					fmt.Sprintf("  %s  buz\n", core.MarkedOptionIcon),
				},
				"\n",
			),
		},
	}

	for _, test := range tests {
		r, w, err := os.Pipe()
		assert.Nil(t, err, test.title)

		test.prompt.WithStdio(terminal.Stdio{Out: w})
		test.data.MultiSelect = test.prompt
		err = test.prompt.Render(
			MultiSelectQuestionTemplate,
			test.data,
		)
		assert.Nil(t, err, test.title)

		w.Close()
		var buf bytes.Buffer
		io.Copy(&buf, r)

		assert.Contains(t, buf.String(), test.expected, test.title)
	}
}

func TestMultiSelectPrompt(t *testing.T) {
	tests := []PromptTest{
		{
			"Test MultiSelect prompt interaction",
			&MultiSelect{
				Message: "What days do you prefer:",
				Options: []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"},
			},
			func(c *expect.Console) {
				c.ExpectString("What days do you prefer:  [Use arrows to move, type to filter]")
				// Select Monday.
				c.Send(string(terminal.KeyArrowDown))
				c.SendLine(" ")
				c.ExpectEOF()
			},
			[]string{"Monday"},
		},
		{
			"Test MultiSelect prompt interaction with default",
			&MultiSelect{
				Message: "What days do you prefer:",
				Options: []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"},
				Default: []string{"Tuesday", "Thursday"},
			},
			func(c *expect.Console) {
				c.ExpectString("What days do you prefer:  [Use arrows to move, type to filter]")
				c.SendLine("")
				c.ExpectEOF()
			},
			[]string{"Tuesday", "Thursday"},
		},
		{
			"Test MultiSelect prompt interaction overriding default",
			&MultiSelect{
				Message: "What days do you prefer:",
				Options: []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"},
				Default: []string{"Tuesday", "Thursday"},
			},
			func(c *expect.Console) {
				c.ExpectString("What days do you prefer:  [Use arrows to move, type to filter]")
				// Deselect Tuesday.
				c.Send(string(terminal.KeyArrowDown))
				c.Send(string(terminal.KeyArrowDown))
				c.SendLine(" ")
				c.ExpectEOF()
			},
			[]string{"Thursday"},
		},
		{
			"Test MultiSelect prompt interaction and prompt for help",
			&MultiSelect{
				Message: "What days do you prefer:",
				Options: []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"},
				Help:    "Saturday is best",
			},
			func(c *expect.Console) {
				c.ExpectString("What days do you prefer:  [Use arrows to move, type to filter, ? for more help]")
				c.Send("?")
				c.ExpectString("Saturday is best")
				// Select Saturday
				c.Send(string(terminal.KeyArrowUp))
				c.SendLine(" ")
				c.ExpectEOF()
			},
			[]string{"Saturday"},
		},
		{
			"Test MultiSelect prompt interaction with page size",
			&MultiSelect{
				Message:  "What days do you prefer:",
				Options:  []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"},
				PageSize: 1,
			},
			func(c *expect.Console) {
				c.ExpectString("What days do you prefer:  [Use arrows to move, type to filter]")
				// Select Monday.
				c.Send(string(terminal.KeyArrowDown))
				c.SendLine(" ")
				c.ExpectEOF()
			},
			[]string{"Monday"},
		},
		{
			"Test MultiSelect prompt interaction with vim mode",
			&MultiSelect{
				Message: "What days do you prefer:",
				Options: []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"},
				VimMode: true,
			},
			func(c *expect.Console) {
				c.ExpectString("What days do you prefer:  [Use arrows to move, type to filter]")
				// Select Tuesday.
				c.Send("jj ")
				// Select Thursday.
				c.Send("jj ")
				// Select Saturday.
				c.Send("jj ")
				c.SendLine("")
				c.ExpectEOF()
			},
			[]string{"Tuesday", "Thursday", "Saturday"},
		},
		{
			"Test MultiSelect prompt interaction with filter",
			&MultiSelect{
				Message: "What days do you prefer:",
				Options: []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"},
			},
			func(c *expect.Console) {
				c.ExpectString("What days do you prefer:  [Use arrows to move, type to filter]")
				// Filter down to Tuesday.
				c.Send("Tues")
				// Select Tuesday.
				c.Send(" ")
				c.SendLine("")
				c.ExpectEOF()
			},
			[]string{"Tuesday"},
		},
		{
			"Test MultiSelect prompt interaction with filter is case-insensitive",
			&MultiSelect{
				Message: "What days do you prefer:",
				Options: []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"},
			},
			func(c *expect.Console) {
				c.ExpectString("What days do you prefer:  [Use arrows to move, type to filter]")
				// Filter down to Tuesday.
				c.Send("tues")
				// Select Tuesday.
				c.Send(" ")
				c.SendLine("")
				c.ExpectEOF()
			},
			[]string{"Tuesday"},
		},
		{
			"Test MultiSelect prompt interaction with custom filter",
			&MultiSelect{
				Message: "What days do you prefer:",
				Options: []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"},
				FilterFn: func(filter string, options []string) (filtered []string) {
					result := DefaultFilterFn(filter, options)
					for _, v := range result {
						if len(v) >= 7 {
							filtered = append(filtered, v)
						}
					}
					return
				},
			},
			func(c *expect.Console) {
				c.ExpectString("What days do you prefer:")
				// Filter down to days which names are longer than 7 runes
				c.Send("day")
				// Select Wednesday.
				c.Send(string(terminal.KeyArrowDown))
				c.SendLine(" ")
				c.ExpectEOF()
			},
			[]string{"Wednesday"},
		},
		{
			"Test MultiSelect clears input on select",
			&MultiSelect{
				Message: "What days do you prefer:",
				Options: []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"},
			},
			func(c *expect.Console) {
				c.ExpectString("What days do you prefer:  [Use arrows to move, type to filter]")
				// Filter down to Tuesday.
				c.Send("Tues")
				// Select Tuesday.
				c.Send(" ")
				// Filter down to Tuesday.
				c.Send("Tues")
				// Deselect Tuesday.
				c.Send(" ")
				c.SendLine("")
				c.ExpectEOF()
			},
			[]string{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			RunPromptTest(t, test)
		})
	}
}

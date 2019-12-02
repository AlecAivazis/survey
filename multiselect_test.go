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

	"github.com/AlecAivazis/survey/v2/core"
	"github.com/AlecAivazis/survey/v2/terminal"
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

	pagePrompt := MultiSelect{
		Message:  "Pick your words:",
		Options:  []string{"foo", "bar", "baz", "buz"},
		PageSize: 2,
	}

	tests := []struct {
		title    string
		prompt   MultiSelect
		data     MultiSelectTemplateData
		expected string
	}{
		{
			"question output",
			prompt,
			MultiSelectTemplateData{
				SelectedIndex: 2,
				PageEntries:   core.OptionAnswerList(prompt.Options),
				Checked:       map[int]bool{1: true, 3: true},
			},
			strings.Join(
				[]string{
					fmt.Sprintf("%s Pick your words:  [Use arrows to move, enter to select, type to filter]", defaultIcons().Question.Text),
					fmt.Sprintf("  %s  foo", defaultIcons().UnmarkedOption.Text),
					fmt.Sprintf("  %s  bar", defaultIcons().MarkedOption.Text),
					fmt.Sprintf("%s %s  baz", defaultIcons().SelectFocus.Text, defaultIcons().UnmarkedOption.Text),
					fmt.Sprintf("  %s  buz\n", defaultIcons().MarkedOption.Text),
				},
				"\n",
			),
		},
		{
			"answer output",
			prompt,
			MultiSelectTemplateData{
				Answer:     "foo, buz",
				ShowAnswer: true,
			},
			fmt.Sprintf("%s Pick your words: foo, buz\n", defaultIcons().Question.Text),
		},
		{
			"help hidden",
			helpfulPrompt,
			MultiSelectTemplateData{
				SelectedIndex: 2,
				PageEntries:   core.OptionAnswerList(prompt.Options),
				Checked:       map[int]bool{1: true, 3: true},
			},
			strings.Join(
				[]string{
					fmt.Sprintf("%s Pick your words:  [Use arrows to move, enter to select, type to filter, %s for more help]", defaultIcons().Question.Text, string(defaultPromptConfig().HelpInput)),
					fmt.Sprintf("  %s  foo", defaultIcons().UnmarkedOption.Text),
					fmt.Sprintf("  %s  bar", defaultIcons().MarkedOption.Text),
					fmt.Sprintf("%s %s  baz", defaultIcons().SelectFocus.Text, defaultIcons().UnmarkedOption.Text),
					fmt.Sprintf("  %s  buz\n", defaultIcons().MarkedOption.Text),
				},
				"\n",
			),
		},
		{
			"question outputhelp shown",
			helpfulPrompt,
			MultiSelectTemplateData{
				SelectedIndex: 2,
				PageEntries:   core.OptionAnswerList(prompt.Options),
				Checked:       map[int]bool{1: true, 3: true},
				ShowHelp:      true,
			},
			strings.Join(
				[]string{
					fmt.Sprintf("%s This is helpful", defaultIcons().Help.Text),
					fmt.Sprintf("%s Pick your words:  [Use arrows to move, enter to select, type to filter]", defaultIcons().Question.Text),
					fmt.Sprintf("  %s  foo", defaultIcons().UnmarkedOption.Text),
					fmt.Sprintf("  %s  bar", defaultIcons().MarkedOption.Text),
					fmt.Sprintf("%s %s  baz", defaultIcons().SelectFocus.Text, defaultIcons().UnmarkedOption.Text),
					fmt.Sprintf("  %s  buz\n", defaultIcons().MarkedOption.Text),
				},
				"\n",
			),
		},
		{
			"marked on paginating",
			pagePrompt,
			MultiSelectTemplateData{
				SelectedIndex: 0,
				PageEntries:   core.OptionAnswerList(pagePrompt.Options)[1:3], /* show unmarked items(bar, baz)*/
				Checked:       map[int]bool{0: true},                          /* foo marked */
			},
			strings.Join(
				[]string{
					fmt.Sprintf("%s Pick your words:  [Use arrows to move, enter to select, type to filter]", defaultIcons().Question.Text),
					fmt.Sprintf("%s %s  bar", defaultIcons().SelectFocus.Text, defaultIcons().UnmarkedOption.Text),
					fmt.Sprintf("  %s  baz", defaultIcons().UnmarkedOption.Text),
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

		// set the icon set
		test.data.Config = defaultPromptConfig()

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
			"basic interaction",
			&MultiSelect{
				Message: "What days do you prefer:",
				Options: []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"},
			},
			func(c *expect.Console) {
				c.ExpectString("What days do you prefer:  [Use arrows to move, enter to select, type to filter]")
				// Select Monday.
				c.Send(string(terminal.KeyArrowDown))
				c.SendLine(" ")
				c.ExpectEOF()
			},
			[]core.OptionAnswer{core.OptionAnswer{Value: "Monday", Index: 1}},
		},
		{
			"default value as []string",
			&MultiSelect{
				Message: "What days do you prefer:",
				Options: []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"},
				Default: []string{"Tuesday", "Thursday"},
			},
			func(c *expect.Console) {
				c.ExpectString("What days do you prefer:  [Use arrows to move, enter to select, type to filter]")
				c.SendLine("")
				c.ExpectEOF()
			},
			[]core.OptionAnswer{
				core.OptionAnswer{Value: "Tuesday", Index: 2},
				core.OptionAnswer{Value: "Thursday", Index: 4},
			},
		},
		{
			"default value as []int",
			&MultiSelect{
				Message: "What days do you prefer:",
				Options: []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"},
				Default: []int{2, 4},
			},
			func(c *expect.Console) {
				c.ExpectString("What days do you prefer:  [Use arrows to move, enter to select, type to filter]")
				c.SendLine("")
				c.ExpectEOF()
			},
			[]core.OptionAnswer{
				core.OptionAnswer{Value: "Tuesday", Index: 2},
				core.OptionAnswer{Value: "Thursday", Index: 4},
			},
		},
		{
			"overriding default",
			&MultiSelect{
				Message: "What days do you prefer:",
				Options: []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"},
				Default: []string{"Tuesday", "Thursday"},
			},
			func(c *expect.Console) {
				c.ExpectString("What days do you prefer:  [Use arrows to move, enter to select, type to filter]")
				// Deselect Tuesday.
				c.Send(string(terminal.KeyArrowDown))
				c.Send(string(terminal.KeyArrowDown))
				c.SendLine(" ")
				c.ExpectEOF()
			},
			[]core.OptionAnswer{core.OptionAnswer{Value: "Thursday", Index: 4}},
		},
		{
			"prompt for help",
			&MultiSelect{
				Message: "What days do you prefer:",
				Options: []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"},
				Help:    "Saturday is best",
			},
			func(c *expect.Console) {
				c.ExpectString("What days do you prefer:  [Use arrows to move, enter to select, type to filter, ? for more help]")
				c.Send("?")
				c.ExpectString("Saturday is best")
				// Select Saturday
				c.Send(string(terminal.KeyArrowUp))
				c.SendLine(" ")
				c.ExpectEOF()
			},
			[]core.OptionAnswer{core.OptionAnswer{Value: "Saturday", Index: 6}},
		},
		{
			"page size",
			&MultiSelect{
				Message:  "What days do you prefer:",
				Options:  []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"},
				PageSize: 1,
			},
			func(c *expect.Console) {
				c.ExpectString("What days do you prefer:  [Use arrows to move, enter to select, type to filter]")
				// Select Monday.
				c.Send(string(terminal.KeyArrowDown))
				c.SendLine(" ")
				c.ExpectEOF()
			},
			[]core.OptionAnswer{core.OptionAnswer{Value: "Monday", Index: 1}},
		},
		{
			"vim mode",
			&MultiSelect{
				Message: "What days do you prefer:",
				Options: []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"},
				VimMode: true,
			},
			func(c *expect.Console) {
				c.ExpectString("What days do you prefer:  [Use arrows to move, enter to select, type to filter]")
				// Select Tuesday.
				c.Send("jj ")
				// Select Thursday.
				c.Send("jj ")
				// Select Saturday.
				c.Send("jj ")
				c.SendLine("")
				c.ExpectEOF()
			},
			[]core.OptionAnswer{
				core.OptionAnswer{Value: "Tuesday", Index: 2},
				core.OptionAnswer{Value: "Thursday", Index: 4},
				core.OptionAnswer{Value: "Saturday", Index: 6},
			},
		},
		{
			"filter interaction",
			&MultiSelect{
				Message: "What days do you prefer:",
				Options: []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"},
			},
			func(c *expect.Console) {
				c.ExpectString("What days do you prefer:  [Use arrows to move, enter to select, type to filter]")
				// Filter down to Tuesday.
				c.Send("Tues")
				// Select Tuesday.
				c.Send(" ")
				c.SendLine("")
				c.ExpectEOF()
			},
			[]core.OptionAnswer{core.OptionAnswer{Value: "Tuesday", Index: 2}},
		},
		{
			"filter is case-insensitive",
			&MultiSelect{
				Message: "What days do you prefer:",
				Options: []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"},
			},
			func(c *expect.Console) {
				c.ExpectString("What days do you prefer:  [Use arrows to move, enter to select, type to filter]")
				// Filter down to Tuesday.
				c.Send("tues")
				// Select Tuesday.
				c.Send(" ")
				c.SendLine("")
				c.ExpectEOF()
			},
			[]core.OptionAnswer{core.OptionAnswer{Value: "Tuesday", Index: 2}},
		},
		{
			"custom filter",
			&MultiSelect{
				Message: "What days do you prefer:",
				Options: []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"},
				Filter: func(filterValue string, optValue string, index int) bool {
					return strings.Contains(optValue, filterValue) && len(optValue) >= 7
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
			[]core.OptionAnswer{core.OptionAnswer{Value: "Wednesday", Index: 3}},
		},
		{
			"clears input on select",
			&MultiSelect{
				Message: "What days do you prefer:",
				Options: []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"},
			},
			func(c *expect.Console) {
				c.ExpectString("What days do you prefer:  [Use arrows to move, enter to select, type to filter]")
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
			[]core.OptionAnswer{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			RunPromptTest(t, test)
		})
	}
}

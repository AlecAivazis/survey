package survey

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

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

	descriptions := []string{"oof", "rab", "zab", "zub"}

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
					fmt.Sprintf("%s Pick your words:  [Use arrows to move, space to select, <right> to all, <left> to none, type to filter]", defaultIcons().Question.Text),
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
					fmt.Sprintf("%s Pick your words:  [Use arrows to move, space to select, <right> to all, <left> to none, type to filter, %s for more help]", defaultIcons().Question.Text, string(defaultPromptConfig().HelpInput)),
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
					fmt.Sprintf("%s Pick your words:  [Use arrows to move, space to select, <right> to all, <left> to none, type to filter]", defaultIcons().Question.Text),
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
					fmt.Sprintf("%s Pick your words:  [Use arrows to move, space to select, <right> to all, <left> to none, type to filter]", defaultIcons().Question.Text),
					fmt.Sprintf("%s %s  bar", defaultIcons().SelectFocus.Text, defaultIcons().UnmarkedOption.Text),
					fmt.Sprintf("  %s  baz", defaultIcons().UnmarkedOption.Text),
				},
				"\n",
			),
		},
		{
			"description all",
			prompt,
			MultiSelectTemplateData{
				SelectedIndex: 2,
				PageEntries:   core.OptionAnswerList(prompt.Options),
				Checked:       map[int]bool{1: true, 3: true},
				Description: func(value string, index int) string {
					return descriptions[index]
				},
			},
			strings.Join(
				[]string{
					fmt.Sprintf("%s Pick your words:  [Use arrows to move, space to select, <right> to all, <left> to none, type to filter]", defaultIcons().Question.Text),
					fmt.Sprintf("  %s  foo - oof", defaultIcons().UnmarkedOption.Text),
					fmt.Sprintf("  %s  bar - rab", defaultIcons().MarkedOption.Text),
					fmt.Sprintf("%s %s  baz - zab", defaultIcons().SelectFocus.Text, defaultIcons().UnmarkedOption.Text),
					fmt.Sprintf("  %s  buz - zub\n", defaultIcons().MarkedOption.Text),
				},
				"\n",
			),
		},
		{
			"description even",
			prompt,
			MultiSelectTemplateData{
				SelectedIndex: 2,
				PageEntries:   core.OptionAnswerList(prompt.Options),
				Checked:       map[int]bool{1: true, 3: true},
				Description: func(value string, index int) string {

					if index%2 == 0 {
						return descriptions[index]
					}

					return ""
				},
			},
			strings.Join(
				[]string{
					fmt.Sprintf("%s Pick your words:  [Use arrows to move, space to select, <right> to all, <left> to none, type to filter]", defaultIcons().Question.Text),
					fmt.Sprintf("  %s  foo - oof", defaultIcons().UnmarkedOption.Text),
					fmt.Sprintf("  %s  bar", defaultIcons().MarkedOption.Text),
					fmt.Sprintf("%s %s  baz - zab", defaultIcons().SelectFocus.Text, defaultIcons().UnmarkedOption.Text),
					fmt.Sprintf("  %s  buz\n", defaultIcons().MarkedOption.Text),
				},
				"\n",
			),
		},
		{
			"description never",
			prompt,
			MultiSelectTemplateData{
				SelectedIndex: 2,
				PageEntries:   core.OptionAnswerList(prompt.Options),
				Checked:       map[int]bool{1: true, 3: true},
				Description: func(value string, index int) string {
					return ""
				},
			},
			strings.Join(
				[]string{
					fmt.Sprintf("%s Pick your words:  [Use arrows to move, space to select, <right> to all, <left> to none, type to filter]", defaultIcons().Question.Text),
					fmt.Sprintf("  %s  foo", defaultIcons().UnmarkedOption.Text),
					fmt.Sprintf("  %s  bar", defaultIcons().MarkedOption.Text),
					fmt.Sprintf("%s %s  baz", defaultIcons().SelectFocus.Text, defaultIcons().UnmarkedOption.Text),
					fmt.Sprintf("  %s  buz\n", defaultIcons().MarkedOption.Text),
				},
				"\n",
			),
		},
		{
			"description repeat value",
			prompt,
			MultiSelectTemplateData{
				SelectedIndex: 2,
				PageEntries:   core.OptionAnswerList(prompt.Options),
				Checked:       map[int]bool{1: true, 3: true},
				Description: func(value string, index int) string {
					return value
				},
			},
			strings.Join(
				[]string{
					fmt.Sprintf("%s Pick your words:  [Use arrows to move, space to select, <right> to all, <left> to none, type to filter]", defaultIcons().Question.Text),
					fmt.Sprintf("  %s  foo - foo", defaultIcons().UnmarkedOption.Text),
					fmt.Sprintf("  %s  bar - bar", defaultIcons().MarkedOption.Text),
					fmt.Sprintf("%s %s  baz - baz", defaultIcons().SelectFocus.Text, defaultIcons().UnmarkedOption.Text),
					fmt.Sprintf("  %s  buz - buz\n", defaultIcons().MarkedOption.Text),
				},
				"\n",
			),
		},
		{
			"description print index",
			prompt,
			MultiSelectTemplateData{
				SelectedIndex: 2,
				PageEntries:   core.OptionAnswerList(prompt.Options),
				Checked:       map[int]bool{1: true, 3: true},
				Description: func(value string, index int) string {
					return fmt.Sprint(index)
				},
			},
			strings.Join(
				[]string{
					fmt.Sprintf("%s Pick your words:  [Use arrows to move, space to select, <right> to all, <left> to none, type to filter]", defaultIcons().Question.Text),
					fmt.Sprintf("  %s  foo - 0", defaultIcons().UnmarkedOption.Text),
					fmt.Sprintf("  %s  bar - 1", defaultIcons().MarkedOption.Text),
					fmt.Sprintf("%s %s  baz - 2", defaultIcons().SelectFocus.Text, defaultIcons().UnmarkedOption.Text),
					fmt.Sprintf("  %s  buz - 3\n", defaultIcons().MarkedOption.Text),
				},
				"\n",
			),
		},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			r, w, err := os.Pipe()
			assert.NoError(t, err)

			test.prompt.WithStdio(terminal.Stdio{Out: w})
			test.data.MultiSelect = test.prompt

			// set the icon set
			test.data.Config = defaultPromptConfig()

			err = test.prompt.Render(
				MultiSelectQuestionTemplate,
				test.data,
			)
			assert.NoError(t, err)

			assert.NoError(t, w.Close())
			var buf bytes.Buffer
			_, err = io.Copy(&buf, r)
			assert.NoError(t, err)

			assert.Contains(t, buf.String(), test.expected)
		})
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
			func(c expectConsole) {
				c.ExpectString("What days do you prefer:  [Use arrows to move, space to select, <right> to all, <left> to none, type to filter]")
				// Select Monday.
				c.Send(string(terminal.KeyArrowDown))
				c.SendLine(" ")
				c.ExpectEOF()
			},
			[]core.OptionAnswer{{Value: "Monday", Index: 1}},
		},
		{
			"cycle to next when tab send",
			&MultiSelect{
				Message: "What days do you prefer:",
				Options: []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"},
			},
			func(c expectConsole) {
				c.ExpectString("What days do you prefer:  [Use arrows to move, space to select, <right> to all, <left> to none, type to filter]")
				// Select Monday.
				c.Send(string(terminal.KeyTab))
				c.Send(" ")
				c.Send(string(terminal.KeyArrowDown))
				c.SendLine(" ")
				c.ExpectEOF()
			},
			[]core.OptionAnswer{
				{Value: "Monday", Index: 1},
				{Value: "Tuesday", Index: 2},
			},
		},
		{
			"default value as []string",
			&MultiSelect{
				Message: "What days do you prefer:",
				Options: []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"},
				Default: []string{"Tuesday", "Thursday"},
			},
			func(c expectConsole) {
				c.ExpectString("What days do you prefer:  [Use arrows to move, space to select, <right> to all, <left> to none, type to filter]")
				c.SendLine("")
				c.ExpectEOF()
			},
			[]core.OptionAnswer{
				{Value: "Tuesday", Index: 2},
				{Value: "Thursday", Index: 4},
			},
		},
		{
			"default value as []int",
			&MultiSelect{
				Message: "What days do you prefer:",
				Options: []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"},
				Default: []int{2, 4},
			},
			func(c expectConsole) {
				c.ExpectString("What days do you prefer:  [Use arrows to move, space to select, <right> to all, <left> to none, type to filter]")
				c.SendLine("")
				c.ExpectEOF()
			},
			[]core.OptionAnswer{
				{Value: "Tuesday", Index: 2},
				{Value: "Thursday", Index: 4},
			},
		},
		{
			"overriding default",
			&MultiSelect{
				Message: "What days do you prefer:",
				Options: []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"},
				Default: []string{"Tuesday", "Thursday"},
			},
			func(c expectConsole) {
				c.ExpectString("What days do you prefer:  [Use arrows to move, space to select, <right> to all, <left> to none, type to filter]")
				// Deselect Tuesday.
				c.Send(string(terminal.KeyArrowDown))
				c.Send(string(terminal.KeyArrowDown))
				c.SendLine(" ")
				c.ExpectEOF()
			},
			[]core.OptionAnswer{{Value: "Thursday", Index: 4}},
		},
		{
			"prompt for help",
			&MultiSelect{
				Message: "What days do you prefer:",
				Options: []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"},
				Help:    "Saturday is best",
			},
			func(c expectConsole) {
				c.ExpectString("What days do you prefer:  [Use arrows to move, space to select, <right> to all, <left> to none, type to filter, ? for more help]")
				c.Send("?")
				c.ExpectString("Saturday is best")
				// Select Saturday
				c.Send(string(terminal.KeyArrowUp))
				c.SendLine(" ")
				c.ExpectEOF()
			},
			[]core.OptionAnswer{{Value: "Saturday", Index: 6}},
		},
		{
			"page size",
			&MultiSelect{
				Message:  "What days do you prefer:",
				Options:  []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"},
				PageSize: 1,
			},
			func(c expectConsole) {
				c.ExpectString("What days do you prefer:  [Use arrows to move, space to select, <right> to all, <left> to none, type to filter]")
				// Select Monday.
				c.Send(string(terminal.KeyArrowDown))
				c.SendLine(" ")
				c.ExpectEOF()
			},
			[]core.OptionAnswer{{Value: "Monday", Index: 1}},
		},
		{
			"vim mode",
			&MultiSelect{
				Message: "What days do you prefer:",
				Options: []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"},
				VimMode: true,
			},
			func(c expectConsole) {
				c.ExpectString("What days do you prefer:  [Use arrows to move, space to select, <right> to all, <left> to none, type to filter]")
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
				{Value: "Tuesday", Index: 2},
				{Value: "Thursday", Index: 4},
				{Value: "Saturday", Index: 6},
			},
		},
		{
			"filter interaction",
			&MultiSelect{
				Message: "What days do you prefer:",
				Options: []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"},
			},
			func(c expectConsole) {
				c.ExpectString("What days do you prefer:  [Use arrows to move, space to select, <right> to all, <left> to none, type to filter]")
				// Filter down to Tuesday.
				c.Send("Tues")
				// Select Tuesday.
				c.Send(" ")
				c.SendLine("")
				c.ExpectEOF()
			},
			[]core.OptionAnswer{{Value: "Tuesday", Index: 2}},
		},
		{
			"filter is case-insensitive",
			&MultiSelect{
				Message: "What days do you prefer:",
				Options: []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"},
			},
			func(c expectConsole) {
				c.ExpectString("What days do you prefer:  [Use arrows to move, space to select, <right> to all, <left> to none, type to filter]")
				// Filter down to Tuesday.
				c.Send("tues")
				// Select Tuesday.
				c.Send(" ")
				c.SendLine("")
				c.ExpectEOF()
			},
			[]core.OptionAnswer{{Value: "Tuesday", Index: 2}},
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
			func(c expectConsole) {
				c.ExpectString("What days do you prefer:")
				// Filter down to days which names are longer than 7 runes
				c.Send("day")
				// Select Wednesday.
				c.Send(string(terminal.KeyArrowDown))
				c.SendLine(" ")
				c.ExpectEOF()
			},
			[]core.OptionAnswer{{Value: "Wednesday", Index: 3}},
		},
		{
			"clears input on select",
			&MultiSelect{
				Message: "What days do you prefer:",
				Options: []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"},
			},
			func(c expectConsole) {
				c.ExpectString("What days do you prefer:  [Use arrows to move, space to select, <right> to all, <left> to none, type to filter]")
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
		{
			"select all",
			&MultiSelect{
				Message: "What days do you prefer:",
				Options: []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"},
			},
			func(c expectConsole) {
				c.ExpectString("What days do you prefer:  [Use arrows to move, space to select, <right> to all, <left> to none, type to filter]")
				// Select all
				c.Send(string(terminal.KeyArrowRight))
				c.SendLine("")
				c.ExpectEOF()
			},
			[]core.OptionAnswer{
				{Value: "Sunday", Index: 0},
				{Value: "Monday", Index: 1},
				{Value: "Tuesday", Index: 2},
				{Value: "Wednesday", Index: 3},
				{Value: "Thursday", Index: 4},
				{Value: "Friday", Index: 5},
				{Value: "Saturday", Index: 6},
			},
		},
		{
			"select none",
			&MultiSelect{
				Message: "What days do you prefer:",
				Options: []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"},
			},
			func(c expectConsole) {
				c.ExpectString("What days do you prefer:  [Use arrows to move, space to select, <right> to all, <left> to none, type to filter]")
				// Select first
				c.Send(" ")
				// Select second
				c.Send(string(terminal.KeyArrowDown))
				c.Send(" ")
				// Deselect all
				c.Send(string(terminal.KeyArrowLeft))
				c.SendLine("")
				c.ExpectEOF()
			},
			[]core.OptionAnswer{},
		},
		{
			"select all with filter",
			&MultiSelect{
				Message: "What days do you prefer:",
				Options: []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"},
			},
			func(c expectConsole) {
				c.ExpectString("What days do you prefer:  [Use arrows to move, space to select, <right> to all, <left> to none, type to filter]")
				// Send filter
				c.Send("tu")
				// Select all
				c.Send(string(terminal.KeyArrowRight))
				c.SendLine("")
				c.ExpectEOF()
			},
			[]core.OptionAnswer{
				{Value: "Tuesday", Index: 2},
				{Value: "Saturday", Index: 6},
			},
		},
		{
			"select all with filter and select others without filter",
			&MultiSelect{
				Message: "What days do you prefer:",
				Options: []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"},
			},
			func(c expectConsole) {
				c.ExpectString("What days do you prefer:  [Use arrows to move, space to select, <right> to all, <left> to none, type to filter]")
				// Select first
				c.Send(" ")
				// Select second
				c.Send(string(terminal.KeyArrowDown))
				c.Send(" ")
				// Send filter
				c.Send("tu")
				// Select all
				c.Send(string(terminal.KeyArrowRight))
				c.SendLine("")
				c.ExpectEOF()
			},
			[]core.OptionAnswer{
				{Value: "Sunday", Index: 0},
				{Value: "Monday", Index: 1},
				{Value: "Tuesday", Index: 2},
				{Value: "Saturday", Index: 6},
			},
		},
		{
			"select all with filter and deselect one without filter",
			&MultiSelect{
				Message: "What days do you prefer:",
				Options: []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"},
			},
			func(c expectConsole) {
				c.ExpectString("What days do you prefer:  [Use arrows to move, space to select, <right> to all, <left> to none, type to filter]")
				// Send filter
				c.Send("tu")
				// Select all
				c.Send(string(terminal.KeyArrowRight))
				// Deselect second
				c.Send(string(terminal.KeyArrowDown))
				c.Send(string(terminal.KeyArrowDown))
				c.Send(" ")
				c.SendLine("")
				c.ExpectEOF()
			},
			[]core.OptionAnswer{
				{Value: "Saturday", Index: 6},
			},
		},
		{
			"delete filter word",
			&MultiSelect{
				Message: "What days do you prefer:",
				Options: []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"},
			},
			func(c expectConsole) {
				c.ExpectString("What days do you prefer:  [Use arrows to move, space to select, <right> to all, <left> to none, type to filter]")
				// Filter down to 'Sunday'
				c.Send("su")
				// Delete 'u'
				c.Send(string(terminal.KeyDelete))
				// Filter down to 'Saturday'
				c.Send("at")
				// Select 'Saturday'
				c.Send(string(terminal.KeyArrowDown))
				c.Send(" ")
				c.SendLine("")
				c.ExpectEOF()
			},
			[]core.OptionAnswer{
				{Value: "Saturday", Index: 6},
			},
		},
		{
			"delete filter word in rune",
			&MultiSelect{
				Message: "今天中午吃什么？",
				Options: []string{"青椒牛肉丝", "小炒肉", "小煎鸡"},
			},
			func(c expectConsole) {
				c.ExpectString("今天中午吃什么？  [Use arrows to move, space to select, <right> to all, <left> to none, type to filter]")
				// Filter down to 小炒肉.
				c.Send("小炒")
				// Filter down to 小炒肉 and 小煎鸡.
				c.Send(string(terminal.KeyBackspace))
				// Filter down to 小煎鸡.
				c.Send("煎")
				// Select 小煎鸡.
				c.Send(" ")
				c.SendLine("")
				c.ExpectEOF()
			},
			[]core.OptionAnswer{
				{Value: "小煎鸡", Index: 2},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			RunPromptTest(t, test)
		})
	}
}

func TestMultiSelectPromptKeepFilter(t *testing.T) {
	tests := []PromptTest{
		{
			"multi select with filter keep",
			&MultiSelect{
				Message: "What color do you prefer:",
				Options: []string{"green", "red", "light-green", "blue", "black", "yellow", "purple"},
			},
			func(c expectConsole) {
				c.ExpectString("What color do you prefer:  [Use arrows to move, space to select, <right> to all, <left> to none, type to filter]")
				// Filter down to green
				c.Send("green")
				// Select green.
				c.Send(" ")
				// Select light-green.
				c.Send(string(terminal.KeyArrowDown))
				c.Send(" ")
				c.SendLine("")
				c.ExpectEOF()
			},
			[]core.OptionAnswer{
				{Value: "green", Index: 0},
				{Value: "light-green", Index: 2},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			RunPromptTestKeepFilter(t, test)
		})
	}
}

func TestMultiSelectPromptRemoveSelectAll(t *testing.T) {
	tests := []PromptTest{
		{
			"multi select with remove select all option",
			&MultiSelect{
				Message: "What color do you prefer:",
				Options: []string{"green", "red", "light-green", "blue", "black", "yellow", "purple"},
			},
			func(c expectConsole) {
				c.ExpectString("What color do you prefer:  [Use arrows to move, space to select, <left> to none, type to filter]")
				// Select the first option "green"
				c.Send(" ")

				// attempt to select all (this shouldn't do anything)
				c.Send(string(terminal.KeyArrowRight))

				// end the session
				c.SendLine("")
				c.ExpectEOF()
			},
			[]core.OptionAnswer{ // we should only have one option selected, not all of them
				{Value: "green", Index: 0},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			RunPromptTestRemoveSelectAll(t, test)
		})
	}
}

func TestMultiSelectPromptRemoveSelectNone(t *testing.T) {
	tests := []PromptTest{
		{
			"multi select with remove select none option",
			&MultiSelect{
				Message: "What color do you prefer:",
				Options: []string{"green", "red", "light-green", "blue", "black", "yellow", "purple"},
			},
			func(c expectConsole) {
				c.ExpectString("What color do you prefer:  [Use arrows to move, space to select, <right> to all, type to filter]")
				// Select the first option "green"
				c.Send(" ")

				// attempt to unselect all (this shouldn't do anything)
				c.Send(string(terminal.KeyArrowLeft))

				// end the session
				c.SendLine("")
				c.ExpectEOF()
			},
			[]core.OptionAnswer{ // we should only have one option selected, not all of them
				{Value: "green", Index: 0},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			RunPromptTestRemoveSelectNone(t, test)
		})
	}
}

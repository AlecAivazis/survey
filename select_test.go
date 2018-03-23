package survey

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/AlecAivazis/survey.v1/core"
	"gopkg.in/AlecAivazis/survey.v1/terminal"
)

func init() {
	// disable color output for all prompts to simplify testing
	core.DisableColor = true
}

func TestSelectRender(t *testing.T) {
	prompt := NewSingleSelect()
	prompt.SetMessage("Pick your word:").
		AddOption("foo", nil, false).
		AddOption("bar", nil, false).
		AddOption("baz", nil, true).
		AddOption("buz", nil, false)

	helpfulPrompt := prompt
	helpfulPrompt.SetHelp("This is helpful")

	tests := []struct {
		title    string
		prompt   *Select
		data     SelectTemplateData
		expected string
	}{
		{
			"Test Select question output",
			prompt,
			SelectTemplateData{SelectedIndex: 2, PageEntries: prompt.Options},
			`? Pick your word:  [Use arrows to move, type to filter]
  foo
  bar
❯ baz
  buz
`,
		},
		{
			"Test Select answer output",
			prompt,
			SelectTemplateData{Answer: prompt.Options[3], ShowAnswer: true, PageEntries: prompt.Options},
			"? Pick your word: buz\n",
		},
		{
			"Test Select question output with help hidden",
			helpfulPrompt,
			SelectTemplateData{SelectedIndex: 2, PageEntries: prompt.Options},
			`? Pick your word:  [Use arrows to move, type to filter, ? for more help]
  foo
  bar
❯ baz
  buz
`,
		},
		{
			"Test Select question output with help shown",
			helpfulPrompt,
			SelectTemplateData{SelectedIndex: 2, ShowHelp: true, PageEntries: prompt.Options},
			`ⓘ This is helpful
? Pick your word:  [Use arrows to move, type to filter]
  foo
  bar
❯ baz
  buz
`,
		},
	}

	outputBuffer := bytes.NewBufferString("")
	terminal.Stdout = outputBuffer

	for _, test := range tests {
		outputBuffer.Reset()
		test.data.Select = *test.prompt
		err := test.prompt.Render(
			SelectQuestionTemplate,
			test.data,
		)
		assert.Nil(t, err, test.title)
		assert.Equal(t, test.expected, outputBuffer.String(), test.title)
	}
}

func TestSelectInterfaceValues(t *testing.T) {
	type value struct {
		Item string
		Other int
	}
	prompt := NewSingleSelect()
	prompt.SetMessage("Pick your word:").
		AddOption("foo", value{"foo", 0}, false).
		AddOption("bar", value{"bar", 5}, false).
		AddOption("baz", value{"baz", 100}, true).
		AddOption("buz", value{"buz", 999}, false)

	helpfulPrompt := prompt
	helpfulPrompt.SetHelp("This is helpful")

	tests := []struct {
		title    string
		prompt   *Select
		data     SelectTemplateData
		expected string
	}{
		{
			"Test Select question output",
			prompt,
			SelectTemplateData{SelectedIndex: 2, PageEntries: prompt.Options},
			`? Pick your word:  [Use arrows to move, type to filter]
  foo
  bar
❯ baz
  buz
`,
		},
		{
			"Test Select answer output",
			prompt,
			SelectTemplateData{Answer: prompt.Options[3], ShowAnswer: true, PageEntries: prompt.Options},
			"? Pick your word: buz\n",
		},
		{
			"Test Select question output with help hidden",
			helpfulPrompt,
			SelectTemplateData{SelectedIndex: 2, PageEntries: prompt.Options},
			`? Pick your word:  [Use arrows to move, type to filter, ? for more help]
  foo
  bar
❯ baz
  buz
`,
		},
		{
			"Test Select question output with help shown",
			helpfulPrompt,
			SelectTemplateData{SelectedIndex: 2, ShowHelp: true, PageEntries: prompt.Options},
			`ⓘ This is helpful
? Pick your word:  [Use arrows to move, type to filter]
  foo
  bar
❯ baz
  buz
`,
		},
	}

	outputBuffer := bytes.NewBufferString("")
	terminal.Stdout = outputBuffer

	for _, test := range tests {
		outputBuffer.Reset()
		test.data.Select = *test.prompt
		err := test.prompt.Render(
			SelectQuestionTemplate,
			test.data,
		)
		assert.Nil(t, err, test.title)
		assert.Equal(t, test.expected, outputBuffer.String(), test.title)
	}
}

func TestSelectionPagination_tooFew(t *testing.T) {
	prompt := NewSingleSelect()
	prompt.SetMessage("Pick your word:").
		AddOption("choice1", nil, false).
		AddOption("choice2", nil, false).
		AddOption("choice3", nil, false).
		SetPageSize(4)
	// the current selection
	prompt.selectedIndex = 3

	// compute the page info
	page, idx := prompt.Paginate(prompt.Options)

	// make sure we see the full list of options
	assert.Equal(t, prompt.Options, page)
	// with the second index highlighted (no change)
	assert.Equal(t, 3, idx)
}

func TestSelectionPagination_firstHalf(t *testing.T) {
	prompt := NewSingleSelect()
	prompt.SetMessage("Pick your word:").
		AddOption("choice1", nil, false).
		AddOption("choice2", nil, false).
		AddOption("choice3", nil, false).
		AddOption("choice4", nil, false).
		AddOption("choice5", nil, false).
		AddOption("choice6", nil, false).
		SetPageSize(4)
	// the current selection
	prompt.selectedIndex = 2

	// compute the page info
	page, idx := prompt.Paginate(prompt.Options)

	// we should see the first three options
	assert.Equal(t, prompt.Options[0:4], page)
	// with the second index highlighted
	assert.Equal(t, 2, idx)
}

func TestSelectionPagination_middle(t *testing.T) {
	prompt := NewSingleSelect()
	prompt.SetMessage("Pick your word:").
		AddOption("choice1", nil, false).
		AddOption("choice2", nil, false).
		AddOption("choice3", nil, false).
		AddOption("choice4", nil, false).
		AddOption("choice5", nil, false).
		SetPageSize(2)
	// the current selection
	prompt.selectedIndex = 3

	// compute the page info
	page, idx := prompt.Paginate(prompt.Options)

	// we should see the first three options
	assert.Equal(t, prompt.Options[2:4], page)
	// with the second index highlighted
	assert.Equal(t, 1, idx)
}

func TestSelectionPagination_lastHalf(t *testing.T) {
	prompt := NewSingleSelect()
	prompt.SetMessage("Pick your word:").
		AddOption("choice1", nil, false).
		AddOption("choice2", nil, false).
		AddOption("choice3", nil, false).
		AddOption("choice4", nil, false).
		AddOption("choice5", nil, false).
		SetPageSize(3)
	// the current selection
	prompt.selectedIndex = 4

	// compute the page info
	page, idx := prompt.Paginate(prompt.Options)

	// we should see the last three options
	assert.Equal(t, prompt.Options[2:5], page)
	// we should be at the bottom of the list
	assert.Equal(t, 2, idx)
}

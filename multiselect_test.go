package survey

import (
	"testing"

	"github.com/AlecAivazis/survey/core"
	"github.com/stretchr/testify/assert"
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

	tests := []struct {
		title    string
		prompt   MultiSelect
		template string
		data     MultiSelectTemplateData
		expected string
	}{
		{
			"Test MultiSelect options output",
			prompt,
			MultiSelectOptionsTemplate,
			MultiSelectTemplateData{SelectedIndex: 2, Checked: map[int]bool{1: true, 3: true}},
			`  ◯  foo
  ◉  bar
❯ ◯  baz
  ◉  buz
`,
		},
		{
			"Test MultiSelect question output",
			prompt,
			MultiSelectQuestionTemplate,
			MultiSelectTemplateData{},
			"? Pick your words:",
		},
		{
			"Test MultiSelect answer output",
			prompt,
			MultiSelectQuestionTemplate,
			MultiSelectTemplateData{Answer: "foo, buz"},
			"? Pick your words: foo, buz",
		},
	}

	for _, test := range tests {
		test.data.MultiSelect = test.prompt
		actual, err := core.RunTemplate(
			test.template,
			test.data,
		)
		assert.Nil(t, err, test.title)
		assert.Equal(t, test.expected, actual, test.title)
	}
}

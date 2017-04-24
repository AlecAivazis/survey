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
		prompt   MultiSelect
		template string
		data     MultiSelectTemplateData
		expected string
	}{
		{
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
			prompt,
			MultiSelectQuestionTemplate,
			MultiSelectTemplateData{},
			"? Pick your words:",
		},
		{
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
		assert.Nil(t, err)
		assert.Equal(t, test.expected, actual)
	}
}

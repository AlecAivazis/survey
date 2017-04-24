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

func TestSelectRender(t *testing.T) {

	prompt := Select{
		Message: "Pick your word:",
		Options: []string{"foo", "bar", "baz", "buz"},
		Default: "baz",
	}

	tests := []struct {
		title    string
		prompt   Select
		template string
		data     SelectTemplateData
		expected string
	}{
		{
			"Test Select choices output",
			prompt,
			SelectChoicesTemplate,
			SelectTemplateData{SelectedIndex: 2},
			`  foo
  bar
> baz
  buz
`,
		},
		{
			"Test Select question output",
			prompt,
			SelectQuestionTemplate,
			SelectTemplateData{},
			"? Pick your word:",
		},
		{
			"Test Select answer output",
			prompt,
			SelectQuestionTemplate,
			SelectTemplateData{Answer: "buz"},
			"? Pick your word: buz",
		},
	}

	for _, test := range tests {
		test.data.Select = test.prompt
		actual, err := core.RunTemplate(
			test.template,
			test.data,
		)
		assert.Nil(t, err, test.title)
		assert.Equal(t, test.expected, actual, test.title)
	}
}

package survey

import (
	"bytes"
	"testing"

	"github.com/AlecAivazis/survey/core"
	"github.com/AlecAivazis/survey/terminal"
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
		data     MultiSelectTemplateData
		expected string
	}{
		{
			prompt,
			MultiSelectTemplateData{SelectedIndex: 2, Checked: map[int]bool{1: true, 3: true}},
			`? Pick your words:
  ◯  foo
  ◉  bar
❯ ◯  baz
  ◉  buz
`,
		},
		{
			prompt,
			MultiSelectTemplateData{Answer: "foo, buz"},
			"? Pick your words: foo, buz\n",
		},
	}

	outputBuffer := bytes.NewBufferString("")
	terminal.AnsiStdout = outputBuffer

	for _, test := range tests {
		outputBuffer.Reset()
		test.data.MultiSelect = test.prompt
		err := test.prompt.render(
			MultiSelectQuestionTemplate,
			test.data,
		)
		assert.Nil(t, err)
		assert.Equal(t, test.expected, outputBuffer.String())
	}
}

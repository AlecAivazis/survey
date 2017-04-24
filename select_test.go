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

func TestSelectRender(t *testing.T) {

	prompt := Select{
		Message: "Pick your word:",
		Options: []string{"foo", "bar", "baz", "buz"},
		Default: "baz",
	}

	tests := []struct {
		prompt   Select
		data     SelectTemplateData
		expected string
	}{
		{
			prompt,
			SelectTemplateData{SelectedIndex: 2},
			`? Pick your word:
  foo
  bar
> baz
  buz
`,
		},
		{
			prompt,
			SelectTemplateData{Answer: "buz"},
			"? Pick your word: buz\n",
		},
	}

	outputBuffer := bytes.NewBufferString("")
	terminal.AnsiStdout = outputBuffer

	for _, test := range tests {
		outputBuffer.Reset()
		test.data.Select = test.prompt
		err := test.prompt.render(
			SelectQuestionTemplate,
			test.data,
		)
		assert.Nil(t, err)
		assert.Equal(t, test.expected, outputBuffer.String())
	}
}

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

func TestInputRender(t *testing.T) {

	prompt := Input{
		Message: "What is your favorite month:",
		Default: "April",
	}

	tests := []struct {
		prompt   Input
		data     InputTemplateData
		expected string
	}{
		{
			prompt,
			InputTemplateData{},
			"? What is your favorite month: (April) ",
		},
		{
			prompt,
			InputTemplateData{Answer: "October"},
			"? What is your favorite month: October\n",
		},
	}

	outputBuffer := bytes.NewBufferString("")
	terminal.AnsiStdout = outputBuffer

	for _, test := range tests {
		outputBuffer.Reset()
		test.data.Input = test.prompt
		err := test.prompt.render(
			InputQuestionTemplate,
			test.data,
		)
		assert.Nil(t, err)
		assert.Equal(t, test.expected, outputBuffer.String())
	}
}

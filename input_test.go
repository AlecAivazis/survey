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

	tests := []struct {
		title    string
		prompt   Input
		data     InputTemplateData
		expected string
	}{
		{
			"Test Input question output without default",
			Input{Message: "What is your favorite month:"},
			InputTemplateData{},
			"? What is your favorite month: ",
		},
		{
			"Test Input question output with default",
			Input{Message: "What is your favorite month:", Default: "April"},
			InputTemplateData{},
			"? What is your favorite month: (April) ",
		},
		{
			"Test Input answer output",
			Input{Message: "What is your favorite month:"},
			InputTemplateData{Answer: "October"},
			"? What is your favorite month: October\n",
		},
	}

	outputBuffer := bytes.NewBufferString("")
	terminal.Stdout = outputBuffer

	for _, test := range tests {
		outputBuffer.Reset()
		test.data.Input = test.prompt
		err := test.prompt.render(
			InputQuestionTemplate,
			test.data,
		)
		assert.Nil(t, err, test.title)
		assert.Equal(t, test.expected, outputBuffer.String(), test.title)
	}
}

package survey

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/AlecAivazis/survey/core"
	"github.com/AlecAivazis/survey/terminal"
	"github.com/stretchr/testify/assert"
)

func init() {
	// disable color output for all prompts to simplify testing
	core.DisableColor = true
}

func TestConfirmRender(t *testing.T) {

	testError := fmt.Errorf("TEST")
	tests := []struct {
		prompt   Confirm
		data     ConfirmTemplateData
		expected string
	}{
		{
			Confirm{Message: "Is pizza your favorite food?", Default: true},
			ConfirmTemplateData{},
			`? Is pizza your favorite food? (Y/n) `,
		},
		{
			Confirm{Message: "Is pizza your favorite food?", Default: false},
			ConfirmTemplateData{},
			`? Is pizza your favorite food? (y/N) `,
		},
		{
			Confirm{Message: "Is pizza your favorite food?"},
			ConfirmTemplateData{Answer: "Yes"},
			"? Is pizza your favorite food? Yes\n",
		},
		{
			Confirm{Message: "Is pizza your favorite food?"},
			ConfirmTemplateData{Answer: "Yes", Error: &testError},
			`✘ Sorry, your reply was invalid: TEST
? Is pizza your favorite food? Yes
`,
		},
	}

	outputBuffer := bytes.NewBufferString("")
	terminal.AnsiStdout = outputBuffer

	for _, test := range tests {
		outputBuffer.Reset()
		test.data.Confirm = test.prompt
		err := test.prompt.render(
			ConfirmQuestionTemplate,
			test.data,
		)
		assert.Nil(t, err)
		assert.Equal(t, test.expected, outputBuffer.String())
	}
}

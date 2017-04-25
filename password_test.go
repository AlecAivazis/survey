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

func TestPasswordRender(t *testing.T) {

	prompt := Password{
		Message: "Tell me your secret:",
	}

	tests := []struct {
		title    string
		prompt   Password
		expected string
	}{
		{
			"Test Password question output",
			prompt,
			"? Tell me your secret: ",
		},
	}

	outputBuffer := bytes.NewBufferString("")
	terminal.Stdout = outputBuffer

	for _, test := range tests {
		outputBuffer.Reset()
		err := test.prompt.render(
			PasswordQuestionTemplate,
			&test.prompt,
		)
		assert.Nil(t, err, test.title)
		assert.Equal(t, test.expected, outputBuffer.String(), test.title)
	}
}

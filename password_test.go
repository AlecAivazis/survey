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
		prompt   Password
		expected string
	}{
		{
			prompt,
			"? Tell me your secret: ",
		},
	}

	outputBuffer := bytes.NewBufferString("")
	terminal.AnsiStdout = outputBuffer

	for _, test := range tests {
		outputBuffer.Reset()
		err := test.prompt.render(
			PasswordQuestionTemplate,
			&test.prompt,
		)
		assert.Nil(t, err)
		assert.Equal(t, test.expected, outputBuffer.String())
	}
}

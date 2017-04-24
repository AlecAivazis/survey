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

	for _, test := range tests {
		actual, err := core.RunTemplate(
			PasswordQuestionTemplate,
			&test.prompt,
		)
		assert.Nil(t, err)
		assert.Equal(t, test.expected, actual)
	}
}

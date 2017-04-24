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
			"? What is your favorite month: October",
		},
	}

	for _, test := range tests {
		test.data.Input = test.prompt
		actual, err := core.RunTemplate(
			InputQuestionTemplate,
			test.data,
		)
		assert.Nil(t, err)
		assert.Equal(t, test.expected, actual)
	}
}

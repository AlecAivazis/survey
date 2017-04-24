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

func TestConfirmRender(t *testing.T) {
	tests := []struct {
		title    string
		prompt   Confirm
		data     ConfirmTemplateData
		expected string
	}{
		{
			"Test Connfirm question output with default true",
			Confirm{Message: "Is pizza your favorite food?", Default: true},
			ConfirmTemplateData{},
			`? Is pizza your favorite food? (Y/n) `,
		},
		{
			"Test Confirm question output with default false",
			Confirm{Message: "Is pizza your favorite food?", Default: false},
			ConfirmTemplateData{},
			`? Is pizza your favorite food? (y/N) `,
		},
		{
			"Test Confirm answer output",
			Confirm{Message: "Is pizza your favorite food?"},
			ConfirmTemplateData{Answer: "Yes"},
			"? Is pizza your favorite food? Yes",
		},
	}

	for _, test := range tests {
		test.data.Confirm = test.prompt
		actual, err := core.RunTemplate(
			ConfirmQuestionTemplate,
			test.data,
		)
		assert.Nil(t, err, test.title)
		assert.Equal(t, test.expected, actual, test.title)
	}
}

package survey

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/Netflix/go-expect"
	"github.com/stretchr/testify/assert"
	"gopkg.in/AlecAivazis/survey.v1/core"
	"gopkg.in/AlecAivazis/survey.v1/terminal"
)

func init() {
	// disable color output for all prompts to simplify testing
	core.DisableColor = true
}

func TestTextRender(t *testing.T) {

	tests := []struct {
		title    string
		prompt   Text
		data     TextTemplateData
		expected string
	}{
		{
			"Test info message output",
			Text{Message: "This is a info message!"},
			TextTemplateData{},
			"",
		},
		{
			"Test info message output",
			Text{Message: "This is a info message!", Level: Info},
			TextTemplateData{},
			"",
		},
		{
			"Test info message output",
			Text{Message: "This is a warning message!", Level: Warning},
			TextTemplateData{},
			"",
		},
		{
			"Test info message output",
			Text{Message: "This is a danger message!", Level: Danger},
			TextTemplateData{},
			"",
		},
	}

	for _, test := range tests {
		r, w, err := os.Pipe()
		assert.Nil(t, err, test.title)

		test.prompt.WithStdio(terminal.Stdio{Out: w})
		test.data.Text = test.prompt
		err = test.prompt.Render(
			InfoTemplate,
			test.data,
		)
		assert.Nil(t, err, test.title)

		w.Close()
		var buf bytes.Buffer
		io.Copy(&buf, r)

		assert.Contains(t, buf.String(), test.expected, test.title)
	}
}

func TestTextPrompt(t *testing.T) {
	tests := []PromptTest{
		{
			"Test info message interaction",
			&Text{
				Message: "This is a info message!",
				Level:   Info,
			},
			func(c *expect.Console) {
				c.ExpectString("This is a info message!")
				c.ExpectEOF()
			},
			"",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			RunPromptTest(t, test)
		})
	}
}

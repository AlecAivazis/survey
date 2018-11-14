package survey

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"

	expect "github.com/Netflix/go-expect"
	"github.com/stretchr/testify/assert"
	"gopkg.in/AlecAivazis/survey.v1/core"
	"gopkg.in/AlecAivazis/survey.v1/terminal"
)

func init() {
	// disable color output for all prompts to simplify testing
	core.DisableColor = true
}

func TestMultilineInputRender(t *testing.T) {

	tests := []struct {
		title    string
		prompt   MultilineInput
		data     MultilineInputTemplateData
		expected string
	}{
		{
			"Test MultilineInput question output without default",
			MultilineInput{Message: "What is your favorite month:"},
			MultilineInputTemplateData{},
			fmt.Sprintf("%s What is your favorite month: [Enter 2 empty lines to finish]", core.QuestionIcon),
		},
		{
			"Test MultilineInput question output with default",
			MultilineInput{Message: "What is your favorite month:", Default: "April"},
			MultilineInputTemplateData{},
			fmt.Sprintf("%s What is your favorite month: (April) [Enter 2 empty lines to finish]", core.QuestionIcon),
		},
		{
			"Test MultilineInput answer output",
			MultilineInput{Message: "What is your favorite month:"},
			MultilineInputTemplateData{Answer: "October", ShowAnswer: true},
			fmt.Sprintf("%s What is your favorite month: October\n", core.QuestionIcon),
		},
		{
			"Test MultilineInput question output without default but with help hidden",
			MultilineInput{Message: "What is your favorite month:", Help: "This is helpful"},
			MultilineInputTemplateData{},
			fmt.Sprintf("%s What is your favorite month: [Enter 2 empty lines to finish]", string(core.HelpInputRune)),
		},
		{
			"Test MultilineInput question output with default and with help hidden",
			MultilineInput{Message: "What is your favorite month:", Default: "April", Help: "This is helpful"},
			MultilineInputTemplateData{},
			fmt.Sprintf("%s What is your favorite month: (April) [Enter 2 empty lines to finish]", string(core.HelpInputRune)),
		},
		{
			"Test MultilineInput question output without default but with help shown",
			MultilineInput{Message: "What is your favorite month:", Help: "This is helpful"},
			MultilineInputTemplateData{ShowHelp: true},
			fmt.Sprintf("%s This is helpful\n%s What is your favorite month: [Enter 2 empty lines to finish]", core.HelpIcon, core.QuestionIcon),
		},
		{
			"Test MultilineInput question output with default and with help shown",
			MultilineInput{Message: "What is your favorite month:", Default: "April", Help: "This is helpful"},
			MultilineInputTemplateData{ShowHelp: true},
			fmt.Sprintf("%s This is helpful\n%s What is your favorite month: (April) [Enter 2 empty lines to finish]", core.HelpIcon, core.QuestionIcon),
		},
	}

	for _, test := range tests {
		r, w, err := os.Pipe()
		assert.Nil(t, err, test.title)

		test.prompt.WithStdio(terminal.Stdio{Out: w})
		test.data.MultilineInput = test.prompt
		err = test.prompt.Render(
			MultilineInputQuestionTemplate,
			test.data,
		)
		assert.Nil(t, err, test.title)

		w.Close()
		var buf bytes.Buffer
		io.Copy(&buf, r)

		assert.Contains(t, buf.String(), test.expected, test.title)
	}
}

func TestMultilineInputPrompt(t *testing.T) {
	tests := []PromptTest{
		{
			"Test MultilineInput prompt interaction",
			&MultilineInput{
				Message: "What is your name?",
			},
			func(c *expect.Console) {
				c.ExpectString("What is your name?")
				c.SendLine("Larry Bird\nI guess...\nnot sure\n\n")
				c.ExpectEOF()
			},
			"Larry Bird\nI guess...\nnot sure",
		},
		{
			"Test MultilineInput prompt interaction with default",
			&MultilineInput{
				Message: "What is your name?",
				Default: "Johnny Appleseed",
			},
			func(c *expect.Console) {
				c.ExpectString("What is your name?")
				c.SendLine("\n\n")
				c.ExpectEOF()
			},
			"Johnny Appleseed",
		},
		{
			"Test MultilineInput prompt interaction overriding default",
			&MultilineInput{
				Message: "What is your name?",
				Default: "Johnny Appleseed",
			},
			func(c *expect.Console) {
				c.ExpectString("What is your name?")
				c.SendLine("Larry Bird\n\n")
				c.ExpectEOF()
			},
			"Larry Bird",
		},
		{
			"Test MultilineInput does not implement help interaction",
			&MultilineInput{
				Message: "What is your name?",
				Help:    "It might be Satoshi Nakamoto",
			},
			func(c *expect.Console) {
				c.ExpectString("What is your name?")
				c.SendLine("?")
				c.SendLine("Satoshi Nakamoto\n\n")
				c.ExpectEOF()
			},
			"?\nSatoshi Nakamoto",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			RunPromptTest(t, test)
		})
	}
}

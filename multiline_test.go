package survey

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/AlecAivazis/survey/v2/core"
	"github.com/AlecAivazis/survey/v2/terminal"
	expect "github.com/Netflix/go-expect"
	"github.com/stretchr/testify/assert"
)

func init() {
	// disable color output for all prompts to simplify testing
	core.DisableColor = true
}

func TestMultilineRender(t *testing.T) {

	tests := []struct {
		title    string
		prompt   Multiline
		data     MultilineTemplateData
		expected string
	}{
		{
			"Test Multiline question output without default",
			Multiline{Message: "What is your favorite month:"},
			MultilineTemplateData{Icons: &defaultIconSet},
			fmt.Sprintf("%s What is your favorite month: [Enter 2 empty lines to finish]", defaultIconSet.Question),
		},
		{
			"Test Multiline question output with default",
			Multiline{Message: "What is your favorite month:", Default: "April"},
			MultilineTemplateData{Icons: &defaultIconSet},
			fmt.Sprintf("%s What is your favorite month: (April) [Enter 2 empty lines to finish]", defaultIconSet.Question),
		},
		{
			"Test Multiline answer output",
			Multiline{Message: "What is your favorite month:"},
			MultilineTemplateData{Answer: "October", ShowAnswer: true, Icons: &defaultIconSet},
			fmt.Sprintf("%s What is your favorite month: \nOctober", defaultIconSet.Question),
		},
		{
			"Test Multiline question output without default but with help hidden",
			Multiline{Message: "What is your favorite month:", Help: "This is helpful"},
			MultilineTemplateData{Icons: &defaultIconSet},
			fmt.Sprintf("%s What is your favorite month: [Enter 2 empty lines to finish]", string(defaultIconSet.HelpInput)),
		},
		{
			"Test Multiline question output with default and with help hidden",
			Multiline{Message: "What is your favorite month:", Default: "April", Help: "This is helpful"},
			MultilineTemplateData{Icons: &defaultIconSet},
			fmt.Sprintf("%s What is your favorite month: (April) [Enter 2 empty lines to finish]", string(defaultIconSet.HelpInput)),
		},
		{
			"Test Multiline question output without default but with help shown",
			Multiline{Message: "What is your favorite month:", Help: "This is helpful"},
			MultilineTemplateData{ShowHelp: true, Icons: &defaultIconSet},
			fmt.Sprintf("%s This is helpful\n%s What is your favorite month: [Enter 2 empty lines to finish]", defaultIconSet.Help, defaultIconSet.Question),
		},
		{
			"Test Multiline question output with default and with help shown",
			Multiline{Message: "What is your favorite month:", Default: "April", Help: "This is helpful"},
			MultilineTemplateData{ShowHelp: true, Icons: &defaultIconSet},
			fmt.Sprintf("%s This is helpful\n%s What is your favorite month: (April) [Enter 2 empty lines to finish]", defaultIconSet.Help, defaultIconSet.Question),
		},
	}

	for _, test := range tests {
		r, w, err := os.Pipe()
		assert.Nil(t, err, test.title)

		test.prompt.WithStdio(terminal.Stdio{Out: w})
		test.data.Multiline = test.prompt
		err = test.prompt.Render(
			MultilineQuestionTemplate,
			test.data,
		)
		assert.Nil(t, err, test.title)

		w.Close()
		var buf bytes.Buffer
		io.Copy(&buf, r)

		assert.Contains(t, buf.String(), test.expected, test.title)
	}
}

func TestMultilinePrompt(t *testing.T) {
	tests := []PromptTest{
		{
			"Test Multiline prompt interaction",
			&Multiline{
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
			"Test Multiline prompt interaction with default",
			&Multiline{
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
			"Test Multiline prompt interaction overriding default",
			&Multiline{
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
			"Test Multiline does not implement help interaction",
			&Multiline{
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

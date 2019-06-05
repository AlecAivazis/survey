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

func TestMultiInputRender(t *testing.T) {

	tests := []struct {
		title    string
		prompt   MultiInput
		data     MultiInputTemplateData
		expected string
	}{
		{
			"Test Input question output without default",
			MultiInput{Message: "What is your favorite month:"},
			MultiInputTemplateData{},
			fmt.Sprintf("%s What is your favorite month: ", core.QuestionIcon),
		},
		{
			"Test Input question output with default",
			MultiInput{Message: "What is your favorite month:", Default: []string{"April"}},
			MultiInputTemplateData{},
			fmt.Sprintf("%s What is your favorite month: ([April]) ", core.QuestionIcon),
		},
		{
			"Test Input answer output",
			MultiInput{Message: "What is your favorite month:"},
			MultiInputTemplateData{Answer: []string{"October"}, ShowAnswer: true},
			fmt.Sprintf("%s What is your favorite month: [October]\n", core.QuestionIcon),
		},
		{
			"Test Input question output without default but with help hidden",
			MultiInput{Message: "What is your favorite month:", Help: "This is helpful"},
			MultiInputTemplateData{},
			fmt.Sprintf("%s What is your favorite month: [%s for help] ", core.QuestionIcon, string(core.HelpInputRune)),
		},
		{
			"Test Input question output with default and with help hidden",
			MultiInput{Message: "What is your favorite month:", Default: []string{"April"}, Help: "This is helpful"},
			MultiInputTemplateData{},
			fmt.Sprintf("%s What is your favorite month: [%s for help] ([April]) ", core.QuestionIcon, string(core.HelpInputRune)),
		},
		{
			"Test Input question output without default but with help shown",
			MultiInput{Message: "What is your favorite month:", Help: "This is helpful"},
			MultiInputTemplateData{ShowHelp: true},
			fmt.Sprintf("%s This is helpful\n%s What is your favorite month: ", core.HelpIcon, core.QuestionIcon),
		},
		{
			"Test Input question output with default and with help shown",
			MultiInput{Message: "What is your favorite month:", Default: []string{"April"}, Help: "This is helpful"},
			MultiInputTemplateData{ShowHelp: true},
			fmt.Sprintf("%s This is helpful\n%s What is your favorite month: ([April]) ", core.HelpIcon, core.QuestionIcon),
		},
	}

	for _, test := range tests {
		r, w, err := os.Pipe()
		assert.Nil(t, err, test.title)

		test.prompt.WithStdio(terminal.Stdio{Out: w})
		test.data.MultiInput = test.prompt
		err = test.prompt.Render(
			InputQuestionTemplate,
			test.data,
		)
		assert.Nil(t, err, test.title)

		w.Close()
		var buf bytes.Buffer
		io.Copy(&buf, r)

		assert.Contains(t, buf.String(), test.expected, test.title)
	}
}

func TestMultiInputPrompt(t *testing.T) {
	tests := []PromptTest{
		{
			"Test MultiInput prompt interaction",
			&MultiInput{
				Message: "Enter the names of your friends:",
			},
			func(c *expect.Console) {
				c.ExpectString("Enter the names of your friends:")
				c.ExpectString("\n")
				c.ExpectString("#1:")
				c.SendLine("Larry Bird")
				c.SendLine("")
				c.ExpectEOF()
			},
			[]string{"Larry Bird"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			RunPromptTest(t, test)
		})
	}
}

package survey

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/AlecAivazis/survey/v2/core"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/stretchr/testify/assert"
)

func init() {
	// disable color output for all prompts to simplify testing
	core.DisableColor = true
}

func TestEditorRender(t *testing.T) {
	tests := []struct {
		title    string
		prompt   Editor
		data     EditorTemplateData
		expected string
	}{
		{
			"Test Editor question output without default",
			Editor{Message: "What is your favorite month:"},
			EditorTemplateData{},
			fmt.Sprintf("%s What is your favorite month: [Enter to launch editor] ", defaultIcons().Question.Text),
		},
		{
			"Test Editor question output with default",
			Editor{Message: "What is your favorite month:", Default: "April"},
			EditorTemplateData{},
			fmt.Sprintf("%s What is your favorite month: (April) [Enter to launch editor] ", defaultIcons().Question.Text),
		},
		{
			"Test Editor question output with HideDefault",
			Editor{Message: "What is your favorite month:", Default: "April", HideDefault: true},
			EditorTemplateData{},
			fmt.Sprintf("%s What is your favorite month: [Enter to launch editor] ", defaultIcons().Question.Text),
		},
		{
			"Test Editor answer output",
			Editor{Message: "What is your favorite month:"},
			EditorTemplateData{Answer: "October", ShowAnswer: true},
			fmt.Sprintf("%s What is your favorite month: October\n", defaultIcons().Question.Text),
		},
		{
			"Test Editor question output without default but with help hidden",
			Editor{Message: "What is your favorite month:", Help: "This is helpful"},
			EditorTemplateData{},
			fmt.Sprintf("%s What is your favorite month: [%s for help] [Enter to launch editor] ", defaultIcons().Question.Text, string(defaultPromptConfig().HelpInput)),
		},
		{
			"Test Editor question output with default and with help hidden",
			Editor{Message: "What is your favorite month:", Default: "April", Help: "This is helpful"},
			EditorTemplateData{},
			fmt.Sprintf("%s What is your favorite month: [%s for help] (April) [Enter to launch editor] ", defaultIcons().Question.Text, string(defaultPromptConfig().HelpInput)),
		},
		{
			"Test Editor question output without default but with help shown",
			Editor{Message: "What is your favorite month:", Help: "This is helpful"},
			EditorTemplateData{ShowHelp: true},
			fmt.Sprintf("%s This is helpful\n%s What is your favorite month: [Enter to launch editor] ", defaultIcons().Help.Text, defaultIcons().Question.Text),
		},
		{
			"Test Editor question output with default and with help shown",
			Editor{Message: "What is your favorite month:", Default: "April", Help: "This is helpful"},
			EditorTemplateData{ShowHelp: true},
			fmt.Sprintf("%s This is helpful\n%s What is your favorite month: (April) [Enter to launch editor] ", defaultIcons().Help.Text, defaultIcons().Question.Text),
		},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			r, w, err := os.Pipe()
			assert.NoError(t, err)

			test.prompt.WithStdio(terminal.Stdio{Out: w})
			test.data.Editor = test.prompt

			// set the icon set
			test.data.Config = defaultPromptConfig()

			err = test.prompt.Render(
				EditorQuestionTemplate,
				test.data,
			)
			assert.NoError(t, err)

			assert.NoError(t, w.Close())
			var buf bytes.Buffer
			_, err = io.Copy(&buf, r)
			assert.NoError(t, err)

			assert.Contains(t, buf.String(), test.expected)
		})
	}
}

func TestEditorPrompt(t *testing.T) {
	if _, err := exec.LookPath("vi"); err != nil {
		t.Skip("warning: vi not found in PATH")
	}

	tests := []PromptTest{
		{
			"Test Editor prompt interaction",
			&Editor{
				Editor:  "vi",
				Message: "Edit git commit message",
			},
			func(c expectConsole) {
				c.ExpectString("Edit git commit message [Enter to launch editor]")
				c.SendLine("")
				time.Sleep(time.Millisecond)
				c.Send("ccAdd editor prompt tests\x1b")
				c.SendLine(":wq!")
				c.ExpectEOF()
			},
			"Add editor prompt tests\n",
		},
		{
			"Test Editor prompt interaction with default",
			&Editor{
				Editor:  "vi",
				Message: "Edit git commit message",
				Default: "No comment",
			},
			func(c expectConsole) {
				c.ExpectString("Edit git commit message (No comment) [Enter to launch editor]")
				c.SendLine("")
				time.Sleep(time.Millisecond)
				c.SendLine(":q!")
				c.ExpectEOF()
			},
			"No comment",
		},
		{
			"Test Editor prompt interaction overriding default",
			&Editor{
				Editor:  "vi",
				Message: "Edit git commit message",
				Default: "No comment",
			},
			func(c expectConsole) {
				c.ExpectString("Edit git commit message (No comment) [Enter to launch editor]")
				c.SendLine("")
				time.Sleep(time.Millisecond)
				c.Send("ccAdd editor prompt tests\x1b")
				c.SendLine(":wq!")
				c.ExpectEOF()
			},
			"Add editor prompt tests\n",
		},
		{
			"Test Editor prompt interaction hiding default",
			&Editor{
				Editor:      "vi",
				Message:     "Edit git commit message",
				Default:     "No comment",
				HideDefault: true,
			},
			func(c expectConsole) {
				c.ExpectString("Edit git commit message [Enter to launch editor]")
				c.SendLine("")
				time.Sleep(time.Millisecond)
				c.SendLine(":q!")
				c.ExpectEOF()
			},
			"No comment",
		},
		{
			"Test Editor prompt interaction and prompt for help",
			&Editor{
				Editor:  "vi",
				Message: "Edit git commit message",
				Help:    "Describe your git commit",
			},
			func(c expectConsole) {
				c.ExpectString(
					fmt.Sprintf(
						"Edit git commit message [%s for help] [Enter to launch editor]",
						string(defaultPromptConfig().HelpInput),
					),
				)
				c.SendLine("?")
				c.ExpectString("Describe your git commit")
				c.SendLine("")
				time.Sleep(time.Millisecond)
				c.Send("ccAdd editor prompt tests\x1b")
				c.SendLine(":wq!")
				c.ExpectEOF()
			},
			"Add editor prompt tests\n",
		},
		{
			"Test Editor prompt interaction with default and append default",
			&Editor{
				Editor:        "vi",
				Message:       "Edit git commit message",
				Default:       "No comment",
				AppendDefault: true,
			},
			func(c expectConsole) {
				c.ExpectString("Edit git commit message (No comment) [Enter to launch editor]")
				c.SendLine("")
				c.ExpectString("No comment")
				c.SendLine("dd")
				c.SendLine(":wq!")
				c.ExpectEOF()
			},
			"",
		},
		{
			"Test Editor prompt interaction with editor args",
			&Editor{
				Editor:  "vi --",
				Message: "Edit git commit message",
			},
			func(c expectConsole) {
				c.ExpectString("Edit git commit message [Enter to launch editor]")
				c.SendLine("")
				time.Sleep(time.Millisecond)
				c.Send("ccAdd editor prompt tests\x1b")
				c.SendLine(":wq!")
				c.ExpectEOF()
			},
			"Add editor prompt tests\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			RunPromptTest(t, test)
		})
	}
}

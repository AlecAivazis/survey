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

func TestInputRender(t *testing.T) {

	suggestFn := func(string) (s []string) { return s }

	tests := []struct {
		title    string
		prompt   Input
		data     InputTemplateData
		expected string
	}{
		{
			"Test Input question output without default",
			Input{Message: "What is your favorite month:"},
			InputTemplateData{},
			fmt.Sprintf("%s What is your favorite month: ", defaultIcons().Question.Text),
		},
		{
			"Test Input question output with default",
			Input{Message: "What is your favorite month:", Default: "April"},
			InputTemplateData{},
			fmt.Sprintf("%s What is your favorite month: (April) ", defaultIcons().Question.Text),
		},
		{
			"Test Input answer output",
			Input{Message: "What is your favorite month:"},
			InputTemplateData{ShowAnswer: true, Answer: "October"},
			fmt.Sprintf("%s What is your favorite month: October\n", defaultIcons().Question.Text),
		},
		{
			"Test Input question output without default but with help hidden",
			Input{Message: "What is your favorite month:", Help: "This is helpful"},
			InputTemplateData{},
			fmt.Sprintf("%s What is your favorite month: [%s for help] ", defaultIcons().Question.Text, string(defaultPromptConfig().HelpInput)),
		},
		{
			"Test Input question output with default and with help hidden",
			Input{Message: "What is your favorite month:", Default: "April", Help: "This is helpful"},
			InputTemplateData{},
			fmt.Sprintf("%s What is your favorite month: [%s for help] (April) ", defaultIcons().Question.Text, string(defaultPromptConfig().HelpInput)),
		},
		{
			"Test Input question output without default but with help shown",
			Input{Message: "What is your favorite month:", Help: "This is helpful"},
			InputTemplateData{ShowHelp: true},
			fmt.Sprintf("%s This is helpful\n%s What is your favorite month: ", defaultIcons().Help.Text, defaultIcons().Question.Text),
		},
		{
			"Test Input question output with default and with help shown",
			Input{Message: "What is your favorite month:", Default: "April", Help: "This is helpful"},
			InputTemplateData{ShowHelp: true},
			fmt.Sprintf("%s This is helpful\n%s What is your favorite month: (April) ", defaultIcons().Help.Text, defaultIcons().Question.Text),
		},
		{
			"Test Input question output with completion",
			Input{Message: "What is your favorite month:", Suggest: suggestFn},
			InputTemplateData{},
			fmt.Sprintf("%s What is your favorite month: [%s for suggestions] ", defaultIcons().Question.Text, string(defaultPromptConfig().SuggestInput)),
		},
		{
			"Test Input question output with suggestions and help hidden",
			Input{Message: "What is your favorite month:", Suggest: suggestFn, Help: "This is helpful"},
			InputTemplateData{},
			fmt.Sprintf("%s What is your favorite month: [%s for help, %s for suggestions] ", defaultIcons().Question.Text, string(defaultPromptConfig().HelpInput), string(defaultPromptConfig().SuggestInput)),
		},
		{
			"Test Input question output with suggestions and default and help hidden",
			Input{Message: "What is your favorite month:", Suggest: suggestFn, Help: "This is helpful", Default: "April"},
			InputTemplateData{},
			fmt.Sprintf("%s What is your favorite month: [%s for help, %s for suggestions] (April) ", defaultIcons().Question.Text, string(defaultPromptConfig().HelpInput), string(defaultPromptConfig().SuggestInput)),
		},
		{
			"Test Input question output with suggestions shown",
			Input{Message: "What is your favorite month:", Suggest: suggestFn},
			InputTemplateData{
				Answer:        "February",
				PageEntries:   core.OptionAnswerList([]string{"January", "February", "March", "etc..."}),
				SelectedIndex: 1,
			},
			fmt.Sprintf(
				"%s What is your favorite month: February [Use arrows to move, enter to select, type to continue]\n"+
					"  January\n%s February\n  March\n  etc...\n",
				defaultIcons().Question.Text, defaultPromptConfig().Icons.SelectFocus.Text,
			),
		},
	}

	for _, test := range tests {
		r, w, err := os.Pipe()
		assert.Nil(t, err, test.title)

		test.prompt.WithStdio(terminal.Stdio{Out: w})
		test.data.Input = test.prompt

		// set the runtime config
		test.data.Config = defaultPromptConfig()

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

func TestInputPrompt(t *testing.T) {

	tests := []PromptTest{
		{
			"Test Input prompt interaction",
			&Input{
				Message: "What is your name?",
			},
			func(c *expect.Console) {
				c.ExpectString("What is your name?")
				c.SendLine("Larry Bird")
				c.ExpectEOF()
			},
			"Larry Bird",
		},
		{
			"Test Input prompt interaction with default",
			&Input{
				Message: "What is your name?",
				Default: "Johnny Appleseed",
			},
			func(c *expect.Console) {
				c.ExpectString("What is your name?")
				c.SendLine("")
				c.ExpectEOF()
			},
			"Johnny Appleseed",
		},
		{
			"Test Input prompt interaction overriding default",
			&Input{
				Message: "What is your name?",
				Default: "Johnny Appleseed",
			},
			func(c *expect.Console) {
				c.ExpectString("What is your name?")
				c.SendLine("Larry Bird")
				c.ExpectEOF()
			},
			"Larry Bird",
		},
		{
			"Test Input prompt interaction and prompt for help",
			&Input{
				Message: "What is your name?",
				Help:    "It might be Satoshi Nakamoto",
			},
			func(c *expect.Console) {
				c.ExpectString("What is your name?")
				c.SendLine("?")
				c.ExpectString("It might be Satoshi Nakamoto")
				c.SendLine("Satoshi Nakamoto")
				c.ExpectEOF()
			},
			"Satoshi Nakamoto",
		},
		{
			// https://en.wikipedia.org/wiki/ANSI_escape_code
			// Device Status Report - Reports the cursor position (CPR) to the
			// application as (as though typed at the keyboard) ESC[n;mR, where n is the
			// row and m is the column.
			"Test Input prompt with R matching DSR",
			&Input{
				Message: "What is your name?",
			},
			func(c *expect.Console) {
				c.ExpectString("What is your name?")
				c.SendLine("R")
				c.ExpectEOF()
			},
			"R",
		},
		{
			"Test Input prompt interaction when delete",
			&Input{
				Message: "What is your name?",
			},
			func(c *expect.Console) {
				c.ExpectString("What is your name?")
				c.Send("Johnny ")
				c.Send(string(terminal.KeyDelete))
				c.SendLine("")
				c.ExpectEOF()
			},
			"Johnny",
		},
		{
			"Test Input prompt interaction when delete rune",
			&Input{
				Message: "What is your name?",
			},
			func(c *expect.Console) {
				c.ExpectString("What is your name?")
				c.Send("小明")
				c.Send(string(terminal.KeyDelete))
				c.SendLine("")
				c.ExpectEOF()
			},
			"小",
		},
		{
			"Test Input prompt interaction when ask for suggestion with empty value",
			&Input{
				Message: "What is your favorite month?",
				Suggest: func(string) []string {
					return []string{"January", "February"}
				},
			},
			func(c *expect.Console) {
				c.ExpectString("What is your favorite month?")
				c.Send(string(terminal.KeyTab))
				c.ExpectString("January")
				c.ExpectString("February")
				c.SendLine("")
				c.ExpectEOF()
			},
			"January",
		},
		{
			"Test Input prompt interaction when ask for suggestion with some value",
			&Input{
				Message: "What is your favorite month?",
				Suggest: func(string) []string {
					return []string{"February"}
				},
			},
			func(c *expect.Console) {
				c.ExpectString("What is your favorite month?")
				c.Send("feb")
				c.Send(string(terminal.KeyTab))
				c.SendLine("")
				c.ExpectEOF()
			},
			"February",
		},
		{
			"Test Input prompt interaction when ask for suggestion with some value, choosing the second one",
			&Input{
				Message: "What is your favorite month?",
				Suggest: func(string) []string {
					return []string{"January", "February", "March"}
				},
			},
			func(c *expect.Console) {
				c.ExpectString("What is your favorite month?")
				c.Send(string(terminal.KeyTab))
				c.Send(string(terminal.KeyArrowDown))
				c.Send(string(terminal.KeyArrowDown))
				c.SendLine("")
				c.ExpectEOF()
			},
			"March",
		},
		{
			"Test Input prompt interaction when ask for suggestion with some value, choosing the second one",
			&Input{
				Message: "What is your favorite month?",
				Suggest: func(string) []string {
					return []string{"January", "February", "March"}
				},
			},
			func(c *expect.Console) {
				c.ExpectString("What is your favorite month?")
				c.Send(string(terminal.KeyTab))
				c.Send(string(terminal.KeyArrowDown))
				c.Send(string(terminal.KeyArrowDown))
				c.Send(string(terminal.KeyArrowUp))
				c.SendLine("")
				c.ExpectEOF()
			},
			"February",
		},
		{
			"Test Input prompt interaction when ask for suggestion, complementing it and get new suggestions",
			&Input{
				Message: "Where to save it?",
				Suggest: func(complete string) []string {
					if complete == "" {
						return []string{"folder1/", "folder2/", "folder3/"}
					}
					return []string{"folder3/file1.txt", "folder3/file2.txt"}
				},
			},
			func(c *expect.Console) {
				c.ExpectString("Where to save it?")
				c.Send(string(terminal.KeyTab))
				c.ExpectString("folder1/")
				c.Send(string(terminal.KeyArrowDown))
				c.Send(string(terminal.KeyArrowDown))
				c.Send("f")
				c.Send(string(terminal.KeyTab))
				c.ExpectString("folder3/file2.txt")
				c.Send(string(terminal.KeyArrowDown))
				c.SendLine("")
				c.ExpectEOF()
			},
			"folder3/file2.txt",
		},
		{
			"Test Input prompt interaction when asked suggestions, but abort suggestions",
			&Input{
				Message: "Wanna a suggestion?",
				Suggest: func(string) []string {
					return []string{"suggest1", "suggest2"}
				},
			},
			func(c *expect.Console) {
				c.ExpectString("Wanna a suggestion?")
				c.Send("typed answer")
				c.Send(string(terminal.KeyTab))
				c.ExpectString("suggest1")
				c.Send(string(terminal.KeyEscape))
				c.ExpectString("typed answer")
				c.SendLine("")
				c.ExpectEOF()
			},
			"typed answer",
		},
		{
			"Test Input prompt interaction with suggestions, when tabbed with list being shown, should select next suggestion",
			&Input{
				Message: "Choose the special one:",
				Suggest: func(string) []string {
					return []string{"suggest1", "suggest2", "special answer"}
				},
			},
			func(c *expect.Console) {
				c.ExpectString("Choose the special one:")
				c.Send("s")
				c.Send(string(terminal.KeyTab))
				c.ExpectString("suggest1")
				c.ExpectString("suggest2")
				c.ExpectString("special answer")
				c.Send(string(terminal.KeyTab))
				c.Send(string(terminal.KeyTab))
				c.SendLine("")
				c.ExpectEOF()
			},
			"special answer",
		},
		{
			"Test Input prompt must allow moving cursor using right and left arrows",
			&Input{Message: "Filename to save:"},
			func(c *expect.Console) {
				c.ExpectString("Filename to save:")
				c.Send("essay.txt")
				c.Send(string(terminal.KeyArrowLeft))
				c.Send(string(terminal.KeyArrowLeft))
				c.Send(string(terminal.KeyArrowLeft))
				c.Send(string(terminal.KeyArrowLeft))
				c.Send("_final")
				c.Send(string(terminal.KeyArrowRight))
				c.Send(string(terminal.KeyArrowRight))
				c.Send(string(terminal.KeyArrowRight))
				c.Send(string(terminal.KeyArrowRight))
				c.Send(string(terminal.KeyBackspace))
				c.Send(string(terminal.KeyBackspace))
				c.Send(string(terminal.KeyBackspace))
				c.Send("md")
				c.Send(string(terminal.KeyArrowLeft))
				c.Send(string(terminal.KeyArrowLeft))
				c.Send(string(terminal.KeyArrowLeft))
				c.SendLine("2")
				c.ExpectEOF()
			},
			"essay_final2.md",
		},
		{
			"Test Input prompt must allow moving cursor using right and left arrows, even after suggestions",
			&Input{Message: "Filename to save:", Suggest: func(string) []string { return []string{".txt", ".csv", ".go"} }},
			func(c *expect.Console) {
				c.ExpectString("Filename to save:")
				c.Send(string(terminal.KeyTab))
				c.ExpectString(".txt")
				c.ExpectString(".csv")
				c.ExpectString(".go")
				c.Send(string(terminal.KeyTab))
				c.Send(string(terminal.KeyArrowLeft))
				c.Send(string(terminal.KeyArrowLeft))
				c.Send(string(terminal.KeyArrowLeft))
				c.Send(string(terminal.KeyArrowLeft))
				c.Send(string(terminal.KeyArrowLeft))
				c.Send("newtable")
				c.SendLine("")
				c.ExpectEOF()
			},
			"newtable.csv",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			RunPromptTest(t, test)
		})
	}
}

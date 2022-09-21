package survey

import (
	"errors"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/AlecAivazis/survey/v2/core"
	"github.com/AlecAivazis/survey/v2/terminal"
	expect "github.com/Netflix/go-expect"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	// disable color output for all prompts to simplify testing
	core.DisableColor = true
}

type expectConsole interface {
	ExpectString(string)
	ExpectEOF()
	SendLine(string)
	Send(string)
}

type consoleWithErrorHandling struct {
	console *expect.Console
	t       *testing.T
}

func (c *consoleWithErrorHandling) ExpectString(s string) {
	if _, err := c.console.ExpectString(s); err != nil {
		c.t.Helper()
		c.t.Fatalf("ExpectString(%q) = %v", s, err)
	}
}

func (c *consoleWithErrorHandling) SendLine(s string) {
	if _, err := c.console.SendLine(s); err != nil {
		c.t.Helper()
		c.t.Fatalf("SendLine(%q) = %v", s, err)
	}
}

func (c *consoleWithErrorHandling) Send(s string) {
	if _, err := c.console.Send(s); err != nil {
		c.t.Helper()
		c.t.Fatalf("Send(%q) = %v", s, err)
	}
}

func (c *consoleWithErrorHandling) ExpectEOF() {
	if _, err := c.console.ExpectEOF(); err != nil {
		c.t.Helper()
		c.t.Fatalf("ExpectEOF() = %v", err)
	}
}

type PromptTest struct {
	name      string
	prompt    Prompt
	procedure func(expectConsole)
	expected  interface{}
}

func RunPromptTest(t *testing.T, test PromptTest) {
	t.Helper()
	var answer interface{}
	RunTest(t, test.procedure, func(stdio terminal.Stdio) error {
		var err error
		if p, ok := test.prompt.(wantsStdio); ok {
			p.WithStdio(stdio)
		}

		answer, err = test.prompt.Prompt(defaultPromptConfig())
		return err
	})
	require.Equal(t, test.expected, answer)
}

func RunPromptTestKeepFilter(t *testing.T, test PromptTest) {
	t.Helper()
	var answer interface{}
	RunTest(t, test.procedure, func(stdio terminal.Stdio) error {
		var err error
		if p, ok := test.prompt.(wantsStdio); ok {
			p.WithStdio(stdio)
		}
		config := defaultPromptConfig()
		config.KeepFilter = true
		answer, err = test.prompt.Prompt(config)
		return err
	})
	require.Equal(t, test.expected, answer)
}

func RunPromptTestRemoveSelectAll(t *testing.T, test PromptTest) {
	t.Helper()
	var answer interface{}
	RunTest(t, test.procedure, func(stdio terminal.Stdio) error {
		var err error
		if p, ok := test.prompt.(wantsStdio); ok {
			p.WithStdio(stdio)
		}
		config := defaultPromptConfig()
		config.RemoveSelectAll = true
		answer, err = test.prompt.Prompt(config)
		return err
	})
	require.Equal(t, test.expected, answer)
}

func RunPromptTestRemoveSelectNone(t *testing.T, test PromptTest) {
	t.Helper()
	var answer interface{}
	RunTest(t, test.procedure, func(stdio terminal.Stdio) error {
		var err error
		if p, ok := test.prompt.(wantsStdio); ok {
			p.WithStdio(stdio)
		}
		config := defaultPromptConfig()
		config.RemoveSelectNone = true
		answer, err = test.prompt.Prompt(config)
		return err
	})
	require.Equal(t, test.expected, answer)
}

func TestPagination_tooFew(t *testing.T) {
	// a small list of options
	choices := core.OptionAnswerList([]string{"choice1", "choice2", "choice3"})

	// a page bigger than the total number
	pageSize := 4
	// the current selection
	sel := 3

	// compute the page info
	page, idx := paginate(pageSize, choices, sel)

	// make sure we see the full list of options
	assert.Equal(t, choices, page)
	// with the second index highlighted (no change)
	assert.Equal(t, 3, idx)
}

func TestPagination_firstHalf(t *testing.T) {
	// the choices for the test
	choices := core.OptionAnswerList([]string{"choice1", "choice2", "choice3", "choice4", "choice5", "choice6"})

	// section the choices into groups of 4 so the choice is somewhere in the middle
	// to verify there is no displacement of the page
	pageSize := 4
	// test the second item
	sel := 2

	// compute the page info
	page, idx := paginate(pageSize, choices, sel)

	// we should see the first three options
	assert.Equal(t, choices[0:4], page)
	// with the second index highlighted
	assert.Equal(t, 2, idx)
}

func TestPagination_middle(t *testing.T) {
	// the choices for the test
	choices := core.OptionAnswerList([]string{"choice0", "choice1", "choice2", "choice3", "choice4", "choice5"})

	// section the choices into groups of 3
	pageSize := 2
	// test the second item so that we can verify we are in the middle of the list
	sel := 3

	// compute the page info
	page, idx := paginate(pageSize, choices, sel)

	// we should see the first three options
	assert.Equal(t, choices[2:4], page)
	// with the second index highlighted
	assert.Equal(t, 1, idx)
}

func TestPagination_lastHalf(t *testing.T) {
	// the choices for the test
	choices := core.OptionAnswerList([]string{"choice0", "choice1", "choice2", "choice3", "choice4", "choice5"})

	// section the choices into groups of 3
	pageSize := 3
	// test the last item to verify we're not in the middle
	sel := 5

	// compute the page info
	page, idx := paginate(pageSize, choices, sel)

	// we should see the first three options
	assert.Equal(t, choices[3:6], page)
	// we should be at the bottom of the list
	assert.Equal(t, 2, idx)
}

func TestAsk(t *testing.T) {
	if _, err := exec.LookPath("vi"); err != nil {
		t.Skip("vi not found in PATH")
	}

	tests := []struct {
		name      string
		questions []*Question
		procedure func(expectConsole)
		expected  map[string]interface{}
	}{
		{
			"Test Ask for all prompts",
			[]*Question{
				{
					Name: "pizza",
					Prompt: &Confirm{
						Message: "Is pizza your favorite food?",
					},
				},
				{
					Name: "commit-message",
					Prompt: &Editor{
						Editor:  "vi",
						Message: "Edit git commit message",
					},
				},
				/* TODO gets stuck
				{
					Name: "commit-message-validated",
					Prompt: &Editor{
						Editor:  "vi",
						Message: "Edit git commit message",
					},
					Validate: func(v interface{}) error {
						s := v.(string)
						if strings.Contains(s, "invalid") {
							return fmt.Errorf("invalid error message")
						}
						return nil
					},
				},
				*/
				{
					Name: "name",
					Prompt: &Input{
						Message: "What is your name?",
					},
				},
				/* TODO gets stuck
				{
					Name: "day",
					Prompt: &MultiSelect{
						Message: "What days do you prefer:",
						Options: []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"},
					},
				},
				*/
				{
					Name: "password",
					Prompt: &Password{
						Message: "Please type your password",
					},
				},
				{
					Name: "color",
					Prompt: &Select{
						Message: "Choose a color:",
						Options: []string{"red", "blue", "green", "yellow"},
					},
				},
			},
			func(c expectConsole) {
				// Confirm
				c.ExpectString("Is pizza your favorite food? (y/N)")
				c.SendLine("Y")

				// Editor
				c.ExpectString("Edit git commit message [Enter to launch editor]")
				c.SendLine("")
				time.Sleep(time.Millisecond)
				c.Send("ccAdd editor prompt tests\x1b")
				c.SendLine(":wq!")

				/* TODO gets stuck
				// Editor validated
				c.ExpectString("Edit git commit message [Enter to launch editor]")
				c.SendLine("")
				time.Sleep(time.Millisecond)
				c.Send("i invalid input first try\x1b")
				c.SendLine(":wq!")
				time.Sleep(time.Millisecond)
				c.ExpectString("invalid error message")
				c.ExpectString("Edit git commit message [Enter to launch editor]")
				c.SendLine("")
				time.Sleep(time.Millisecond)
				c.ExpectString("first try")
				c.Send("ccAdd editor prompt tests, but validated\x1b")
				c.SendLine(":wq!")
				*/

				// Input
				c.ExpectString("What is your name?")
				c.SendLine("Johnny Appleseed")

				/* TODO gets stuck
				// MultiSelect
				c.ExpectString("What days do you prefer:  [Use arrows to move, space to select, type to filter]")
				// Select Monday.
				c.Send(string(terminal.KeyArrowDown))
				c.Send(" ")
				// Select Wednesday.
				c.Send(string(terminal.KeyArrowDown))
				c.Send(string(terminal.KeyArrowDown))
				c.SendLine(" ")
				*/

				// Password
				c.ExpectString("Please type your password")
				c.Send("secret")
				c.SendLine("")

				// Select
				c.ExpectString("Choose a color:  [Use arrows to move, type to filter]")
				c.SendLine("yellow")
				c.ExpectEOF()
			},
			map[string]interface{}{
				"pizza":          true,
				"commit-message": "Add editor prompt tests\n",
				/* TODO
				"commit-message-validated": "Add editor prompt tests, but validated\n",
				*/
				"name": "Johnny Appleseed",
				/* TODO
				"day":                      []string{"Monday", "Wednesday"},
				*/
				"password": "secret",
				"color":    core.OptionAnswer{Index: 3, Value: "yellow"},
			},
		},
		{
			"Test Ask with validate survey.Required",
			[]*Question{
				{
					Name: "name",
					Prompt: &Input{
						Message: "What is your name?",
					},
					Validate: Required,
				},
			},
			func(c expectConsole) {
				c.ExpectString("What is your name?")
				c.SendLine("")
				c.ExpectString("Sorry, your reply was invalid: Value is required")
				time.Sleep(time.Millisecond)
				c.SendLine("Johnny Appleseed")
				c.ExpectEOF()
			},
			map[string]interface{}{
				"name": "Johnny Appleseed",
			},
		},
		{
			"Test Ask with transformer survey.ToLower",
			[]*Question{
				{
					Name: "name",
					Prompt: &Input{
						Message: "What is your name?",
					},
					Transform: ToLower,
				},
			},
			func(c expectConsole) {
				c.ExpectString("What is your name?")
				c.SendLine("Johnny Appleseed")
				c.ExpectString("What is your name? johnny appleseed")
				c.ExpectEOF()
			},
			map[string]interface{}{
				"name": "johnny appleseed",
			},
		},
	}

	for _, test := range tests {
		// Capture range variable.
		test := test
		t.Run(test.name, func(t *testing.T) {
			answers := make(map[string]interface{})
			RunTest(t, test.procedure, func(stdio terminal.Stdio) error {
				return Ask(test.questions, &answers, WithStdio(stdio.In, stdio.Out, stdio.Err))
			})
			require.Equal(t, test.expected, answers)
		})
	}
}

func TestAsk_returnsErrorIfTargetIsNil(t *testing.T) {
	// pass an empty place to leave the answers
	err := Ask([]*Question{}, nil)

	// if we didn't get an error
	if err == nil {
		// the test failed
		t.Error("Did not encounter error when asking with no where to record.")
	}
}

func Test_computeCursorOffset_MultiSelect(t *testing.T) {
	tests := []struct {
		name      string
		ix        int
		opts      []core.OptionAnswer
		termWidth int
		want      int
	}{
		{
			name: "no opts",
			ix:   0,
			opts: []core.OptionAnswer{},
			want: 0,
		},
		{
			name: "one opt",
			ix:   0,
			opts: core.OptionAnswerList([]string{"one"}),
			want: 1,
		},
		{
			name: "multiple opt",
			opts: core.OptionAnswerList([]string{"one", "two"}),
			ix:   0,
			want: 2,
		},
		{
			name: "first choice",
			opts: core.OptionAnswerList([]string{"one", "two", "three", "four", "five"}),
			ix:   0,
			want: 5,
		},
		{
			name: "mid choice",
			opts: core.OptionAnswerList([]string{"one", "two", "three", "four", "five"}),
			ix:   2,
			want: 3,
		},
		{
			name: "last choice",
			opts: core.OptionAnswerList([]string{"one", "two", "three", "four", "five"}),
			ix:   4,
			want: 1,
		},
		{
			name: "wide choices, uneven",
			opts: core.OptionAnswerList([]string{
				"wide one wide one wide one",
				"two", "three",
				"wide four wide four wide four",
				"five", "six"}),
			termWidth: 20,
			ix:        0,
			want:      8,
		},
		{
			name: "wide choices, even",
			opts: core.OptionAnswerList([]string{
				"wide one wide one wide one",
				"two", "three",
				"012345678901",
				"five", "six"}),
			termWidth: 20,
			ix:        0,
			want:      7,
		},
		{
			name: "wide choices, wide before idx",
			opts: core.OptionAnswerList([]string{
				"wide one wide one wide one",
				"wide two wide two wide two",
				"three", "four", "five", "six"}),
			termWidth: 20,
			ix:        2,
			want:      4,
		},
	}
	for _, tt := range tests {
		if tt.termWidth == 0 {
			tt.termWidth = 100
		}
		tmpl := MultiSelectQuestionTemplate
		data := MultiSelectTemplateData{
			SelectedIndex: tt.ix,
			Config:        defaultPromptConfig(),
		}
		t.Run(tt.name, func(t *testing.T) {
			if got := computeCursorOffset(tmpl, data, tt.opts, tt.ix, tt.termWidth); got != tt.want {
				t.Errorf("computeCursorOffset() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_computeCursorOffset_Select(t *testing.T) {
	tests := []struct {
		name      string
		ix        int
		opts      []core.OptionAnswer
		termWidth int
		want      int
	}{
		{
			name: "no opts",
			ix:   0,
			opts: []core.OptionAnswer{},
			want: 0,
		},
		{
			name: "one opt",
			ix:   0,
			opts: core.OptionAnswerList([]string{"one"}),
			want: 1,
		},
		{
			name: "multiple opt",
			opts: core.OptionAnswerList([]string{"one", "two"}),
			ix:   0,
			want: 2,
		},
		{
			name: "first choice",
			opts: core.OptionAnswerList([]string{"one", "two", "three", "four", "five"}),
			ix:   0,
			want: 5,
		},
		{
			name: "mid choice",
			opts: core.OptionAnswerList([]string{"one", "two", "three", "four", "five"}),
			ix:   2,
			want: 3,
		},
		{
			name: "last choice",
			opts: core.OptionAnswerList([]string{"one", "two", "three", "four", "five"}),
			ix:   4,
			want: 1,
		},
		{
			name: "wide choices, uneven",
			opts: core.OptionAnswerList([]string{
				"wide one wide one wide one",
				"two", "three",
				"wide four wide four wide four",
				"five", "six"}),
			termWidth: 20,
			ix:        0,
			want:      8,
		},
		{
			name: "wide choices, even",
			opts: core.OptionAnswerList([]string{
				"wide one wide one wide one",
				"two", "three",
				"01234567890123456",
				"five", "six"}),
			termWidth: 20,
			ix:        0,
			want:      7,
		},
		{
			name: "wide choices, wide before idx",
			opts: core.OptionAnswerList([]string{
				"wide one wide one wide one",
				"wide two wide two wide two",
				"three", "four", "five", "six"}),
			termWidth: 20,
			ix:        2,
			want:      4,
		},
	}
	for _, tt := range tests {
		if tt.termWidth == 0 {
			tt.termWidth = 100
		}
		tmpl := SelectQuestionTemplate
		data := SelectTemplateData{
			SelectedIndex: tt.ix,
			Config:        defaultPromptConfig(),
		}
		t.Run(tt.name, func(t *testing.T) {
			if got := computeCursorOffset(tmpl, data, tt.opts, tt.ix, tt.termWidth); got != tt.want {
				t.Errorf("computeCursorOffset() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAsk_Validation(t *testing.T) {
	p := &mockPrompt{
		answers: []string{"", "company", "COM", "com"},
	}

	var res struct {
		TLDN string
	}
	err := Ask([]*Question{
		{
			Name:   "TLDN",
			Prompt: p,
			Validate: func(v interface{}) error {
				s := v.(string)
				if strings.ToLower(s) != s {
					return errors.New("value contains uppercase characters")
				}
				return nil
			},
		},
	}, &res, WithValidator(MinLength(1)), WithValidator(MaxLength(5)))
	if err != nil {
		t.Fatalf("Ask() = %v", err)
	}

	if res.TLDN != "com" {
		t.Errorf("answer: %q, want %q", res.TLDN, "com")
	}
	if p.cleanups != 1 {
		t.Errorf("cleanups: %d, want %d", p.cleanups, 1)
	}
	if err1 := p.printedErrors[0].Error(); err1 != "value is too short. Min length is 1" {
		t.Errorf("printed error 1: %q, want %q", err1, "value is too short. Min length is 1")
	}
	if err2 := p.printedErrors[1].Error(); err2 != "value is too long. Max length is 5" {
		t.Errorf("printed error 2: %q, want %q", err2, "value is too long. Max length is 5")
	}
	if err3 := p.printedErrors[2].Error(); err3 != "value contains uppercase characters" {
		t.Errorf("printed error 2: %q, want %q", err3, "value contains uppercase characters")
	}
}

type mockPrompt struct {
	index         int
	answers       []string
	cleanups      int
	printedErrors []error
}

func (p *mockPrompt) Prompt(*PromptConfig) (interface{}, error) {
	if p.index >= len(p.answers) {
		return nil, errors.New("no more answers")
	}
	val := p.answers[p.index]
	p.index++
	return val, nil
}

func (p *mockPrompt) Cleanup(*PromptConfig, interface{}) error {
	p.cleanups++
	return nil
}

func (p *mockPrompt) Error(_ *PromptConfig, err error) error {
	p.printedErrors = append(p.printedErrors, err)
	return nil
}

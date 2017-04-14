package survey

import (
	"testing"

	"github.com/alecaivazis/survey/core"
)

func init() {
	// disable color output for all prompts to simplify testing
	core.DisableColor = true
}

func TestCanFormatSelectOptions(t *testing.T) {

	prompt := &Select{
		Options: []string{"foo", "bar", "baz", "buz"},
		Default: "baz",
	}
	// TODO: figure out a way for the test to actually test this bit of code
	prompt.SelectedIndex = 2

	actual, err := core.RunTemplate(
		selectChoicesTemplate,
		SelectTemplateData{Select: *prompt},
	)

	if err != nil {
		t.Errorf("Failed to run template to format choice options: %s", err)
	}

	expected := `  foo
  bar
> baz
  buz
`

	if actual != expected {
		t.Errorf("Formatted choice options were not formatted correctly. Found:\n%s\nExpected:\n%s", actual, expected)
	}
}

func TestSelectFormatQuestion(t *testing.T) {

	prompt := &Select{
		Message: "Pick your word:",
		Options: []string{"foo", "bar", "baz", "buz"},
		Default: "baz",
	}

	actual, err := core.RunTemplate(
		selectQuestionTemplate,
		SelectTemplateData{Select: *prompt},
	)
	if err != nil {
		t.Errorf("Failed to run template to format Select question: %s", err)
	}

	expected := `? Pick your word: `

	if actual != expected {
		t.Errorf("Formatted Select question was not formatted correctly. Found:\n%s\nExpected:\n%s", actual, expected)
	}
}

func TestSelectFormatAnswer(t *testing.T) {

	prompt := &Select{
		Message: "Pick your word:",
		Options: []string{"foo", "bar", "baz", "buz"},
		Default: "baz",
	}

	actual, err := core.RunTemplate(
		selectQuestionTemplate,
		SelectTemplateData{Select: *prompt, Answer: "buz"},
	)
	if err != nil {
		t.Errorf("Failed to run template to format Select answer: %s", err)
	}

	expected := `? Pick your word: buz`

	if actual != expected {
		t.Errorf("Formatted Select answer was not formatted correctly. Found:\n%s\nExpected:\n%s", actual, expected)
	}
}

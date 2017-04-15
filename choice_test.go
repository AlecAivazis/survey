package survey

import (
	"testing"
)

func init() {
	// disable color output for all prompts to simplify testing
	DisableColor = true
}

func TestCanFormatSelectOptions(t *testing.T) {

	prompt := &Choice{
		Choices: []string{"foo", "bar", "baz", "buz"},
		Default: "baz",
	}
	// TODO: figure out a way for the test to actually test this bit of code
	prompt.SelectedIndex = 2

	actual, err := RunTemplate(
		SelectChoicesTemplate,
		SelectTemplateData{Select: *prompt},
	)

	if err != nil {
		t.Errorf("Failed to run template to format choice choices: %s", err)
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

	prompt := &Choice{
		Message: "Pick your word:",
		Choices: []string{"foo", "bar", "baz", "buz"},
		Default: "baz",
	}

	actual, err := RunTemplate(
		SelectQuestionTemplate,
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

	prompt := &Choice{
		Message: "Pick your word:",
		Choices: []string{"foo", "bar", "baz", "buz"},
		Default: "baz",
	}

	actual, err := RunTemplate(
		SelectQuestionTemplate,
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

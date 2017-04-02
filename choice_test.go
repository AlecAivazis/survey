package survey

import "testing"

func init() {
	// disable color output for all prompts to simplify testing
	DisableColor = true
}

func TestCanFormatChoiceOptions(t *testing.T) {

	prompt := &Choice{
		Choices: []string{"foo", "bar", "baz", "buz"},
	}

	actual, err := runTemplate(
		ChoiceChoicesTemplate,
		ChoiceTemplateData{Choice: *prompt, Selected: 2},
	)

	if err != nil {
		t.Errorf("Failed to run template to format choice options: %s", err)
	}

	expected := `  foo
  bar
➤ baz
  buz
`

	if actual != expected {
		t.Errorf("Formatted choice options were not formatted correctly. Found:\n%s\nExpected:\n%s", actual, expected)
	}
}

func TestChoiceFormatQuestion(t *testing.T) {

	prompt := &Choice{
		Message: "Pick your word:",
		Choices: []string{"foo", "bar", "baz", "buz"},
		Default: "baz",
	}

	actual, err := runTemplate(
		ChoiceQuestionTemplate,
		ChoiceTemplateData{Choice: *prompt},
	)
	if err != nil {
		t.Errorf("Failed to run template to format choice question: %s", err)
	}

	expected := `? Pick your word: `

	if actual != expected {
		t.Errorf("Formatted choice question was not formatted correctly. Found:\n%s\nExpected:\n%s", actual, expected)
	}
}

func TestChoiceFormatAnswer(t *testing.T) {

	prompt := &Choice{
		Message: "Pick your word:",
		Choices: []string{"foo", "bar", "baz", "buz"},
		Default: "baz",
	}

	actual, err := runTemplate(
		ChoiceQuestionTemplate,
		ChoiceTemplateData{Choice: *prompt, Answer: "buz"},
	)
	if err != nil {
		t.Errorf("Failed to run template to format choice answer: %s", err)
	}

	expected := `? Pick your word: buz`

	if actual != expected {
		t.Errorf("Formatted choice answer was not formatted correctly. Found:\n%s\nExpected:\n%s", actual, expected)
	}
}

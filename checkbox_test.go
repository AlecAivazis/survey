package survey

import "testing"

func init() {
	// disable color output for all prompts to simplify testing
	DisableColor = true
}

func TestCanFormatCheckboxOptions(t *testing.T) {

	prompt := &Checkbox{
		Options:  []string{"foo", "bar", "baz", "buz"},
		Defaults: []string{"bar", "buz"},
	}

	actual, err := runTemplate(
		CheckboxOptionsTemplate,
		CheckboxTemplateData{Checkbox: *prompt, Selected: 2, Checked: map[int]bool{1: true, 3: true}},
	)

	if err != nil {
		t.Errorf("Failed to run template to format checkbox options: %s", err)
	}

	expected := ` ◯ foo
 ◉ bar
❯◯ baz
 ◉ buz
`

	if actual != expected {
		t.Errorf("Formatted checkbox options were not formatted correctly. Found:\n%s\nExpected:\n%s", actual, expected)
	}
}

func TestCheckboxFormatQuestion(t *testing.T) {

	prompt := &Checkbox{
		Message:  "Pick your words:",
		Options:  []string{"foo", "bar", "baz", "buz"},
		Defaults: []string{"bar", "buz"},
	}

	actual, err := runTemplate(
		CheckboxQuestionTemplate,
		CheckboxTemplateData{Checkbox: *prompt},
	)
	if err != nil {
		t.Errorf("Failed to run template to format checkbox question: %s", err)
	}

	expected := `? Pick your words: `

	if actual != expected {
		t.Errorf("Formatted checkbox question was not formatted correctly. Found:\n%s\nExpected:\n%s", actual, expected)
	}
}

func TestCheckboxFormatAnswer(t *testing.T) {

	prompt := &Checkbox{
		Message:  "Pick your words:",
		Options:  []string{"foo", "bar", "baz", "buz"},
		Defaults: []string{"bar", "buz"},
	}

	actual, err := runTemplate(
		CheckboxQuestionTemplate,
		CheckboxTemplateData{Checkbox: *prompt, Answer: "bar,buz"},
	)
	if err != nil {
		t.Errorf("Failed to run template to format checkbox answer: %s", err)
	}

	expected := `? Pick your words: bar,buz`

	if actual != expected {
		t.Errorf("Formatted checkbox answer was not formatted correctly. Found:\n%s\nExpected:\n%s", actual, expected)
	}
}

package survey

import "testing"

func init() {
	// disable color output for all prompts to simplify testing
	DisableColor = true
}

func TestCanFormatMultiChoiceOptions(t *testing.T) {

	prompt := &MultiChoice{
		Options:  []string{"foo", "bar", "baz", "buz"},
		Defaults: []string{"bar", "buz"},
	}

	actual, err := RunTemplate(
		MultiChoiceOptionsTemplate,
		MultiChoiceTemplateData{
			MultiChoice:   *prompt,
			SelectedIndex: 2,
			Checked:       map[int]bool{1: true, 3: true},
		},
	)

	if err != nil {
		t.Errorf("Failed to run template to format checkbox options: %s", err)
	}

	expected := `  ◯  foo
  ◉  bar
❯ ◯  baz
  ◉  buz
`

	if actual != expected {
		t.Errorf("Formatted checkbox options were not formatted correctly. Found:\n%s\nExpected:\n%s", actual, expected)
	}
}

func TestMultiChoiceFormatQuestion(t *testing.T) {

	prompt := &MultiChoice{
		Message:  "Pick your words:",
		Options:  []string{"foo", "bar", "baz", "buz"},
		Defaults: []string{"bar", "buz"},
	}

	actual, err := RunTemplate(
		MultiChoiceQuestionTemplate,
		MultiChoiceTemplateData{MultiChoice: *prompt},
	)
	if err != nil {
		t.Errorf("Failed to run template to format checkbox question: %s", err)
	}

	expected := `? Pick your words: `

	if actual != expected {
		t.Errorf("Formatted checkbox question was not formatted correctly. Found:\n%s\nExpected:\n%s", actual, expected)
	}
}

func TestMultiChoiceFormatAnswer(t *testing.T) {

	prompt := &MultiChoice{
		Message:  "Pick your words:",
		Options:  []string{"foo", "bar", "baz", "buz"},
		Defaults: []string{"bar", "buz"},
	}

	actual, err := RunTemplate(
		MultiChoiceQuestionTemplate,
		MultiChoiceTemplateData{MultiChoice: *prompt, Answer: []string{"foo", "buz"}},
	)
	if err != nil {
		t.Errorf("Failed to run template to format checkbox answer: %s", err)
	}

	expected := `? Pick your words: ["foo" "buz"]`

	if actual != expected {
		t.Errorf("Formatted checkbox answer was not formatted correctly. Found:\n%s\nExpected:\n%s", actual, expected)
	}
}

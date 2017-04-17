package survey

import (
	"strings"
	"testing"

	"github.com/alecaivazis/survey/core"
)

func init() {
	// disable color output for all prompts to simplify testing
	core.DisableColor = true
}

func TestCanFormatMultiSelectOptions(t *testing.T) {

	prompt := &MultiSelect{
		Options: []string{"foo", "bar", "baz", "buz"},
		Default: []string{"bar", "buz"},
	}

	actual, err := core.RunTemplate(
		MultiSelectOptionsTemplate,
		MultiSelectTemplateData{
			MultiSelect:   *prompt,
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

func TestMultiSelectFormatQuestion(t *testing.T) {

	prompt := &MultiSelect{
		Message: "Pick your words:",
		Options: []string{"foo", "bar", "baz", "buz"},
		Default: []string{"bar", "buz"},
	}

	actual, err := core.RunTemplate(
		MultiSelectQuestionTemplate,
		MultiSelectTemplateData{MultiSelect: *prompt},
	)
	if err != nil {
		t.Errorf("Failed to run template to format checkbox question: %s", err)
	}

	expected := `? Pick your words: `

	if actual != expected {
		t.Errorf("Formatted checkbox question was not formatted correctly. Found:\n%s\nExpected:\n%s", actual, expected)
	}
}

func TestMultiSelectFormatAnswer(t *testing.T) {

	prompt := &MultiSelect{
		Message: "Pick your words:",
		Options: []string{"foo", "bar", "baz", "buz"},
		Default: []string{"bar", "buz"},
	}

	actual, err := core.RunTemplate(
		MultiSelectQuestionTemplate,
		MultiSelectTemplateData{MultiSelect: *prompt, Answer: strings.Join([]string{"foo", "buz"}, ", ")},
	)
	if err != nil {
		t.Errorf("Failed to run template to format checkbox answer: %s", err)
	}

	expected := `? Pick your words: foo, buz`

	if actual != expected {
		t.Errorf("Formatted checkbox answer was not formatted correctly. Found:\n%s\nExpected:\n%s", actual, expected)
	}
}

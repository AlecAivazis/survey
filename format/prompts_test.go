package format

import (
	"fmt"
	"testing"
)

func TestCanFormatChoiceOption(t *testing.T) {
	// the string to format
	str := "hello"
	// make sure there is a tab before the option
	if ChoiceOption(str, false) != fmt.Sprintf("%s%s", ChoiceNotSelected, str) {
		t.Error("Could not format option")
	}
}

func TestCanFormatSelectedChoiceOption(t *testing.T) {
	// the string to format
	str := "hello"
	// make sure there is a tab before the option
	if ChoiceOption(str, true) != fmt.Sprintf("%s%s", ChoiceSelected, str) {
		t.Error("Could not format selected option")
	}
}

func TestCanFormatAsk(t *testing.T) {
	// the string to format
	str := "hello"
	// make sure there is a tab before the option
	if Ask(str) != fmt.Sprintf("%s%s", Question, str) {
		t.Error("Could not format ask")
	}
}

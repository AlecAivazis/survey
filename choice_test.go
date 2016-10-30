package probe

import (
	"fmt"
	"testing"

	"github.com/alecaivazis/probe/format"
)

func TestCanFormatChoiceOption(t *testing.T) {
	// the string to format
	str := "hello"
	// make sure there is a tab before the option
	if format.FormatChoiceOption(str, false) != fmt.Sprintf("%s%s", format.ChoiceNotSelected, str) {
		t.Error("Could not format option")
	}
}
func TestCanFormatSelectedChoiceOption(t *testing.T) {
	// the string to format
	str := "hello"
	// make sure there is a tab before the option
	if format.FormatChoiceOption(str, true) != fmt.Sprintf("%s%s", format.ChoiceSelected, str) {
		t.Error("Could not format selected option")
	}
}

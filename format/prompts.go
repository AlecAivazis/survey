package format

import (
	"fmt"
)

// FormatChoiceOption prepares the string to be printed like an option in a
// choice list.
func FormatChoiceOption(opt string, selected bool) string {
	// the tab to use depends on wether the option is selected
	var tab string
	if selected {
		tab = ChoiceSelected
	} else {
		tab = ChoiceNotSelected
	}
	return fmt.Sprintf("%s%s", tab, opt)
}

// FormatAsk prepares a string to be printed like the first line
// of a prompt
func FormatAsk(q string) string {
	return fmt.Sprintf("%s%v", Question, q)
}

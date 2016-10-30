package format

import (
	"fmt"
)

// FormatChoiceOption prepares the string to be printed like an option in a
// choice list.
func ChoiceOption(opt string, selected bool) string {
	// if we are rendering the selected option
	if selected {
		// paint the line blue
		return fmt.Sprint(SelectedColor, ChoiceSelected, opt, ResetColor)
	} else {
		// if its not selected, treat it like normal
		return fmt.Sprint(ChoiceNotSelected, opt)
	}
}

// FormatAsk prepares a string to be printed like the first line
// of a prompt
func Ask(q string) string {
	return fmt.Sprintf("%s%v", Question, q)
}

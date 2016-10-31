package format

import (
	"fmt"
	"github.com/ttacon/chalk"
)

// FormatChoiceOption prepares the string to be printed like an option in a
// choice list.
func ChoiceOption(opt string, selected bool) string {
	// if we are rendering the selected option
	if selected {
		// paint the line blue
		return fmt.Sprint(SelectedColor, ChoiceSelected, opt, ResetFormat)
	} else {
		// if its not selected, treat it like normal
		return fmt.Sprint(ChoiceNotSelected, opt)
	}
}

// FormatAsk prepares a string to be printed like the first line
// of a prompt
func Ask(q string, def string) string {
	// get the message for the question
	msg := chalk.Bold.TextStyle(fmt.Sprint(QuestionColor, Question, ResetFormat, q))
	// the default default message is empty
	defMsg := ""
	// if the user passed a default value
	if def != "" {
		// the default message
		defMsg = fmt.Sprint(DefaultColor, "(", def, ") ", ResetFormat)
	}

	// return the combination of the two message
	return fmt.Sprintf("%s%s", msg, defMsg)
}

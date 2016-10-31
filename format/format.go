package format

import "github.com/ttacon/chalk"

var (
	Tab               = "  "
	Question          = "? "
	ChoiceSelected    = "> "
	ChoiceNotSelected = Tab
	SelectedColor     = chalk.Cyan
	ResetFormat       = chalk.Reset
	QuestionColor     = chalk.Green
	Error             = chalk.Red
	DefaultColor      = "\u001b[37m"
)

package format

import "github.com/ttacon/chalk"

var (
	Tab               = "  "
	Question          = "? "
	ChoiceSelected    = "> "
	Error             = "! "
	ChoiceNotSelected = Tab
	SelectedColor     = chalk.Cyan
	ResetFormat       = chalk.Reset
	QuestionColor     = chalk.Green
	ErrorColor        = chalk.Red
	DefaultColor      = "\u001b[37m"
)

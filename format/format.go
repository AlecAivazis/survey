package format

import "github.com/ttacon/chalk"

var (
	TabEmpty          = "   "
	Question          = " ? "
	ChoiceSelected    = " > "
	ChoiceNotSelected = TabEmpty
	SelectedColor     = chalk.Cyan
	ResetFormat       = chalk.Reset
	QuestionColor     = chalk.Green
)

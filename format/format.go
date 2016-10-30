package format

import "github.com/ttacon/chalk"

var (
	TabEmpty          = "   "
	Question          = " ? "
	ChoiceSelected    = " > "
	ChoiceNotSelected = TabEmpty
	SelectedColor     = chalk.Blue
	ResetColor        = chalk.Reset
)

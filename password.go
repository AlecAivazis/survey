package survey

import (
	"os"

	"github.com/AlecAivazis/survey/core"
	"github.com/AlecAivazis/survey/terminal"
)

// Password is like a normal Input but the text shows up as *'s and
// there is no default.
type Password struct {
	Message string
}

// Templates with Color formatting. See Documentation: https://github.com/mgutz/ansi#style-format
var PasswordQuestionTemplate = `
{{- color "green+hb"}}{{ QuestionIcon }} {{color "reset"}}
{{- color "default+hb"}}{{ .Message }} {{color "reset"}}`

func (p *Password) Prompt() (line interface{}, err error) {
	// render the question template
	out, err := core.RunTemplate(
		PasswordQuestionTemplate,
		*p,
	)
	terminal.Print(out)
	if err != nil {
		return "", err
	}

	rr := terminal.NewRuneReader(os.Stdin)
	rr.SetTermMode()
	defer rr.RestoreTermMode()
	return rr.ReadLine('*')
}

// Cleanup hides the string with a fixed number of characters.
func (prompt *Password) Cleanup(val interface{}) error {
	return nil
}

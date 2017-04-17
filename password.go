package survey

import (
	"github.com/chzyer/readline"

	"github.com/alecaivazis/survey/core"
)

// Password is like a normal Input but the text shows up as *'s and
// there is no default.
type Password struct {
	Message string
}

// Templates with Color formatting. See Documentation: https://github.com/mgutz/ansi#style-format
var PasswordQuestionTemplate = `
{{- color "green+hb"}}? {{color "reset"}}
{{- color "default+hb"}}{{ .Message }} {{color "reset"}}`

func (p *Password) Prompt(rl *readline.Instance) (line interface{}, err error) {
	// render the question template
	out, err := core.RunTemplate(
		PasswordQuestionTemplate,
		*p,
	)
	if err != nil {
		return "", err
	}

	// a configuration for the password prompt
	config := rl.GenPasswordConfig()
	// use the right prompt (make sure there is an empty space after the prompt)
	config.Prompt = out + " "

	config.MaskRune = '*'

	// ask for the user's Password
	pass, err := rl.ReadPasswordWithConfig(config)
	// we're done here
	return string(pass), err
}

// Cleanup hides the string with a fixed number of characters.
func (prompt *Password) Cleanup(rl *readline.Instance, val interface{}) error {
	return nil
}

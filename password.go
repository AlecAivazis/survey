package survey

import (
	"bytes"
	"fmt"

	"github.com/AlecAivazis/survey/terminal"
	"github.com/chzyer/readline"
)

// Password is like a normal Input but the text shows up as *'s and
// there is no default.
type Password struct {
	renderer
	Message string
}

// Templates with Color formatting. See Documentation: https://github.com/mgutz/ansi#style-format
var PasswordQuestionTemplate = `
{{- color "green+hb"}}? {{color "reset"}}
{{- color "default+hb"}}{{ .Message }} {{color "reset"}}`

type filteredReadlinePasswordOutput struct {
	len int
}

// Write will be called by readline to echo the masked characters to the console.  It assumes
// full control of the output with prompting, but since survey is printing the prompt readline
// does not play nice here.  Write is called with clear-line statements which need to be
// ignored then it is called with the full buffered output. So if you type "hi" this is the
// sequence that will be sent to Write
// 0: "\033[J\033[2K\r"
// 1: " \b"
// 2: "\033[J\033[2K\r"
// 3: "*" (this is the masked "h" from "hi")
// 4: "\033[J\033[2K\r"
// 5: "**" (this is the masked "hi" from "hi")
func (o *filteredReadlinePasswordOutput) Write(p []byte) (n int, err error) {
	// we need to ignore the reset sequence unless we have already printed masked
	// characters in which case we need to move the cursor back for however many
	// characters we have printed previously
	if bytes.Equal(p, []byte("\033[J\033[2K\r")) {
		if o.len > 0 {
			// we have previously printed, so move back some number of characters
			// then clear to end of line
			terminal.CursorBack(o.len)
			terminal.EraseLine(terminal.ERASE_LINE_END)
		}
		return 0, nil
	}

	// when the buffer is empty readline.RuneBuffer.output will print
	// this sequence to move the cusror forward then backwards
	// we dont want to count these so just print and return
	if bytes.Equal(p, []byte(" \b")) {
		return fmt.Fprintf(readline.Stdout, "%s", p)
	}

	o.len = len(p)
	return fmt.Fprintf(readline.Stdout, "%s", p)
}

func (p *Password) Prompt(rl *readline.Instance) (line interface{}, err error) {
	// render the question template
	err = p.render(
		PasswordQuestionTemplate,
		*p,
	)
	if err != nil {
		return "", err
	}

	// a configuration for the password prompt
	config := rl.GenPasswordConfig()
	config.Stdout = &filteredReadlinePasswordOutput{}
	config.MaskRune = '*'
	rl.SetConfig(config)

	// ask for the user's Password
	pass, err := rl.ReadPasswordWithConfig(config)
	// we're done here
	return string(pass), err
}

// Cleanup hides the string with a fixed number of characters.
func (prompt *Password) Cleanup(rl *readline.Instance, val interface{}) error {
	return nil
}

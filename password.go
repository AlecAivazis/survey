package survey

import (
	"fmt"
	"os"
	"unicode"

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

	pass := []rune{}
	for {
		r, _, _ := rr.ReadRune()
		if r == '\r' || r == '\n' {
			terminal.Print("\r\n")
			break
		}
		if r == terminal.KeyInterrupt {
			return "", fmt.Errorf("interrupt")
		}
		if r == terminal.KeyEndTransmission {
			break
		}
		// allow for backspace/delete editing of password
		if r == terminal.KeyBackspace || r == terminal.KeyDelete {
			if len(pass) > 0 {
				pass = pass[:len(pass)-1]
				terminal.CursorBack(1)
				terminal.EraseLine(terminal.ERASE_LINE_END)
			}
			continue
		}
		if unicode.IsPrint(r) {
			pass = append(pass, r)
			terminal.Print("*")
		}
	}
	return string(pass), err
}

// Cleanup hides the string with a fixed number of characters.
func (prompt *Password) Cleanup(val interface{}) error {
	return nil
}

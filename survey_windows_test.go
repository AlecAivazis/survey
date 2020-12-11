package survey

import (
	"testing"

	"github.com/AlecAivazis/survey/v2/terminal"
	expect "github.com/Netflix/go-expect"
)

func RunTest(t *testing.T, procedure func(*expect.Console), test func(terminal.Stdio) error) error {
	t.Skip("Windows does not support psuedoterminals")
	return nil
}

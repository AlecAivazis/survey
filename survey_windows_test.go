package survey

import (
	"testing"

	expect "github.com/Netflix/go-expect"
	"github.com/AlecAivazis/survey/v2/terminal"
)

func RunTest(t *testing.T, procedure func(*expect.Console), test func(terminal.Stdio) error) {
	t.Skip("Windows does not support psuedoterminals")
}

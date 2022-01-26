package survey

import (
	"testing"

	"github.com/AlecAivazis/survey/v2/terminal"
	expect "github.com/Netflix/go-expect"
)

func RunTest(t *testing.T, procedure func(expectConsole), test func(terminal.Stdio) error) {
	t.Skip("Windows does not support psuedoterminals")
}

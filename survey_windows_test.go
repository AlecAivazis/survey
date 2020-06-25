package survey

import (
	"testing"

	"github.com/tomercy/survey/terminal"
	expect "github.com/Netflix/go-expect"
)

func RunTest(t *testing.T, procedure func(*expect.Console), test func(terminal.Stdio) error) {
	t.Skip("Windows does not support psuedoterminals")
}

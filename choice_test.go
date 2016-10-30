package probe

import (
	"strings"
	"testing"

	"github.com/alecaivazis/probe/format"
)

func TestCanFormatChoiceOptions(t *testing.T) {
	// for now, options are just a slice of strings
	opts := []string{"foo", "bar", "baz", "buz"}
	// the current selection
	sel := 2

	// the formatted array
	fOpts := []string{
		format.ChoiceOption(opts[0], false),
		format.ChoiceOption(opts[1], false),
		format.ChoiceOption(opts[2], true),
		format.ChoiceOption(opts[3], false),
	}

	var (
		expected = strings.Join(fOpts, "\n")
		actual   = formatChoiceOptions(opts, sel)
	)

	// since format takes care of the actual formatting, we just need to make
	// sure that we have a formatted option on each line
	if actual != expected {
		t.Errorf("Formatted choice options were not formatted correctly. Found:\n%s", actual)
	}
}

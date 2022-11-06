package terminal

import (
	"testing"
)

func TestRuneWidthInvisible(t *testing.T) {
	var example rune = '⁣'
	expected := 0
	actual := runeWidth(example)
	if actual != expected {
		t.Errorf("Expected '%c' to have width %d, found %d", example, expected, actual)
	}
}

func TestRuneWidthNormal(t *testing.T) {
	var example rune = 'a'
	expected := 1
	actual := runeWidth(example)
	if actual != expected {
		t.Errorf("Expected '%c' to have width %d, found %d", example, expected, actual)
	}
}

func TestRuneWidthWide(t *testing.T) {
	var example rune = '错'
	expected := 2
	actual := runeWidth(example)
	if actual != expected {
		t.Errorf("Expected '%c' to have width %d, found %d", example, expected, actual)
	}
}

func TestStringWidthEmpty(t *testing.T) {
	example := ""
	expected := 0
	actual := StringWidth(example)
	if actual != expected {
		t.Errorf("Expected '%s' to have width %d, found %d", example, expected, actual)
	}
}

func TestStringWidthNormal(t *testing.T) {
	example := "Green"
	expected := 5
	actual := StringWidth(example)
	if actual != expected {
		t.Errorf("Expected '%s' to have width %d, found %d", example, expected, actual)
	}
}

func TestStringWidthFormat(t *testing.T) {
	example := "\033[31mRed\033[0m"
	expected := 3
	actual := StringWidth(example)
	if actual != expected {
		t.Errorf("Expected '%s' to have width %d, found %d", example, expected, actual)
	}

	example = "\033[1;34mbold\033[21mblue\033[0m"
	expected = 8
	actual = StringWidth(example)
	if actual != expected {
		t.Errorf("Expected '%s' to have width %d, found %d", example, expected, actual)
	}
}

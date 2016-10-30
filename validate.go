package probe

import (
	"errors"
)

func Required(str string) error {
	// if the string is empty
	if str == "" {
		// return the error
		return errors.New("Value is required.")
	}
	// nothing was wrong
	return nil
}

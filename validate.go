package probe

import (
	"errors"
)

func NonNull(str string) error {
	// if the string is empty
	if str == "" {
		// return the error
		return errors.New("Empty values not accepted.")
	}
	// nothing was wrong
	return nil
}

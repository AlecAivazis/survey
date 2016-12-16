package survey

import (
	"errors"
	"fmt"
)

// Required does not allow an empty value
func Required(str string) error {
	// if the string is empty
	if str == "" {
		// return the error
		return errors.New("Value is required.")
	}
	// nothing was wrong
	return nil
}

// MaxLength requires that the string is no longer than the specified value
func MaxLength(length int) Validator {
	// return a validator that checks the length of the string
	return func(str string) error {
		// if the string is longer than the given value
		if len(str) > length {
			return fmt.Errorf("Value is too long. Max length is %v.", length)
		}
		// the input is fine
		return nil
	}
}

// MinLength requires that the string is longer or equal in length to the specified value
func MinLength(length int) Validator {
	// return a validator that checks the length of the string
	return func(str string) error {
		// if the string is longer than the given value
		if len(str) < length {
			return fmt.Errorf("Value is too short. Min length is %v.", length)
		}
		// the input is fine
		return nil
	}
}

// ComposeValidators is a variadic function used to create one validator from many.
func ComposeValidators(validators ...Validator) Validator {
	// return a validator that calls each one sequentially
	return func(str string) error {
		// execute each validator
		for _, validator := range validators {
			// if the string is not valid
			if err := validator(str); err != nil {
				// return the error
				return err
			}
		}
		// we passed all validators, the string is valid
		return nil
	}
}

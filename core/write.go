package core

import (
	"errors"
	"reflect"
)

// Write takes a value and copies it to the target
func Write(t interface{}, v interface{}) (err error) {
	// when we're done
	defer func() {
		// if there was a panic
		if r := recover(); r != nil {
			// pass the panic on as an error
			err = r.(error)
		}
	}()

	// the target to write to
	target := reflect.ValueOf(t)
	value := reflect.ValueOf(v)

	// make sure we were handed a point
	if target.Kind() != reflect.Ptr {
		return errors.New("you must pass a pointer as the target of a Write operation")
	}

	// set the value of the object we're pointing to to match the value
	target.Elem().SetBool(value.Bool())

	// nothing went wrong
	return nil
}

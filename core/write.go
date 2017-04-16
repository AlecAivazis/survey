package core

import (
	"errors"
	"reflect"
)

// Write takes a value and copies it to the target
func Write(t interface{}, v interface{}) (err error) {
	// the target to write to
	target := reflect.ValueOf(t)
	value := reflect.ValueOf(v)

	// make sure we were handed a point
	if target.Kind() != reflect.Ptr {
		return errors.New("you must pass a pointer as the target of a Write operation")
	}

	// handle the target based on its prop
	switch target.Elem().Kind() {
	// if we are writing to a string
	case reflect.String:
		err = writeString(target.Elem(), value)
	// if we are writing to a bool
	case reflect.Bool:
		err = writeBool(target.Elem(), value)
	}

	// we're done
	return err
}

func writeBool(target, source reflect.Value) (err error) {
	// make sure we handle the source type
	switch source.Kind() {
	// if we are turning a boolean into a boolean
	case reflect.Bool:
		// just copy the boolean over
		target.SetBool(source.Bool())
	// otherwise its a source we do not recognize
	default:
		err = errors.New("Cannot convert to bool")
	}
	// nothing went wrong
	return err
}

func writeString(target, source reflect.Value) (err error) {
	// make sure we handle the source type
	switch source.Kind() {
	// if we are turning a string into a string
	case reflect.String:
		// just copy the string over
		target.SetString(source.String())
	// otherwise its a source we do not recognize
	default:
		err = errors.New("Cannot convert to string")
	}

	// nothing went wrong
	return err
}

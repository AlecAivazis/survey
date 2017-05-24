package core

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// the tag used to denote the name of the question
const tagName = "survey"

func WriteAnswer(t interface{}, name string, v interface{}) (err error) {
	// the target to write to
	target := reflect.ValueOf(t)
	// the value to write from
	value := reflect.ValueOf(v)

	// make sure we are writing to a pointer
	if target.Kind() != reflect.Ptr {
		return errors.New("you must pass a pointer as the target of a Write operation")
	}
	// the object "inside" of the target pointer
	elem := target.Elem()

	// handle the special types
	switch elem.Kind() {
	// if we are writing to a struct
	case reflect.Struct:
		// get the name of the field that matches the string we  were given
		fieldIndex, err := findFieldIndex(elem, name)
		// if something went wrong
		if err != nil {
			// bubble up
			return err
		}

		// copy the value over to the field
		return copy(elem.Field(fieldIndex), value)
	case reflect.Map:
		mapType := reflect.TypeOf(t).Elem()
		if mapType.Key().Kind() != reflect.String || mapType.Elem().Kind() != reflect.Interface {
			return errors.New("answer maps must be of type map[string]interface")
		}
		mt := *t.(*map[string]interface{})
		mt[name] = value.Interface()
		return nil
	}
	// otherwise just copy the value to the target
	return copy(elem, value)
}

// BUG(AlecAivazis): the current implementation might cause weird conflicts if there are
// two fields with same name that only differ by casing.
func findFieldIndex(s reflect.Value, name string) (int, error) {
	// the type of the value
	sType := s.Type()

	// first look for matching tags so we can overwrite matching field names
	for i := 0; i < sType.NumField(); i++ {
		// the field we are current scanning
		field := sType.Field(i)

		// the value of the survey tag
		tag := field.Tag.Get(tagName)
		// if the tag matches the name we are looking for
		if tag != "" && tag == name {
			// then we found our index
			return i, nil
		}
	}

	// then look for matching names
	for i := 0; i < sType.NumField(); i++ {
		// the field we are current scanning
		field := sType.Field(i)

		// if the name of the field matches what we're looking for
		if strings.ToLower(field.Name) == strings.ToLower(name) {
			return i, nil
		}
	}

	// we didn't find the field
	return -1, fmt.Errorf("could not find field matching %v", name)
}

// Write takes a value and copies it to the target
func copy(t reflect.Value, v reflect.Value) (err error) {
	// if something ends up panicing we need to catch it in a deferred func
	defer func() {
		if r := recover(); r != nil {
			// if we paniced with an error
			if _, ok := r.(error); ok {
				// cast the result to an error object
				err = r.(error)
			} else if _, ok := r.(string); ok {
				// otherwise we could have paniced with a string so wrap it in an error
				err = errors.New(r.(string))
			}
		}
	}()

	// attempt to copy the underlying value to the target
	t.Set(v)

	// we're done
	return
}

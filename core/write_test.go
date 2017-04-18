package core

import (
	"reflect"
	"testing"
)

func TestWrite_returnsErrorIfTargetNotPtr(t *testing.T) {
	// try to copy a value to a non-pointer
	err := WriteAnswer(true, "hello", true)
	// make sure there was an error
	if err == nil {
		t.Error("Did not encounter error when writing to non-pointer.")
	}
}

func TestWrite_canWriteToBool(t *testing.T) {
	// a pointer to hold the boolean value
	ptr := true

	// try to copy a false value to the pointer
	WriteAnswer(&ptr, "", false)

	// if the value is true
	if ptr {
		// the test failed
		t.Error("Could not write a false bool to a pointer")
	}
}

func TestWrite_canWriteString(t *testing.T) {
	// a pointer to hold the boolean value
	ptr := ""

	// try to copy a false value to the pointer
	err := WriteAnswer(&ptr, "", "hello")
	if err != nil {
		t.Error(err)
	}

	// if the value is not what we wrote
	if ptr != "hello" {
		t.Error("Could not write a string value to a pointer")
	}
}

func TestWrite_canWriteSlice(t *testing.T) {
	// a pointer to hold the value
	ptr := []string{}

	// copy in a value
	WriteAnswer(&ptr, "", []string{"hello", "world"})

	// make sure there are two entries
	if len(ptr) != 2 {
		// the test failed
		t.Errorf("Incorrect number of entries in written list. Expected 2, found %v.", len(ptr))
		// dont move on
		return
	}

	// make sure the first entry is hello
	if ptr[0] != "hello" {
		// the test failed
		t.Errorf("incorrect first value in written pointer. expected hello found %v.", ptr[0])
	}

	// make sure the second entry is world
	if ptr[1] != "world" {
		// the test failed
		t.Errorf("incorrect second value in written pointer. expected world found %v.", ptr[0])
	}
}

func TestWrite_recoversInvalidReflection(t *testing.T) {
	// a variable to mutate
	ptr := false

	// write a boolean value to the string
	err := WriteAnswer(&ptr, "", "hello")

	// if there was no error
	if err == nil {
		// the test failed
		t.Error("Did not encounter error when forced invalid write.")
	}
}

func TestWriteAnswer_handlesNonStructValues(t *testing.T) {
	// the value to write to
	ptr := ""

	// write a value to the pointer
	WriteAnswer(&ptr, "", "world")

	// if we didn't change the value appropriate
	if ptr != "world" {
		// the test failed
		t.Error("Did not write value to primitive pointer")
	}
}

func TestWriteAnswer_canMutateStruct(t *testing.T) {
	// the struct to hold the answer
	ptr := struct{ Name string }{}

	// write a value to an existing field
	err := WriteAnswer(&ptr, "name", "world")
	if err != nil {
		// the test failed
		t.Errorf("Encountered error while writing answer: %v", err.Error())
		// we're done here
		return
	}

	// make sure we changed the field
	if ptr.Name != "world" {
		// the test failed
		t.Error("Did not mutate struct field when writing answer.")
	}
}

func TestWriteAnswer_returnsErrWhenFieldNotFound(t *testing.T) {
	// the struct to hold the answer
	ptr := struct{ Name string }{}

	// write a value to an existing field
	err := WriteAnswer(&ptr, "", "world")

	if err == nil {
		// the test failed
		t.Error("Did not encountered error while writing answer to non-existing field.")
	}
}

func TestFindFieldIndex_canFindExportedField(t *testing.T) {
	// create a reflective wrapper over the struct to look through
	val := reflect.ValueOf(struct{ Name string }{})

	// find the field matching "name"
	fieldIndex, err := findFieldIndex(val, "name")
	// if something went wrong
	if err != nil {
		// the test failed
		t.Error(err.Error())
		return
	}

	// make sure we got the right value
	if val.Type().Field(fieldIndex).Name != "Name" {
		// the test failed
		t.Errorf("Did not find the correct field name. Expected 'Name' found %v.", val.Type().Field(fieldIndex).Name)
	}
}

func TestFindFieldIndex_canFindTaggedField(t *testing.T) {
	// the struct to look through
	val := reflect.ValueOf(struct {
		Username string `survey:"name"`
	}{})

	// find the field matching "name"
	fieldIndex, err := findFieldIndex(val, "name")
	// if something went wrong
	if err != nil {
		// the test failed
		t.Error(err.Error())
		return
	}

	// make sure we got the right value
	if val.Type().Field(fieldIndex).Name != "Username" {
		// the test failed
		t.Errorf("Did not find the correct field name. Expected 'Username' found %v.", val.Type().Field(fieldIndex).Name)
	}
}

func TestFindFieldIndex_canHandleCapitalAnswerNames(t *testing.T) {
	// create a reflective wrapper over the struct to look through
	val := reflect.ValueOf(struct{ Name string }{})

	// find the field matching "name"
	fieldIndex, err := findFieldIndex(val, "Name")
	// if something went wrong
	if err != nil {
		// the test failed
		t.Error(err.Error())
		return
	}

	// make sure we got the right value
	if val.Type().Field(fieldIndex).Name != "Name" {
		// the test failed
		t.Errorf("Did not find the correct field name. Expected 'Name' found %v.", val.Type().Field(fieldIndex).Name)
	}
}

func TestFindFieldIndex_tagOverwriteFieldName(t *testing.T) {
	// the struct to look through
	val := reflect.ValueOf(struct {
		Name     string
		Username string `survey:"name"`
	}{})

	// find the field matching "name"
	fieldIndex, err := findFieldIndex(val, "name")
	// if something went wrong
	if err != nil {
		// the test failed
		t.Error(err.Error())
		return
	}

	// make sure we got the right value
	if val.Type().Field(fieldIndex).Name != "Username" {
		// the test failed
		t.Errorf("Did not find the correct field name. Expected 'Username' found %v.", val.Type().Field(fieldIndex).Name)
	}
}

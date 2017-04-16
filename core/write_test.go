package core

import (
	"testing"
)

func TestWrite_returnsErrorIfTargetNotPtr(t *testing.T) {
	// try to copy a value to a non-pointer
	err := Write(true, true)
	// make sure there was an error
	if err == nil {
		t.Error("Did not encounter error when writing to non-pointer.")
	}
}

func TestWrite_canWriteToBool(t *testing.T) {
	// a pointer to hold the boolean value
	ptr := true

	// try to copy a false value to the pointer
	Write(&ptr, false)

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
	err := Write(&ptr, "hello")
	if err != nil {
		t.Error(err)
	}

	// if the value is not what we wrote
	if ptr != "hello" {
		t.Error("Could not write a string value to a pointer")
	}
}

func TestWrite_gracefullyHandlesFailedStringWrites(t *testing.T) {
	// a pointer to hold the boolean value
	ptr := ""
	// try to copy a false value to the pointer
	err := Write(&ptr, false)
	// if the value is try
	if err == nil {
		// the test failed
		t.Error("Did not encouner error when casting boolean to string")
	}
}

// func TestWrite_recoversInvalidReflection(t *testing.T) {
// 	// a variable to mutate
// 	ptr := false

// 	// write a boolean value to the string
// 	err := Write(&ptr, "hello")
// 	fmt.Println(err.Error())
// 	// if there was no error
// 	if err == nil {
// 		// the test failed
// 		t.Error("Did not encounter error when forced invalid write.")
// 	}
// }

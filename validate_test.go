package probe

import "testing"

func TestNonNullCanSucceed(t *testing.T) {
	// a string to test
	str := "hello"
	// if the string is not valid
	if valid := NonNull(str); valid != nil {
		//
		t.Error("Non null returned an error when one wasn't expected.")
	}
}
func TestNonNullCanFail(t *testing.T) {
	// a string to test
	str := ""
	// if the string is valid
	if notValid := NonNull(str); notValid == nil {
		//
		t.Error("Non null did not return an error when one was expected.")
	}
}

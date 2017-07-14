package survey

import (
	"errors"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"testing"

	"github.com/AlecAivazis/survey/core"
	"github.com/stretchr/testify/assert"
)

func init() {
	// disable color output for all prompts to simplify testing
	core.DisableColor = true
}

func TestConverter(t *testing.T) {
	q := []*Question{
		{
			Name: "age",
			Prompt: &Input{
				Message: "What is your age?",
			},
			Validate: Required,
			Convert: func(val interface{}) (interface{}, error) {
				newVal, err := strconv.ParseInt(val.(string), 10, 64)
				if err != nil {
					return nil, errors.New("please enter a valid integer")
				}

				return newVal, nil
			},
		},
	}

	in, _ := ioutil.TempFile("", "")
	defer in.Close()

	os.Stdin = in

	io.WriteString(in, "21\n")
	in.Seek(0, os.SEEK_SET)

	age := int64(0)
	err := Ask(q, &age)

	assert.Nil(t, err)
	assert.Equal(t, int64(21), age)

}

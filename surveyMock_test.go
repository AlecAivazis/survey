package survey

import (
	"testing"
)

func TestMockAskOne(t *testing.T) {
	t.Run("Setting Response", func(t *testing.T) {
		mock := SurveyMock{}
		mock.SetResponse(true)

		prompt := &Confirm{
			Message: "test",
		}

		var response bool
		mock.AskOne(prompt, &response)

		if !response {
			t.Fatalf("Response was false but should have been true!")
		}
	})

}

func TestMockAsk(t *testing.T) {
	t.Run("Setting Response", func(t *testing.T) {
		mock := SurveyMock{}

		test := make(map[string]interface{})
		test["test"] = true

		mock.SetResponse(test)

		questions := []*Question{
			{
				Name: "test",
				Prompt: &Confirm{
					Message: "testing",
				},
			},
		}

		answer := struct {
			Test bool
		}{}
		mock.Ask(questions, &answer)

		if !answer.Test {
			t.Fatalf("Response was false but should have been true!")
		}
	})
}

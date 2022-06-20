package survey

import (
	"io"

	"github.com/AlecAivazis/survey/v2/core"
	"github.com/AlecAivazis/survey/v2/terminal"
)

type SurveyMock struct {
	singleResponse    interface{}
	multipleResponses map[string]interface{}
}

// Not implemented, because it does not affect the mock
func (mock *SurveyMock) WithStdio(in terminal.FileReader, out terminal.FileWriter, err io.Writer) AskOpt {
	return nil
}

// Not implemented, because it does not affect the mock
func (mock *SurveyMock) WithFilter(filter func(filter string, value string, index int) (include bool)) AskOpt {
	return nil
}

// Not implemented, because it does not affect the mock
func (mock *SurveyMock) WithKeepFilter(KeepFilter bool) AskOpt {
	return nil
}

// Not implemented, because it does not affect the mock
func (mock *SurveyMock) WithValidator(v Validator) AskOpt {
	return nil
}

// Not implemented, because it does not affect the mock
func (mock *SurveyMock) WithPageSize(pageSize int) AskOpt {
	return nil
}

// Not implemented, because it does not affect the mock
func (mock *SurveyMock) WithHelpInput(r rune) AskOpt {
	return nil
}

// Not implemented, because it does not affect the mock
func (mock *SurveyMock) WithIcons(setIcons func(*IconSet)) AskOpt {
	return nil
}

// Not implemented, because it does not affect the mock
func (mock *SurveyMock) WithShowCursor(ShowCursor bool) AskOpt {
	return nil
}

func (mock *SurveyMock) AskOne(p Prompt, response interface{}, opts ...AskOpt) error {
	err := core.WriteAnswer(response, "", mock.singleResponse)
	if err != nil {
		// panicing is fine inside a mock
		panic(err)
	}
	return nil
}

func (mock *SurveyMock) Ask(qs []*Question, response interface{}, opts ...AskOpt) error {
	for _, q := range qs {

		err := core.WriteAnswer(response, q.Name, mock.multipleResponses[q.Name])
		if err != nil {
			// panicing is fine inside a mock
			panic(err)
		}

	}

	return nil
}

func (mock *SurveyMock) SetResponse(response interface{}) {
	if val, ok := response.(map[string]interface{}); ok {
		mock.multipleResponses = val
	} else {
		mock.singleResponse = response
	}
}

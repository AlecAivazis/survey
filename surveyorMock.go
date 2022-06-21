package survey

import (
	"github.com/AlecAivazis/survey/v2/core"
)

type SurveyorMock struct {
	singleResponse    interface{}
	multipleResponses map[string]interface{}
}

func (mock *SurveyorMock) AskOne(p Prompt, response interface{}, opts ...AskOpt) error {
	err := core.WriteAnswer(response, "", mock.singleResponse)
	if err != nil {
		// panicing is fine inside a mock
		panic(err)
	}
	return nil
}

func (mock *SurveyorMock) Ask(qs []*Question, response interface{}, opts ...AskOpt) error {
	for _, q := range qs {

		err := core.WriteAnswer(response, q.Name, mock.multipleResponses[q.Name])
		if err != nil {
			// panicing is fine inside a mock
			panic(err)
		}

	}

	return nil
}

func (mock *SurveyorMock) SetResponse(response interface{}) {
	if val, ok := response.(map[string]interface{}); ok {
		mock.multipleResponses = val
	} else {
		mock.singleResponse = response
	}
}

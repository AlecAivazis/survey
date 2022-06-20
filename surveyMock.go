package survey

import (
	"io"

	"github.com/AlecAivazis/survey/v2/terminal"
)

type SurveyMock struct{}

//TODO implement the mock functionality
func (survey *SurveyMock) WithStdio(in terminal.FileReader, out terminal.FileWriter, err io.Writer) AskOpt {
	return nil
}
func (survey *SurveyMock) WithFilter(filter func(filter string, value string, index int) (include bool)) AskOpt {
	return nil
}
func (survey *SurveyMock) WithKeepFilter(KeepFilter bool) AskOpt {
	return nil
}
func (survey *SurveyMock) WithValidator(v Validator) AskOpt {
	return nil
}
func (survey *SurveyMock) WithPageSize(pageSize int) AskOpt {
	return nil
}
func (survey *SurveyMock) WithHelpInput(r rune) AskOpt {
	return nil
}
func (survey *SurveyMock) WithIcons(setIcons func(*IconSet)) AskOpt {
	return nil
}
func (survey *SurveyMock) WithShowCursor(ShowCursor bool) AskOpt {
	return nil
}
func (survey *SurveyMock) AskOne(p Prompt, response interface{}, opts ...AskOpt) error {
	return nil
}
func (survey *SurveyMock) Ask(qs []*Question, response interface{}, opts ...AskOpt) error {
	return nil
}

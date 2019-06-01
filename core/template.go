package core

import (
	"bytes"
	"sync"
	"text/template"

	"github.com/mgutz/ansi"
)

// DisableColor can be used to make testing reliable
var DisableColor = false

// IconSet holds the strings to use for various prompts
type IconSet struct {
	HelpInput      rune
	Error          string
	Help           string
	Question       string
	MarkedOption   string
	UnmarkedOption string
	SelectFocus    string
}

// DefaultIconSet is the default icons used by prompts
var DefaultIconSet = IconSet{
	HelpInput:      '?',
	Error:          "X",
	Help:           "?",
	Question:       "?",
	MarkedOption:   "[x]",
	UnmarkedOption: "[ ]",
	SelectFocus:    ">",
}

var TemplateFuncs = map[string]interface{}{
	// Templates with Color formatting. See Documentation: https://github.com/mgutz/ansi#style-format
	"color": func(color string) string {
		if DisableColor {
			return ""
		}
		return ansi.ColorCode(color)
	},
	"HelpInputRune": func() string {
		return string(DefaultIconSet.HelpInput)
	},
	"ErrorIcon": func() string {
		return DefaultIconSet.Error
	},
	"HelpIcon": func() string {
		return DefaultIconSet.Help
	},
	"QuestionIcon": func() string {
		return DefaultIconSet.Question
	},
	"MarkedOptionIcon": func() string {
		return DefaultIconSet.MarkedOption
	},
	"UnmarkedOptionIcon": func() string {
		return DefaultIconSet.UnmarkedOption
	},
	"SelectFocusIcon": func() string {
		return DefaultIconSet.SelectFocus
	},
}

var (
	memoizedGetTemplate = map[string]*template.Template{}

	memoMutex = &sync.RWMutex{}
)

func getTemplate(tmpl string) (*template.Template, error) {
	memoMutex.RLock()
	if t, ok := memoizedGetTemplate[tmpl]; ok {
		memoMutex.RUnlock()
		return t, nil
	}
	memoMutex.RUnlock()

	t, err := template.New("prompt").Funcs(TemplateFuncs).Parse(tmpl)
	if err != nil {
		return nil, err
	}

	memoMutex.Lock()
	memoizedGetTemplate[tmpl] = t
	memoMutex.Unlock()
	return t, nil
}

func RunTemplate(tmpl string, data interface{}) (string, error) {
	t, err := getTemplate(tmpl)
	if err != nil {
		return "", err
	}
	buf := bytes.NewBufferString("")
	err = t.Execute(buf, data)
	if err != nil {
		return "", err
	}
	return buf.String(), err
}

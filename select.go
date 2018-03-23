package survey

import (
	"gopkg.in/AlecAivazis/survey.v1/core"
	"gopkg.in/AlecAivazis/survey.v1/terminal"
	"math"
	"strings"
	"os"
)

/*
Select is a prompt that presents a list of various options to the user
for them to select using the arrow keys and enter. Response type is a string.

	color := &survey.Option{}
	prompt := survey.NewSingleSelect().SetMessage("Select Color:").
			AddOption("red", nil, false).
			AddOption("blue", nil, false).
			AddOption("green", nil, false)
	survey.AskOne(prompt, &color, nil)
*/
type Select struct {
	core.Renderer
	Message       string
	Options       Options
	Default       *Option
	Help          string
	PageSize      int
	VimMode       bool
	FilterMessage string
	filter        string
	selectedIndex int
	useDefault    bool
	showingHelp   bool
}

/*
NewSingleSelect is a shortcut method to get a Select prompt
 */
func NewSingleSelect() *Select {
	return &Select{
		Options: make(Options, 0),
	}
}

/*
AddOption is a method to add an option to the selection and specify if it is the default value ot not
This returns a Selection interface to allow chaining of these method calls

	color := &survey.Option{}
	prompt := survey.NewSingleSelect().SetMessage("Select Color:").
			AddOption("red", nil, false).
			AddOption("blue", nil, false).
			AddOption("green", nil, false)
	survey.AskOne(prompt, &color, nil)
 */
func (s *Select) AddOption(display string, value interface{}, defaultOption bool) Selection {
	if value == nil {
		value = display
	}

	opt := &Option{display, value}
	s.Options = append(s.Options, opt)
	if defaultOption {
		s.Default = opt
	}
	return s
}

/*
SetMessage is a method to set the prompt message for a selection
This returns a Selection interface to allow chaining of these method calls
 */
func (s *Select) SetMessage(msg string) Selection {
	s.Message = msg
	return s
}

/*
SetHelp is a method to set the prompt help message for a selection
This returns a Selection interface to allow chaining of these method calls
 */
func (s *Select) SetHelp(help string) Selection {
	s.Help = help
	return s
}

/*
SetFilterMessage is a method to set the prompt filter message for a selection
This returns a Selection interface to allow chaining of these method calls
 */
func (s *Select) SetFilterMessage(msg string) Selection {
	s.FilterMessage = msg
	return s
}

/*
SetVimMode is a method to turn on or off VimMode
This returns a Selection interface to allow chaining of these method calls
 */
func (s *Select) SetVimMode(vimMode bool) Selection {
	s.VimMode = vimMode
	return s
}

/*
SetVimMode is a method to turn on or off VimMode
This returns a Selection interface to allow chaining of these method calls
 */
func (s *Select) SetPageSize(pageSize int) Selection {
	s.PageSize = pageSize
	return s
}

// Paginate returns a single page of choices given the page size, the total list of
// possible choices, and the current selected index in the total list.
func (s *Select) Paginate(choices Options) (Options, int) {
	if s.PageSize == 0 {
		s.PageSize = PageSize
	}

	var start, end, max, cursor int
	max = len(choices)
	if max < s.PageSize {
		// if we dont have enough options to fill a page
		start = 0
		end = max
		cursor = s.selectedIndex
	} else if s.selectedIndex < s.PageSize/2 {
		// if we are in the first half page
		start = 0
		end = s.PageSize
		cursor = s.selectedIndex
	} else if max-s.selectedIndex-1 < s.PageSize/2 {
		// if we are in the last half page
		start = max - s.PageSize
		end = max
		cursor = s.selectedIndex - start
	} else {
		// somewhere in the middle
		above := s.PageSize / 2
		below := s.PageSize - above

		cursor = s.PageSize / 2
		start = s.selectedIndex - above
		end = s.selectedIndex + below
	}
	end = int(math.Min(float64(end), float64(max)))

	// return the subset we care about and the index
	return choices[start:end], cursor
}


// the data available to the templates when processing
type SelectTemplateData struct {
	Select
	PageEntries   Options
	SelectedIndex int
	Answer        *Option
	ShowAnswer    bool
	ShowHelp      bool
}

var SelectQuestionTemplate = `
{{- if .ShowHelp }}{{- color "cyan"}}{{ HelpIcon }} {{ .Help }}{{color "reset"}}{{"\n"}}{{end}}
{{- color "green+hb"}}{{ QuestionIcon }} {{color "reset"}}
{{- color "default+hb"}}{{ .Message }}{{ .FilterMessage }}{{color "reset"}}
{{- if .ShowAnswer}} {{ color "cyan"}}{{.Answer.Display}}{{color "reset"}}{{"\n"}}
{{- else}}
  {{- "  "}}{{- color "cyan"}}[Use arrows to move, type to filter{{- if and .Help (not .ShowHelp)}}, {{ HelpInputRune }} for more help{{end}}]{{color "reset"}}
  {{- "\n"}}
  {{- range $ix, $choice := .PageEntries}}
    {{- if eq $ix $.SelectedIndex}}{{color "cyan+b"}}{{ SelectFocusIcon }} {{else}}{{color "default+hb"}}  {{end}}
    {{- $choice.Display }}
    {{- color "reset"}}{{"\n"}}
  {{- end}}
{{- end}}`

// OnChange is called on every keypress.
func (s *Select) OnChange(line []rune, pos int, key rune) (newLine []rune, newPos int, ok bool) {
	options := s.filterOptions()
	oldFilter := s.filter

	// if the user pressed the enter key
	if key == terminal.KeyEnter {
		if s.selectedIndex < len(options) {
			return []rune(options[s.selectedIndex].Display), 0, true
		}
		// if the user pressed the up arrow or 'k' to emulate vim
	} else if key == terminal.KeyArrowUp || (s.VimMode && key == 'k') {
		s.useDefault = false

		// if we are at the top of the list
		if s.selectedIndex == 0 {
			// start from the button
			s.selectedIndex = len(options) - 1
		} else {
			// otherwise we are not at the top of the list so decrement the selected index
			s.selectedIndex--
		}
		// if the user pressed down or 'j' to emulate vim
	} else if key == terminal.KeyArrowDown || (s.VimMode && key == 'j') {
		s.useDefault = false
		// if we are at the bottom of the list
		if s.selectedIndex == len(options)-1 {
			// start from the top
			s.selectedIndex = 0
		} else {
			// increment the selected index
			s.selectedIndex++
		}
		// only show the help message if we have one
	} else if key == core.HelpInputRune && s.Help != "" {
		s.showingHelp = true
	} else if key == terminal.KeyEscape {
		s.VimMode = !s.VimMode
	} else if key == terminal.KeyDeleteWord || key == terminal.KeyDeleteLine {
		s.filter = ""
	} else if key == terminal.KeyDelete || key == terminal.KeyBackspace {
		if s.filter != "" {
			s.filter = s.filter[0 : len(s.filter)-1]
		}
	} else if key >= terminal.KeySpace {
		s.filter += string(key)
		s.VimMode = false
	}

	s.FilterMessage = ""
	if s.filter != "" {
		s.FilterMessage = " " + s.filter
	}
	if oldFilter != s.filter {
		// filter changed
		options = s.filterOptions()
		if len(options) > 0 && len(options) <= s.selectedIndex {
			s.selectedIndex = len(options) - 1
		}
	}

	// figure out the options and index to render

	// TODO if we have started filtering and were looking at the end of a list
	// and we have modified the filter then we should move the page back!
	opts, idx := s.Paginate(options)

	// render the options
	s.Render(
		SelectQuestionTemplate,
		SelectTemplateData{
			Select:        *s,
			SelectedIndex: idx,
			ShowHelp:      s.showingHelp,
			PageEntries:   opts,
		},
	)

	// if we are not pressing ent
	if len(options) <= s.selectedIndex {
		return []rune{}, 0, false
	}
	return []rune(options[s.selectedIndex].Display), 0, true
}

func (s *Select) filterOptions() Options {
	filter := strings.ToLower(s.filter)
	if filter == "" {
		return s.Options
	}
	answer := make(Options, 0)
	for _, o := range s.Options {
		if strings.Contains(strings.ToLower(o.Display), filter) {
			answer = append(answer, o)
		}
	}
	return answer
}

func (s *Select) Prompt() (interface{}, error) {
	// if there are no options to render
	if len(s.Options) == 0 {
		// we failed
		return "", errors.New("please provide options to select from")
	}

	// start off with the first option selected
	sel := 0
	// if there is a default
	if s.Default != nil {
		// find the choice
		for i, opt := range s.Options {
			// if the option correponds to the default
			if opt == s.Default {
				// we found our initial value
				sel = i
				// stop looking
				break
			}
		}
	}
	// save the selected index
	s.selectedIndex = sel

	// figure out the options and index to render
	opts, idx := s.Paginate(s.Options)

	// ask the question
	err := s.Render(
		SelectQuestionTemplate,
		SelectTemplateData{
			Select:        *s,
			PageEntries:   opts,
			SelectedIndex: idx,
		},
	)
	if err != nil {
		return "", err
	}

	// hide the cursor
	terminal.CursorHide()
	// show the cursor when we're done
	defer terminal.CursorShow()

	// by default, use the default value
	s.useDefault = true

	rr := terminal.NewRuneReader(os.Stdin)
	rr.SetTermMode()
	defer rr.RestoreTermMode()
	// start waiting for input
	for {
		r, _, err := rr.ReadRune()
		if err != nil {
			return "", err
		}
		if r == '\r' || r == '\n' {
			break
		}
		if r == terminal.KeyInterrupt {
			return "", terminal.InterruptErr
		}
		if r == terminal.KeyEndTransmission {
			break
		}
		s.OnChange(nil, 0, r)
	}
	options := s.filterOptions()
	s.filter = ""
	s.FilterMessage = ""

	var val *Option
	// if we are supposed to use the default value
	if s.useDefault || s.selectedIndex >= len(options) {
		// if there is a default value
		if s.Default != nil {
			// use the default value
			val = s.Default
		} else if len(options) > 0 {
			// there is no default value so use the first
			val = options[0]
		}
		// otherwise the selected index points to the value
	} else if s.selectedIndex < len(options) {
		val = options[s.selectedIndex]
	}
	return val, err
}

func (s *Select) Cleanup(val interface{}) error {
	return s.Render(
		SelectQuestionTemplate,
		SelectTemplateData{
			Select:     *s,
			Answer:     val.(*Option),
			ShowAnswer: true,
		},
	)
}

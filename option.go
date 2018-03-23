package survey

// Selection is the interface for a prompt that will present a fixed section to the user
type Selection interface {
	Prompt
	AddOption(display string, value interface{}, defaultOption bool) Selection
	SetMessage(msg string) Selection
	SetHelp(help string) Selection
	SetFilterMessage(msg string) Selection
	SetVimMode(vimMode bool) Selection
	SetPageSize(pageSize int) Selection
	Paginate(choices Options) (Options, int)
	OnChange(line []rune, pos int, key rune) (newLine []rune, newPos int, ok bool)
}

// Option is a struct to selerate the display value from the actual value of a selection
type Option struct {
	Display string
	Value interface{}
}

// String to define the default output of Option
func (o *Option) String() string {
	return o.Display
}

// Options alias for []*Option
type Options = []*Option

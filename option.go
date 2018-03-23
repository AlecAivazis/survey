package survey

// Selection is the interface for a prompt that will present a section to the user
type Selection interface {
	Prompt
	AddOption(display string, value interface{}, defaultOption bool) Selection
	SetMessage(msg string) Selection
	SetHelp(help string) Selection
	SetFilterMessage(msg string) Selection
	SetVimMode(vim bool) Selection
	SetPageSize(size int) Selection
	Paginate(choices Options) (Options, int)
	OnChange(line []rune, pos int, key rune) (newLine []rune, newPos int, ok bool)
}

type Option struct {
	Display string
	Value interface{}
}

func (o *Option) String() string {
	return o.Display
}

type Options = []*Option

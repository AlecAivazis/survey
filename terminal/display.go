package terminal

type EraseLineMode int

const (
	ERASE_LINE_END   EraseLineMode = iota
	ERASE_LINE_START EraseLineMode = iota
	ERASE_LINE_ALL   EraseLineMode = iota
)

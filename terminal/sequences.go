package terminal

const (
	KeyArrowLeft       = '\x02'
	KeyArrowRight      = '\x06'
	KeyArrowUp         = '\x10'
	KeyArrowDown       = '\x0e'
	KeySpace           = ' '
	KeyEnter           = '\r'
	KeyBackspace       = '\b'
	KeyDelete          = '\x7f'
	KeyInterrupt       = '\x03'
	KeyEndTransmission = '\x04'
	KeyEscape		   = '\x1b'
	KeyDeleteWord      = '\x17' // Ctrl+W
	KeyDeleteLine      = '\x18' // Ctrl+X
)

func soundBell() {
	Print("\a")
}

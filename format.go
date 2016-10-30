package probe

import "fmt"

// FormatAsk prepares a string to be printed like the first line
// of a prompt
func FormatAsk(q string) string {
	return fmt.Sprintf("? %v", q)
}

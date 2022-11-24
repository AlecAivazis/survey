package log

import (
	"fmt"
	"log"
	"os"
	"strings"
)

var enabled bool

func init() {
	logFile := os.Getenv("SURVEY_LOG_FILE")
	if logFile == "" {
		return
	}

	w, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "survey: enable to open file %q for logging\n", logFile)
		return
	}

	log.SetOutput(w)
	log.Print(strings.Repeat("=", 80))
	enabled = true
}

func Print(s string) {
	if !enabled {
		return
	}
	log.Print(s)
}

func Printf(f string, args ...interface{}) {
	if !enabled {
		return
	}
	log.Printf(f, args...)
}

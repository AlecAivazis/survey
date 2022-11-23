package log

import (
	"log"
	"os"
)

func init() {
	w, err := os.OpenFile("survey.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	log.SetOutput(w)
}

func Print(s string) {
	log.Print(s)
}

func Printf(f string, args ...interface{}) {
	log.Printf(f, args...)
}

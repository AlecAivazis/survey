package main

import (
// "fmt"
// tm "github.com/buger/goterm"
)

type Question struct {
	Name   string
	Prompt Prompt
}

type Prompt interface {
	Prompt() (string, error)
	AskQuestion()
}

func Ask(qs []*Question) map[string]string {
	// the response map
	res := make(map[string]string)
	// go over every question
	for _, q := range qs {
		// print the question
		q.Prompt.AskQuestion()
		// grab the user input and save it
		ans, err := q.Prompt.Prompt()
		// if something went wrong
		if err != nil {
			panic(err)
		}
		// add it to the map
		res[q.Name] = ans
	}
	// return the response
	return res
}

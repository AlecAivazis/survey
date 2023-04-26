//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"path/filepath"
)

func suggestFiles(toComplete string) []string {
	files, _ := filepath.Glob(toComplete + "*")
	return files
}

// the questions to ask
var q = []*survey.Question{
	{
		Name: "file",
		Prompt: &survey.Input{
			Message: "Which file should be read?",
			Suggest: suggestFiles,
			Help:    "Any file; do not need to exist yet",
		},
		Validate: survey.Required,
	},
}

func main() {
	answers := struct {
		File string
	}{}

	// ask the question
	err := survey.Ask(q, &answers)

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	// print the answers
	fmt.Printf("File chosen %s.\n", answers.File)
}

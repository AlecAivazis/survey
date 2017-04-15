package main

import (
	"fmt"

	"github.com/alecaivazis/survey"
)

func main() {
	days := []string{}
	prompt := &survey.MultiChoice{
		Message:  "What days do you prefer:",
		Options:  []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"},
		Defaults: []string{"Saturday", "Sunday"},
		Answer:   &days,
	}

	answer, err := survey.AskOne(prompt)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Printf("response string (json for MultiChoice): %s\n", answer)
	fmt.Printf("response days: %#v\n", days)
}

package main

import (
	"fmt"

	"github.com/alecaivazis/survey"
)

func main() {
	days := []string{}
	prompt := &survey.MultiSelect{
		Message:  "What days do you prefer:",
		Options:  []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"},
		Defaults: []string{"Saturday", "Sunday"},
	}

	err := survey.AskOne(prompt, &days, nil)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Printf("response: %s\n", days)
}

package main

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
)

// the questions to ask
type Meal struct {
	Title   string
	Comment string
}

func main() {
	var meals = []Meal{
		{Title: "Bread", Comment: "Good, not so healthy"},
		{Title: "Eggs", Comment: "Healthy"},
		{Title: "Apple", Comment: ""},
		{Title: "Burger", Comment: "Really?"},
	}

	answers := struct {
		Meal string
	}{}

	titles := make([]string, len(meals))
	for i, m := range meals {
		titles[i] = m.Title
	}
	var qs = []*survey.Question{
		{
			Name: "meal",
			Prompt: &survey.Select{
				Message: "Choose a meal:",
				Options: titles,
				Description: func(value string, index int) string {
					return meals[index].Comment
				},
			},
			Validate: survey.Required,
		},
	}

	// ask the question
	err := survey.Ask(qs, &answers)

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	// print the answers
	fmt.Printf("you picked %s, nice choice.\n", answers.Meal)
}

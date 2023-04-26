//go:build ignore
// +build ignore

package main

import (
	"fmt"
)

type Meal struct {
	Title   string
	Comment string
}

func main() {
	var meals = []Meal{
		{Title: "Bread", Comment: "Contains gluten"},
		{Title: "Eggs", Comment: "Free-range"},
		{Title: "Apple", Comment: ""},
		{Title: "Burger", Comment: "Veggie patties available"},
	}

	titles := make([]string, len(meals))
	for i, m := range meals {
		titles[i] = m.Title
	}
	var qs = &survey.Select{
		Message: "Choose a meal:",
		Options: titles,
		Description: func(value string, index int) string {
			return meals[index].Comment
		},
	}

	answerIndex := 0

	// ask the question
	err := survey.AskOne(qs, &answerIndex)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	meal := meals[answerIndex]
	// print the answers
	fmt.Printf("you picked %s, nice choice.\n", meal.Title)
}

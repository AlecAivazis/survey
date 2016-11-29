# Survey
A library for building interactive prompts. Heavily inspired by the great [inquirer.js](https://github.com/SBoudrias/Inquirer.js/).

[![GoDoc](http://img.shields.io/badge/godoc-reference-5272B4.svg)](https://godoc.org/github.com/alecaivazis/survey)


![](https://zippy.gfycat.com/AmusingBossyArrowworm.gif)

```go
package main

import (
    "fmt"
    "github.com/alecaivazis/survey"
)

// the questions to ask
var qs = []*survey.Question{
    {
        Name:     "name",
        Prompt:   &survey.Input{"What is your name?", ""},
        Validate: survey.Required,
    },
    {
        Name: "color",
        Prompt: &survey.Choice{
            Message: "Choose a color:",
            Choices: []string{"red", "blue", "green"},
            Default: "red",
        },
    },
}

func main() {
    answers, err := survey.Ask(qs)

    if err != nil {
        fmt.Println("\n", err.Error())
        return
    }

    fmt.Printf("%s chose %s.", answers["name"], answers["color"])
}

```

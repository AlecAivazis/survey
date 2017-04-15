# Survey
[![Build Status](https://travis-ci.org/AlecAivazis/survey.svg?branch=feature%2Fpretty)](https://travis-ci.org/AlecAivazis/survey)
[![GoDoc](http://img.shields.io/badge/godoc-reference-5272B4.svg)](https://godoc.org/github.com/alecaivazis/survey)


A library for building interactive prompts. Heavily inspired by the great [inquirer.js](https://github.com/SBoudrias/Inquirer.js/).



![](https://zippy.gfycat.com/AmusingBossyArrowworm.gif)

```go
package main

import (
    "fmt"
    "gopkg.in/alecaivazis/survey.v0"
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
        Prompt: &survey.Select{
            Message: "Choose a color:",
            Options: []string{"red", "blue", "green"},
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

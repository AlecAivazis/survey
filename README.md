# Probe
A library for building interactive prompts. Heavily inspired by the great [inquirer.js](https://github.com/SBoudrias/Inquirer.js/).

```go
package main

import (
    "fmt"
    "github.com/alecaivazis/probe"
)

// the questions to ask
var qs = []*probe.Question{
    {
        Name:   "name",
        Prompt: &probe.Input{"What is your name?"},
        Validate: probe.NonNull,
    },
    {
        Name: "color",
        Prompt: &probe.Choice{
            Message: "Choose a color:",
            Choices:  []string{"red", "blue", "green"},
        },
    },
}

func main() {
    answers, err := probe.Ask(qs)

    if err != nil {
        fmt.Println("\n", err.Error())
        return
    }

    fmt.Printf("%s chose %s.", answers["name"], answers["color"])
}

```

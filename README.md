# Survey
A library for building interactive prompts. Heavily inspired by the great [inquirer.js](https://github.com/SBoudrias/Inquirer.js/).


![survey](https://zippy.gfycat.com/DisastrousDescriptiveGrunion.gif "survey")


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
        Prompt:   &survey.Input{"What is your name?", nil},
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

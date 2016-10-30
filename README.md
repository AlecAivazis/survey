# Probe
A library for building interactive prompts. Heavily inspired by the great [inquirer.js](https://github.com/SBoudrias/Inquirer.js/).

```golang
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
    },
    {
        Name:   "birthday",
        Prompt: &probe.Input{"When is your birthday?"},
    },
}

func main() {
    answers := probe.Ask(qs)

    fmt.Println(answers)
}
```

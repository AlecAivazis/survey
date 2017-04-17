# Survey
[![Build Status](https://travis-ci.org/AlecAivazis/survey.svg?branch=feature%2Fpretty)](https://travis-ci.org/AlecAivazis/survey)
[![GoDoc](http://img.shields.io/badge/godoc-reference-5272B4.svg)](https://godoc.org/github.com/alecaivazis/survey)


A library for building interactive prompts. Heavily inspired by the great [inquirer.js](https://github.com/SBoudrias/Inquirer.js/).



![](https://zippy.gfycat.com/AmusingBossyArrowworm.gif)

```go
package main

import (
    "fmt"
    "gopkg.in/alecaivazis/survey.v1"
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
    // the answers will be written to this struct
    answers := struct {
        Name          string                  // survey will match the question and field names
        FavoriteColor string `survey:"color"` // or you can tag fields to match a specific name
    }

    // perform the questions
    err := survey.Ask(qs, &answers)
    if err != nil {
        fmt.Println(err.Error())
        return
    }

    fmt.Printf("%s chose %s.", answers.Name, answers.FavoriteColor)
}
```

## Examples
Examples can be found in the `examples/` directory. Run them
to see basic behavior:
```bash
go get github.com/alecaivazis/survey

# ... navigate to the repo in your GOPATH

go run examples/simple.go
go run examples/validation.go
```

## Prompts

### Input
<img src="https://media.giphy.com/media/3og0IxS8JsuD9Z8syA/giphy.gif" width="300px"/>

```golang
name := ""
prompt = &survey.Input{
    Message: "What's your name?",
}
survey.AskOne(prompt, &name, nil)
```


### Password
<img src="https://media.giphy.com/media/3o7bu960gXMggMttXG/giphy.gif" width="300px" />

```golang
password := ""
prompt = &survey.Password{"Please type your password"}
survey.AskOne(prompt, &password, nil)
```


### Confirm
<img src="https://media.giphy.com/media/3og0IFvdDIaUgJzbcQ/giphy.gif" width="300px"/>

```golang
name := false
prompt = &survey.Confirm{
    Message: "Do you like pie?",
}
survey.AskOne(prompt, &name, nil)
```


### Select
<img src="https://media.giphy.com/media/l0IykKO3Vdxw3daGA/giphy.gif" width="300px"/>

```golang
color := ""
prompt = &survey.Select{
    Options: []string{"red", "blue" "green"}
}
survey.AskOne(prompt, &color, nil)
```


### MultiSelect
<img src="https://media.giphy.com/media/3o7bukX6PNQJo7JUwo/giphy.gif" width="300px"/>

```golang
days := []string{}
prompt = &survey.MultiSelect{
    Message: "What days do you prefer:",
    Options: []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"},
}
survey.AskOne(prompt, &days, nil)
```

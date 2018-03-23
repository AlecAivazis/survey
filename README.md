# Survey

[![Build Status](https://travis-ci.org/AlecAivazis/survey.svg?branch=feature%2Fpretty)](https://travis-ci.org/AlecAivazis/survey)
[![GoDoc](http://img.shields.io/badge/godoc-reference-5272B4.svg)](https://godoc.org/gopkg.in/AlecAivazis/survey.v1)

A library for building interactive prompts. 

![](https://zippy.gfycat.com/AmusingBossyArrowworm.gif)

```go
package main

import (
    "fmt"
    "gopkg.in/AlecAivazis/survey.v1"
)

// the questions to ask
var qs = []*survey.Question{
    {
        Name:     "name",
        Prompt:   &survey.Input{Message: "What is your name?"},
        Validate: survey.Required,
        Transform: survey.Title,
    },
    {
        Name: "color",
        Prompt: survey.NewSingleSelect().SetMessage("Choose a Color:").
            AddOption("red", nil, true).
            AddOption("blue", nil, false).
            AddOption("green", nil, false),
    },
    {
        Name: "age",
        Prompt:   &survey.Input{Message: "How old are you?"},
    },
}

func main() {
    // the answers will be written to this struct
    answers := struct {
        Name          string                  // survey will match the question and field names
        FavoriteColor *survey.Option `survey:"color"` // or you can tag fields to match a specific name
        Age           int                     // if the types don't match exactly, survey will try to convert for you
    }{}

    // perform the questions
    err := survey.Ask(qs, &answers)
    if err != nil {
        fmt.Println(err.Error())
        return
    }

    fmt.Printf("%s chose %s.", answers.Name, answers.FavoriteColor)
}
```

## Table of Contents

1. [Examples](#examples)
1. [Prompts](#prompts)
   1. [Input](#input)
   1. [Password](#password)
   1. [Confirm](#confirm)
   1. [Select](#select)
   1. [MultiSelect](#multiselect)
   1. [Editor](#editor)
1. [Validation](#validation)
   1. [Built-in Validators](#built-in-validators)
1. [Help Text](#help-text)
   1. [Changing the input rune](#changing-the-input-run)
1. [Custom Types](#custom-types)
1. [Customizing Output](#customizing-output)
1. [Versioning](#versioning)

## Examples

Examples can be found in the `examples/` directory. Run them
to see basic behavior:

```bash
go get gopkg.in/AlecAivazis/survey.v1

cd $GOPATH/src/gopkg.in/AlecAivazis/survey.v1

go run examples/simple.go
go run examples/validation.go
```

## Prompts

### Input

<img src="https://media.giphy.com/media/3og0IxS8JsuD9Z8syA/giphy.gif" width="400px"/>

```golang
name := ""
prompt := &survey.Input{
    Message: "ping",
}
survey.AskOne(prompt, &name, nil)
```

### Password

<img src="https://media.giphy.com/media/26FmQr6mUivkq71GE/giphy.gif" width="400px" />

```golang
password := ""
prompt := &survey.Password{
    Message: "Please type your password",
}
survey.AskOne(prompt, &password, nil)
```

### Confirm

<img src="https://media.giphy.com/media/3oKIPgsUmTp4m3eo4E/giphy.gif" width="400px"/>

```golang
name := false
prompt := &survey.Confirm{
    Message: "Do you like pie?",
}
survey.AskOne(prompt, &name, nil)
```

### Selection
All of the selection prompts have a `SetHelp` field which can be defined to provide more information to your users:

```golang
prompt.SetHelp("This is the help message shown")
```

The user can filter for options by typing while the prompt is active, and you can display a custom help message while filtering.
```golang
prompt.SetFilterMessage("You are now filtering the list")
```

The user can also press `esc` to toggle the ability cycle through the options with the j and k keys to do down and up respectively.
Or you can enable it in code.
```golang
prompt.SetVimMode(true)
```

By default, the Selection prompt is limited to showing 7 options at a time
and will paginate lists of options longer than that. To increase, you can either
change the global `survey.PageSize`, or set the `PageSize` field on the prompt:

```golang
prompt.SetPageSize(10)
```

#### Select

<img src="https://media.giphy.com/media/3oKIPxigmMu5YqpUPK/giphy.gif" width="400px"/>

```golang
color := &survey.Option{}
prompt := survey.NewSingleSelect().SetMessage("Select Color:").
			AddOption("red", nil, false).
			AddOption("blue", nil, false).
			AddOption("green", nil, true),
survey.AskOne(prompt, &color, nil)
```

A more complex example that uses the interface Value for the Options
```golang

type user struct {
	ID int `json:"id"`
	Name string `json:"name"`
	Username string `json:"username"`
	Email string `json:"email"`
	Address *address `json:"address"`
}

type address struct {
	Street string `json:"street"`
	Suite string `json:"suite"`
	City string `json:"city"`
	Zip string `json:"zipcode"`
}
type users = []*user

// the questions to ask
var userPrompt = survey.NewSingleSelect().SetMessage("Select User:")
var simpleQs = []*survey.Question{
	{
		Name: "user",
		Prompt: userPrompt,
		Validate: survey.Required,
	},
}

func init() {
	var (
		userData []byte
		request *http.Request
		response *http.Response
		err error
	)
	httpClient := &http.Client{Timeout: 5*time.Second}
	if request, err = http.NewRequest("GET", "https://jsonplaceholder.typicode.com/users", nil); err != nil {
		fmt.Println(err.Error())
		return
	}
	if response, err = httpClient.Do(request); err != nil {
		fmt.Println(err.Error())
		return
	}
	defer response.Body.Close()
	userData, err = ioutil.ReadAll(response.Body)
	var us users
	if err = json.Unmarshal(userData, &us); err != nil {
		fmt.Println(err.Error())
		return
	}

	for _, _user := range us {
		userPrompt.AddOption(_user.Username, _user, false)
	}

}
func main() {
	answers := struct {
		User *survey.Option
	}{}

	// ask the question
	err := survey.Ask(simpleQs, &answers)

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	_user := answers.User.Value.(*user)
	// Do something with the user that was selected
	...
}
```

### MultiSelect

<img src="https://media.giphy.com/media/3oKIP8lHYFtGeQDH0c/giphy.gif" width="400px"/>

```golang
days := make(survey.Options, 0)
prompt := survey.NewMultiSelect().SetMessage("What days do you prefer:").
        AddOption("Sunday", nil, true).
        AddOption("Monday", nil, false).
        AddOption("Tuesday", nil, false).
        AddOption("Wednesday", nil, false).
        AddOption("Thursday", nil, false).
        AddOption("Friday", nil, true).
        AddOption("Saturday", nil, true)
survey.AskOne(prompt, &days, nil)
```

### Editor

Launches the user's preferred editor (defined by the $EDITOR environment variable) on a
temporary file. Once the user exits their editor, the contents of the temporary file are read in as
the result. If neither of those are present, notepad (on Windows) or vim (Linux or Mac) is used.

## Validation

Validating individual responses for a particular question can be done by defining a
`Validate` field on the `survey.Question` to be validated. This function takes an
`interface{}` type and returns an error to show to the user, prompting them for another
response:

```golang
q := &survey.Question{
    Prompt: &survey.Input{Message: "Hello world validation"},
    Validate: func (val interface{}) error {
        // since we are validating an Input, the assertion will always succeed
        if str, ok := val.(string) ; !ok || len(str) > 10 {
            return errors.New("This response cannot be longer than 10 characters.")
        }
    }
}
```

### Built-in Validators

`survey` comes prepackaged with a few validators to fit common situations. Currently these
validators include:

| name         | valid types | description                                                 | notes                                                                                 |
| ------------ | ----------- | ----------------------------------------------------------- | ------------------------------------------------------------------------------------- |
| Required     | any         | Rejects zero values of the response type                    | Boolean values pass straight through since the zero value (false) is a valid response |
| MinLength(n) | string      | Enforces that a response is at least the given length       |                                                                                       |
| MaxLength(n) | string      | Enforces that a response is no longer than the given length |                                                                                       |

## Help Text

All of the prompts have a `Help` field which can be defined to provide more information to your users:

<img src="https://media.giphy.com/media/l1KVbc5CehW6r7pss/giphy.gif" width="400px" style="margin-top: 8px"/>

```golang
&survey.Input{
    Message: "What is your phone number:",
    Help:    "Phone number should include the area code",
}
```

### Changing the input rune

In some situations, `?` is a perfectly valid response. To handle this, you can change the rune that survey
looks for by setting the `HelpInputRune` variable in `survey/core`:

```golang
import (
    "gopkg.in/AlecAivazis/survey.v1"
    surveyCore "gopkg.in/AlecAivazis/survey.v1/core"
)

number := ""
prompt := &survey.Input{
    Message: "If you have this need, please give me a reasonable message.",
    Help:    "I couldn't come up with one.",
}

surveyCore.HelpIcon = '^'

survey.AskOne(prompt, &number, nil)
```

## Custom Types

survey will assign prompt answers to your custom types if they implement this interface:

```golang
type settable interface {
    WriteAnswer(field string, value interface{}) error
}
```

Here is an example how to use them:

```golang
type MyValue struct {
    value string
}
func (my *MyValue) WriteAnswer(name string, value interface{}) error {
     my.value = value.(string)
}

myval := MyValue{}
survey.AskOne(
    &survey.Input{
        Message: "Enter something:",
    },
    &myval,
    nil,
)
```

## Customizing Output

Customizing the icons and various parts of survey can easily be done by setting the following variables
in `survey/core`:

| name               | default | description                                                   |
| ------------------ | ------- | ------------------------------------------------------------- |
| ErrorIcon          | ✘       | Before an error                                               |
| HelpIcon           | ⓘ       | Before help text                                              |
| QuestionIcon       | ?       | Before the message of a prompt                                |
| SelectFocusIcon    | ❯       | Marks the current focus in `Select` and `MultiSelect` prompts |
| MarkedOptionIcon   | ◉       | Marks a chosen selection in a `MultiSelect` prompt            |
| UnmarkedOptionIcon | ◯       | Marks an unselected option in a `MultiSelect` prompt          |

## Versioning

This project tries to maintain semantic GitHub releases as closely as possible and relies on [gopkg.in](http://labix.org/gopkg.in)
to maintain those releases. Importing version 1 of survey would look like:

```golang
package main

import "gopkg.in/AlecAivazis/survey.v1"
```

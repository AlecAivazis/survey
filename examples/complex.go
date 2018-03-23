package main

import (
	"fmt"

	"gopkg.in/AlecAivazis/survey.v1"
	"net/http"
	"time"
	"io/ioutil"
	"encoding/json"
)

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
	// print the answers
	_user := answers.User.Value.(*user)
	fmt.Printf("%s has the username %s and thier address is %+v\r\n", _user.Name, _user.Username, _user.Address)
}

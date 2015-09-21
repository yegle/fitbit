package main

import (
	"fmt"
	"io/ioutil"
)

type ProfileCommand struct{}

func (ProfileCommand) Help() string {
	return "profile"
}

func (ProfileCommand) Synopsis() string {
	return "show profile"
}

func (ProfileCommand) Run(args []string) int {
	var (
		config *Config
		err    error
	)
	if config, err = LoadConfig(); err != nil {
		ERR.Print(err)
		return 1
	}

	client := config.GetOAuthClient()
	url := "https://api.fitbit.com/1/user/-/foods/log/water/goal.json"
	resp, err := client.Get(url)
	if err != nil {
		ERR.Print(err)
		return 1
	}

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ERR.Println(err)
		return 1
	}
	fmt.Println(string(content))
	return 0
}

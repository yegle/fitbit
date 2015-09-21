package main

import (
	"log"
	"os"

	"github.com/mitchellh/cli"
)

var (
	OUT    = log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
	ERR    = log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile)
	CONFIG *Config
)

func main() {
	c := cli.NewCLI("fitbit", "0.0.1")
	c.Args = os.Args[1:]
	c.Commands = map[string]cli.CommandFactory{
		"login":   func() (cli.Command, error) { return LoginCommand{}, nil },
		"drink":   func() (cli.Command, error) { return DrinkCommand{}, nil },
		"profile": func() (cli.Command, error) { return ProfileCommand{}, nil },
	}
	if exitStatus, err := c.Run(); err != nil {
		log.Println(err)
	} else {
		os.Exit(exitStatus)
	}
}

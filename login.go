package main

import "golang.org/x/oauth2"

const DEFAULT_PORT = 8080

var (
	AuthCodeFlow   = oauth2.SetAuthURLParam("response_type", "code")
	ExpireIn30Days = oauth2.SetAuthURLParam("expires_in", "2592000")
)

type LoginCommand struct{}

func (LoginCommand) Help() string {
	return "login"
}

func (LoginCommand) Synopsis() string {
	return "login command"
}

func (LoginCommand) Run(args []string) int {
	var (
		err    error
		config *Config
		token  *oauth2.Token
		code   string
	)
	if config, err = LoadConfig(); err != nil {
		ERR.Print(err)
		return 1
	}
	oauthConfig := config.GetOAuthConfig()
	url := oauthConfig.AuthCodeURL("state", AuthCodeFlow, ExpireIn30Days)
	OUT.Print("Getting access token. Please go to following URL and grant permission")
	OUT.Print(url)
	if code, err = GetAuthCode(); err != nil {
		ERR.Print(err)
		return 1
	}
	if token, err = oauthConfig.Exchange(oauth2.NoContext, code); err != nil {
		ERR.Print(err)
		return 1
	}
	config.UpdateAccessTokenAndSave(token)
	return 0
}

func GetAuthCode() (string, error) {
	var (
		err       error
		closeChan chan struct{}
		codeChan  chan string
	)
	if closeChan, codeChan, err = NewServer(DEFAULT_PORT); err != nil {
		return "", err
	}

	code := <-codeChan
	closeChan <- struct{}{}
	return code, nil
}

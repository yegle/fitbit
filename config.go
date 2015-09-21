package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path"

	"golang.org/x/oauth2"
)

const CONFIG_FILE = ".fitbit.json"
const CONFIG_FILE_PERMISSION = 0600
const INDENT = "    "

type Config struct {
	oauth2.Token
	ClientID     string
	ClientSecret string
}

func (c *Config) Save() error {
	var b []byte
	var err error
	if b, err = json.MarshalIndent(c, "", INDENT); err != nil {
		return err
	}
	output, err := os.OpenFile(path.Join(os.Getenv("HOME"), CONFIG_FILE), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}

	if _, err = output.Write(b); err != nil {
		return err
	}
	return nil
}

func (c *Config) UpdateAccessTokenAndSave(token *oauth2.Token) error {
	c.Token = *token
	return c.Save()
}

func (c *Config) AuthorizationHeader() []string {
	return []string{fmt.Sprintf("Bearer %s", c.AccessToken)}
}

func (c *Config) GetOAuthConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     c.ClientID,
		ClientSecret: c.ClientSecret,
		Scopes:       []string{"activity", "heartrate", "location", "nutrition", "profile", "settings", "sleep", "social", "weight"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://www.fitbit.com/oauth2/authorize",
			TokenURL: "https://api.fitbit.com/oauth2/token",
		},
		RedirectURL: fmt.Sprintf("http://127.0.0.1:%d/callback", DEFAULT_PORT),
	}
}

func (c *Config) GetOAuthClient() *http.Client {
	return c.GetOAuthConfig().Client(oauth2.NoContext, &c.Token)
}

func LoadConfig() (*Config, error) {
	var (
		f      *os.File
		err    error
		config = &Config{}
	)
	if f, err = os.Open(path.Join(os.Getenv("HOME"), CONFIG_FILE)); err != nil {
		return nil, err
	} else if info, err := f.Stat(); err != nil {
		return nil, err
	} else if perm := info.Mode().Perm(); perm != CONFIG_FILE_PERMISSION {
		return nil, fmt.Errorf("refuse to load config file: %q has permission %o, expect %o", info.Name(), perm, CONFIG_FILE_PERMISSION)
	}

	decoder := json.NewDecoder(f)
	if err := decoder.Decode(config); err != nil {
		return nil, err
	}
	return config, nil
}

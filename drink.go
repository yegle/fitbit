package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/google/go-querystring/query"
)

type DrinkCommand struct{}

func (c DrinkCommand) Help() string {
	return "drink command"
}

func (c DrinkCommand) Synopsis() string {
	return "NUMBER [glass|oz]"
}

type APIDate struct {
	time.Time
}

func (d APIDate) EncodeValues(key string, v *url.Values) error {
	(*v)[key] = []string{d.Format("2006-01-02")}
	return nil
}

type MyFloat64 float64

func (f MyFloat64) EncodeValues(key string, v *url.Values) error {
	if l, ok := (*v)[key]; ok {
		l = append(l, fmt.Sprintf("%.1f", f))
	} else {
		(*v)[key] = []string{fmt.Sprintf("%.1f", f)}
	}
	return nil
}

type LogWaterRequest struct {
	Amount MyFloat64 `url:"amount"`
	Date   APIDate   `url:"date"`
	Unit   string    `url:"unit"`
}

func (c DrinkCommand) Run(args []string) int {
	var (
		config *Config
		err    error
	)
	if config, err = LoadConfig(); err != nil {
		ERR.Println(err)
		return 1
	}
	if len(args) < 2 {
		OUT.Println(c.Synopsis())
		return 1
	}
	amount, err := strconv.ParseFloat(args[0], 64)
	if err != nil {
		ERR.Println(err)
		return 1
	}
	unit := args[1]
	switch unit {
	case "glass":
		amount *= 14
		unit = "fl oz"
	}
	r := &LogWaterRequest{
		Amount: MyFloat64(amount),
		Unit:   unit,
		Date:   APIDate{time.Now()},
	}
	v, err := query.Values(r)
	if err != nil {
		ERR.Println(err)
		return 1
	}
	req, err := http.NewRequest("POST", "https://api.fitbit.com/1/user/-/foods/log/water.json?"+v.Encode(), nil)
	req.ParseForm()
	client := config.GetOAuthClient()
	resp, err := client.Do(req)
	if err != nil {
		ERR.Println(err)
		return 1
	}
	fmt.Printf("%v\n", resp)
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ERR.Println(err)
		return 1
	}
	fmt.Println(string(content))
	return 0
}

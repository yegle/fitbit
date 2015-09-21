package main

import (
	"net/url"
	"strings"
)

type Values struct {
	url.Values
}

func (v Values) Read(p []byte) (n int, err error) {
	return strings.NewReader(v.Values.Encode()).Read(p)
}

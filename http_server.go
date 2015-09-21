package main

import (
	"fmt"
	"net"
	"net/http"
)

var authCodeChan = make(chan string, 1)

func RedirectHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	if e, ok := q["error"]; ok {
		fmt.Fprintf(w, "Error: %q. Description: %q", e, q["error_description"])
	} else if code, ok := q["code"]; !ok {
		fmt.Fprint(w, "Error: no code in query string!")
	} else {
		authCodeChan <- code[0]
		fmt.Fprint(w, "Successfully read the auth code. Close the window and check command line output")
	}
}

func NewServer(port int) (chan struct{}, chan string, error) {
	closeChan := make(chan struct{}, 1)
	var l net.Listener
	var err error
	if l, err = net.Listen("tcp", fmt.Sprintf(":%d", port)); err != nil {
		return nil, nil, err
	}

	go func() {
		<-closeChan
		l.Close()
		close(closeChan)
	}()

	http.HandleFunc("/callback", RedirectHandler)
	go func() {
		http.Serve(l, nil)
	}()
	return closeChan, authCodeChan, nil
}

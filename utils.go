package main

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httputil"
	"time"
)

// helper function to make http call
func HttpCall(method, rurl string, headers [][]string, buf *bytes.Buffer) *http.Response {
	var req *http.Request
	var err error
	if buf != nil {
		// POST request
		req, err = http.NewRequest(method, rurl, buf)
	} else {
		// GET, DELETE requests
		req, err = http.NewRequest(method, rurl, nil)
	}
	if err != nil {
		log.Printf("Unable to make request to %s, error: %s", rurl, err)
	}
	for _, v := range headers {
		if len(v) == 2 {
			req.Header.Set(v[0], v[1])
		}
	}
	if Config.Verbose > 1 {
		dump, err := httputil.DumpRequestOut(req, true)
		if err == nil {
			log.Println("request: ", string(dump))
		}
	}

	timeout := time.Duration(Config.HTTPTimeout) * time.Second
	client := &http.Client{Timeout: timeout}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Unable to get response from %s, error: %s", rurl, err)
	}
	if Config.Verbose > 2 {
		dump, err := httputil.DumpResponse(resp, true)
		if err == nil {
			log.Println("response: ", string(dump))
		}
	}
	log.Println(method, rurl, resp.Status)
	return resp
}

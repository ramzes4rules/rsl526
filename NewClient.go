package main

import (
	"net/http"
	"time"
)

func NewClient() *http.Client {

	client := &http.Client{}
	client.Timeout = 5 * time.Second
	return client

}

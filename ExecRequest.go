package main

import (
	"crypto/tls"
	"io"
	"net/http"
	"strings"
	"time"
)

func ExecRequest2(client *http.Client, url string, json string) error {

	//var timer time.Time
	//timer = time.Now()
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	method := "POST"
	payload := strings.NewReader(json)

	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	//client := &http.Client{}
	//client.Timeout = 5 * time.Second

	//
	//timer = time.Now()
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	//fmt.Printf("Request executed in %d ms\n", time.Since(timer).Milliseconds())

	_, err = io.ReadAll(res.Body)
	if err != nil {
		//fmt.Println(err)
		return err
	}

	//fmt.Println(string(body))

	return nil

}

func ExecRequest(url string, json string) error {

	//var timer time.Time
	//timer = time.Now()
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	method := "POST"
	payload := strings.NewReader(json)

	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	client.Timeout = 5 * time.Second

	//
	//timer = time.Now()
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	//fmt.Printf("Request executed in %d ms\n", time.Since(timer).Milliseconds())

	_, err = io.ReadAll(res.Body)
	if err != nil {
		//fmt.Println(err)
		return err
	}

	//fmt.Println(string(body))

	return nil

}

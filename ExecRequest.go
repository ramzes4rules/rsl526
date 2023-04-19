package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

func ExecRequest2(url string, json string, channel chan error) {

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	method := "POST"
	payload := strings.NewReader(json)

	//
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		channel <- err
	}

	//
	req.Header.Add("Content-Type", "application/json")

	//
	client := &http.Client{}
	client.Timeout = 5 * time.Second

	// call request
	_, err = client.Do(req)
	if err != nil {
		fmt.Printf("Exec request error: %v\n", err)
		channel <- err
	}

	//if res.StatusCode != 200 {
	//	body, err := io.ReadAll(res.Body)
	//	if err != nil {
	//		channel <- err
	//	}
	//	//fmt.Printf("\nBody: %s\n", string(body))
	//	channel <- fmt.Errorf("%s", string(body))
	//}
	channel <- nil
}

func ExecRequest(url string, json string) error {

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	method := "POST"
	payload := strings.NewReader(json)

	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	client.Timeout = 20 * time.Second

	// executing request
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	//fmt.Printf("Request executed in %d ms\n", time.Since(timer).Milliseconds())

	if res.StatusCode != 200 {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("%s", body)
		}
	}

	return nil

}

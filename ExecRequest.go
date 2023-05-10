package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

func ExecRequest2(url string, object string, channel chan error) {

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	method := "POST"
	payload := strings.NewReader(object)

	//
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		channel <- err
		return
	}

	//
	req.Header.Add("Content-Type", "application/json")

	//
	client := &http.Client{}
	client.Timeout = time.Duration(settings.Timeout) * time.Second

	// call request
	res, err := client.Do(req)
	if err != nil {
		//fmt.Printf("Exec request error: %v\n", err)
		channel <- err
		return
	}

	if res.StatusCode != 200 {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			channel <- err
			return
		}
		var response, _ = json.MarshalIndent(map[string]string{"Object": object, "Error": string(body)}, "", "\t")
		channel <- fmt.Errorf(string(response))
		return
	}
	channel <- nil
	return
}

//func ExecRequest(url string, json string) error {
//
//	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
//
//	method := "POST"
//	payload := strings.NewReader(json)
//
//	req, err := http.NewRequest(method, url, payload)
//	if err != nil {
//		return err
//	}
//
//	req.Header.Add("Content-Type", "application/json")
//
//	client := &http.Client{}
//	client.Timeout = 20 * time.Second
//
//	// executing request
//	res, err := client.Do(req)
//	if err != nil {
//		return err
//	}
//
//	// parsing response
//	if res.StatusCode != 200 {
//		body, err := io.ReadAll(res.Body)
//		if err != nil {
//			return fmt.Errorf("%s", err)
//		}
//		return fmt.Errorf("Failed to load:\n %s\nerror: %s", json, string(body))
//	}
//	return nil
//}

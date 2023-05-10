package main

import (
	"crypto/tls"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"
)

type ErrorInfo struct {
	Error  string
	Object any
}

func ExecRequest(url string, object any, channel chan *ErrorInfo) {

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	method := "POST"
	var obj, _ = json.Marshal(object)
	//fmt.Printf("%s", string(obj))
	payload := strings.NewReader(string(obj))

	//
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		channel <- &ErrorInfo{Error: err.Error()}
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
		channel <- &ErrorInfo{Error: err.Error()}
		return
	}

	if res.StatusCode != 200 {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			channel <- &ErrorInfo{Error: err.Error()}
			return
		}
		//var response, _ = json.MarshalIndent(map[string]string{"Object": object, "Error": string(body)}, "", "")
		//fmt.Printf(string(response))
		//channel <- fmt.Errorf(string(response))
		channel <- &ErrorInfo{Error: string(body), Object: string(obj)}
		return
	}
	channel <- nil
	return
}

//func ExecRequest2(url string, object string, channel chan error) {
//
//	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
//
//	method := "POST"
//	payload := strings.NewReader(object)
//
//	//
//	req, err := http.NewRequest(method, url, payload)
//	if err != nil {
//		channel <- err
//		return
//	}
//
//	//
//	req.Header.Add("Content-Type", "application/json")
//
//	//
//	client := &http.Client{}
//	client.Timeout = time.Duration(settings.Timeout) * time.Second
//
//	// call request
//	res, err := client.Do(req)
//	if err != nil {
//		//fmt.Printf("Exec request error: %v\n", err)
//		channel <- err
//		return
//	}
//
//	if res.StatusCode != 200 {
//		body, err := io.ReadAll(res.Body)
//		if err != nil {
//			channel <- err
//			return
//		}
//		//var response, _ = json.MarshalIndent(map[string]string{"Object": object, "Error": string(body)}, "", "")
//		//fmt.Printf(string(response))
//		channel <- fmt.Errorf(string(body))
//		//channel <- b
//		return
//	}
//	channel <- nil
//	return
//}

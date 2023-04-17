package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

var clients []*http.Client

func UploadCustomersAsync() error {

	//
	var timer time.Time
	var global time.Time
	//var waitGroup sync.WaitGroup
	var numbers = 10
	var cycleNumbers int
	var url = fmt.Sprintf("%s/api/customers/customer_import", settings.DestinationHost)

	// creat pool clients
	for i := 0; i < numbers; i++ {
		clients = append(clients, NewClient())
	}

	// read customers file
	fmt.Printf("Reading customers file. Be patient...\n")
	var customers []Customer
	err := ObjectRead(&customers, FileCustomers)
	if err != nil {
		return err
	}
	fmt.Printf("Loaded customers: %d\n", len(customers))

	cycleNumbers = len(customers) / numbers
	if len(customers)%numbers != 0 {
		cycleNumbers++
	}

	global = time.Now()
	//timer = time.Now()
	for i := 0; i < cycleNumbers; i++ {
		timer = time.Now()
		for j := 0; j < numbers; j++ {
			waitGroup.Add(1)
			var customer, err = json.MarshalIndent(customers[i+j], "", "\t")
			if err != nil {
				fmt.Printf("\tFailed to serialize customer '%s': %v\n", customers[i+j].CustomerID, err)
				continue
			}
			//var channel = make(chan struct{})
			go func(url string, customer []byte, number int) {
				//err := ExecRequest2(clients[j], url, string(customer))
				err := ExecRequest(url, string(customer))
				if err != nil {
					fmt.Printf("\tFailed to upload customer '%d': %v\n", number, err)
					return
				}
				waitGroup.Done()
			}(url, customer, i+j)
			//<-channel
		}
		waitGroup.Wait()
		fmt.Printf("Cyrcle %d. Uploaded customers from %d to %d. Time: %f, aver: %6.2f, total time: %10.2f minutes\n",
			i+1, i*numbers+1, i*numbers+numbers, time.Since(timer).Seconds(), float64(numbers)/time.Since(timer).Seconds(), time.Since(global).Minutes())
	}

	return nil
}

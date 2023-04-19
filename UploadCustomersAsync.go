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
	//for i := 0; i < numbers; i++ {
	//	clients = append(clients, NewClient())
	//}

	// read customers file
	fmt.Printf("Time start: %s\n", time.Now().Format(time.ANSIC))
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
			var customer, err = json.MarshalIndent(customers[i*numbers+j], "", "\t")
			if err != nil {
				fmt.Printf("\tFailed to serialize customer '%s': %v\n", customers[i+j].CustomerID, err)
				continue
			}
			//var channel = make(chan struct{})
			j := j
			go func(url string, customer []byte, number int) {
				//err := ExecRequest2(clients[j], url, string(customer))
				//fmt.Printf("i=%d, j=%d\n", i, j)

				err := ExecRequest(url, string(customer))

				if err != nil {
					fmt.Printf("\tFailed to upload customer '%d': %v\n", number, err)
					return
				}
				waitGroup.Done()
			}(url, customer, i*numbers+j)
			//<-channel
		}
		waitGroup.Wait()
		fmt.Printf("\rCyrcle %d. Uploaded customers from %d to %d. Time: %d ms, total: %10.2f min, average: %4.2f objects/second",
			i+1, i*numbers+1, i*numbers+numbers, time.Since(timer).Milliseconds(), time.Since(global).Minutes(), float64(i*numbers+numbers)/time.Since(global).Seconds())
	}

	fmt.Printf("Time finish: %s\n", time.Now().Format(time.ANSIC))
	return nil
}

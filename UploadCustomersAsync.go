package main

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

func UploadCustomersAsync() error {

	//
	var timer time.Time
	var global time.Time
	var waitGroup sync.WaitGroup
	var numbers = 40
	var cycleNumbers int
	var url = fmt.Sprintf("%s/api/customers/customer_import", settings.DestinationHost)

	//
	global = time.Now()

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

	timer = time.Now()
	for i := 0; i < cycleNumbers; i++ {
		timer = time.Now()
		for j := 0; j < numbers; j++ {
			waitGroup.Add(1)
			var customer, err = json.MarshalIndent(customers[i+j], "", "\t")
			if err != nil {
				fmt.Printf("\tFailed to serialize customer '%s': %v\n", customers[i+j].CustomerID, err)
				continue
			}
			go func(url string, customer []byte, number int) {
				defer waitGroup.Done()
				//fmt.Printf("Uploading customer: %d\n", number)
				err := ExecRequest(url, string(customer))
				if err != nil {
					fmt.Printf("\tFailed to upload customer '%d': %v\n", number, err)
					return
				}
			}(url, customer, i+j)
		}
		waitGroup.Wait()
		fmt.Printf("Cyrcle %d. Uploaded customers from %d to %d. Time: %f, total time: %f minutes\n",
			i+1, i*numbers+1, i*numbers+numbers, time.Since(timer).Seconds(), time.Since(global).Minutes())
	}

	return nil
}

package main

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

var waitGroup sync.WaitGroup

func UploadAccountsAsync() error {

	//
	var timer time.Time
	var global time.Time

	var numbers = 100
	var cycleNumbers int
	url := fmt.Sprintf("%s/api/accounts/accrual_to_loyalty_card", settings.DestinationHost)

	//

	// loading list of accounts
	fmt.Printf("Load account fro file %s. Be patient...", FileAccounts)
	var accounts []Account
	err := ObjectRead(&accounts, FileAccounts)
	if err != nil {
		return err
	}
	fmt.Printf("Loaded accounts numbers: %d", len(accounts))

	cycleNumbers = len(accounts) / numbers
	if len(accounts)%numbers != 0 {
		cycleNumbers++
	}

	global = time.Now()
	timer = time.Now()
	for i := 0; i < cycleNumbers; i++ {
		timer = time.Now()
		for j := 0; j < numbers; j++ {
			waitGroup.Add(1)
			var customer, err = json.MarshalIndent(accounts[i+j], "", "\t")
			if err != nil {
				fmt.Printf("\tFailed to serialize account customer: '%s': %v\n", accounts[i+j].LoyaltyCardId, err)
				continue
			}
			//var channel = make(chan struct{})
			go func(url string, customer []byte, number int) {
				defer waitGroup.Done()
				//fmt.Printf("Uploading customer: %d\n", number)
				err := ExecRequest(url, string(customer))
				if err != nil {
					fmt.Printf("\tFailed to upload account '%d': %v\n", number, err)
					return
				}
			}(url, customer, i+j)
			//<-channel
		}
		waitGroup.Wait()
		fmt.Printf("Cyrcle %d. Uploaded accounts from %d to %d. Time: %f, total time: %f minutes\n",
			i+1, i*numbers+1, i*numbers+numbers, time.Since(timer).Seconds(), time.Since(global).Minutes())
	}

	return nil
}

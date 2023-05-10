package main

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var waitGroup sync.WaitGroup

func UploadAccountsAsync() error {

	// declare variables
	var timer time.Time
	var global time.Time
	var loopNumbers int
	var results = make(chan error)
	var url = fmt.Sprintf("%s/api/accounts/accrual_to_loyalty_card", settings.DestinationHost)

	// print start time
	fmt.Printf("Start time: %s\n", time.Now().Format(time.ANSIC))

	// search file(s) to load
	files, _ := filepath.Glob(fmt.Sprintf("%s_?????%s", strings.TrimSuffix(FileAccounts, filepath.Ext(FileAccounts)), filepath.Ext(FileAccounts)))
	if len(files) == 0 {
		files, _ = filepath.Glob(FileAccounts)
	}
	if len(files) == 0 {
		return fmt.Errorf("files to loading not found")
	}

	//
	global = time.Now()
	for _, file := range files {

		// load list of accounts from file
		fmt.Printf("Loading account from file %s\n", file)
		var accounts []Account
		err := ObjectRead(&accounts, file)
		if err != nil {
			return err
		}
		fmt.Printf("Accounts loaded: %d", len(accounts))

		// calculate loops number
		loopNumbers = len(accounts) / settings.PacketSize
		if len(accounts)%settings.PacketSize != 0 {
			loopNumbers++
		}

		// run loop
		for i := 0; i < loopNumbers; i++ {
			timer = time.Now()

			end := settings.PacketSize
			if i+1 == loopNumbers {
				end = len(accounts) % settings.PacketSize
			}

			for j := 0; j < end; j++ {
				var customer, _ = json.MarshalIndent(accounts[i*settings.PacketSize+j], "", "\t")
				go ExecRequest2(url, string(customer), results)
			}

			for j := 0; j < end; j++ {
				err = <-results
				if err != nil {
					fmt.Printf("Error: %v", err)
				}
			}

			fmt.Printf("\rCyrcle %09d. Uploaded cards from %09d to %09d. Time: %05d ms, total: %05.2f min, average: %05.2f objects/second",
				i+1, i*settings.PacketSize+1, i*settings.PacketSize+settings.PacketSize, time.Since(timer).Milliseconds(),
				time.Since(global).Minutes(), float64(i*settings.PacketSize+settings.PacketSize)/time.Since(global).Seconds())

		}
	}
	return nil
}

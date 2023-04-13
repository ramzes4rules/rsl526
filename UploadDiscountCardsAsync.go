package main

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

func UploadDiscountCardsAsync() error {

	//
	var timer time.Time
	var global time.Time
	var waitGroup sync.WaitGroup
	var numbers = 40
	var cycleNumbers int
	url := fmt.Sprintf("%s/api/loyalty_cards/loyalty_card_import", settings.DestinationHost)

	//
	global = time.Now()

	// Загружаем файл кастомеров
	fmt.Printf("Reading file %s. Be patient...\n", FileDiscountCards)
	var discountCards []DiscountCard
	err := ObjectRead(&discountCards, FileDiscountCards)
	if err != nil {
		return err
	}
	fmt.Printf("Numbers of discount card read: %d\n", len(discountCards))

	cycleNumbers = len(discountCards) / numbers
	if len(discountCards)%numbers != 0 {
		cycleNumbers++
	}

	timer = time.Now()
	for i := 0; i < cycleNumbers; i++ {
		timer = time.Now()
		for j := 0; j < numbers; j++ {
			waitGroup.Add(1)
			var customer, err = json.MarshalIndent(discountCards[i+j], "", "\t")
			if err != nil {
				fmt.Printf("\tFailed to serialize card '%s': %v\n", discountCards[i+j].Id, err)
				continue
			}
			go func(url string, customer []byte, number int) {
				defer waitGroup.Done()
				//fmt.Printf("Uploading customer: %d\n", number)
				err := ExecRequest(url, string(customer))
				if err != nil {
					fmt.Printf("\tFailed to upload card '%d': %v\n", number, err)
					return
				}
			}(url, customer, i+j)
		}
		waitGroup.Wait()
		fmt.Printf("Cyrcle %d. Uploaded cards from %d to %d. Time: %f seconds, total time: %f minutes\n",
			i+1, i*numbers+1, i*numbers+numbers, time.Since(timer).Seconds(), time.Since(global).Minutes())
	}

	return nil
}

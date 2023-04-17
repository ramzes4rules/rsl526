package main

import (
	"encoding/json"
	"fmt"
	"time"
)

func UploadDiscountCardsAsync() error {

	//
	var timer time.Time
	var global time.Time
	//var waitGroup sync.WaitGroup
	var numbers = 10
	var cycleNumbers int
	url := fmt.Sprintf("%s/api/loyalty_cards/loyalty_card_import", settings.DestinationHost)

	// creat pool clients
	for i := 0; i < numbers; i++ {
		clients = append(clients, NewClient())
	}

	// read discount cards file
	fmt.Printf("Reading cards from file %s. Be patient...\n", FileDiscountCards)
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

	global = time.Now()
	//timer = time.Now()
	for i := 0; i < cycleNumbers; i++ {
		timer = time.Now()
		for j := 0; j < numbers; j++ {
			waitGroup.Add(1)
			var card, err = json.MarshalIndent(discountCards[i+j], "", "\t")
			if err != nil {
				fmt.Printf("\tFailed to serialize card '%s': %v\n", discountCards[i+j].Id, err)
				continue
			}
			go func(url string, card []byte, number int) {

				//fmt.Printf("Uploading card: %d\n", number)
				err := ExecRequest(url, string(card))
				if err != nil {
					fmt.Printf("\tFailed to upload card '%d': %v\n", number, err)
					return
				}
				waitGroup.Done()
			}(url, card, i+j)
		}
		waitGroup.Wait()
		fmt.Printf("\rCyrcle %d. Uploaded cards from %d to %d. Time: %f, aver: %6.2f, total time: %10.2f minutes",
			i+1, i*numbers+1, i*numbers+numbers, time.Since(timer).Seconds(), float64(numbers)/time.Since(timer).Seconds(), time.Since(global).Minutes())

	}

	return nil
}

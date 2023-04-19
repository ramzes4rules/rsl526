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
	var numbers = 5
	var cycleNumbers int
	url := fmt.Sprintf("%s/api/loyalty_cards/loyalty_card_import", settings.DestinationHost)

	// creat pool clients
	//for i := 0; i < numbers; i++ {
	//	clients = append(clients, NewClient())
	//}

	fmt.Printf("Time start: %s\n", time.Now().Format(time.ANSIC))

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
			var card, err = json.MarshalIndent(discountCards[i*numbers+j], "", "\t")
			if err != nil {
				fmt.Printf("\tFailed to serialize card '%s': %v\n", discountCards[i+j].Id, err)
				continue
			}
			go func(url string, card []byte, number int) {

				//fmt.Printf("Uploading card: %d\n", number)
				err := ExecRequest(url, string(card))
				if err != nil {
					fmt.Printf("\nFailed to upload CardID='%s': %v\n", discountCards[number].Id, err)
					return
				}
				waitGroup.Done()
			}(url, card, i*numbers+j)
		}
		waitGroup.Wait()
		fmt.Printf("\rCyrcle %d. Uploaded cards from %d to %d. Time: %d ms, total: %10.2f min, average: %4.2f objects/second",
			i+1, i*numbers+1, i*numbers+numbers, time.Since(timer).Milliseconds(), time.Since(global).Minutes(), float64(i*numbers+numbers)/time.Since(global).Seconds())

	}

	fmt.Printf("\nTime finish: %s\n", time.Now().Format(time.ANSIC))
	return nil
}

func UploadDiscountCardsAsyncC() error {

	// declare variables
	var timer time.Time
	var global time.Time
	var numbers = 5
	sizes := make(chan error)
	var cycleNumbers int
	var url = fmt.Sprintf("%s/api/loyalty_cards/loyalty_card_import", settings.DestinationHost)

	// print start time
	fmt.Printf("Start time: %s\n", time.Now().Format(time.ANSIC))

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

	//
	global = time.Now()
	for i := 0; i < cycleNumbers; i++ {
		timer = time.Now()

		// executing parallel request
		for j := 0; j < numbers; j++ {
			var card, _ = json.MarshalIndent(discountCards[i*numbers+j], "", "\t")
			go ExecRequest2(url, string(card), sizes)
		}

		// waiting for result
		for j := 0; j < numbers; j++ {
			err = <-sizes
		}

		//
		fmt.Printf("\rCyrcle %9d. Uploaded cards from %9d to %9d. Time: %5d ms, total: %5.2f min, average: %5.2f objects/second",
			i+1, i*numbers+1, i*numbers+numbers, time.Since(timer).Milliseconds(), time.Since(global).Minutes(), float64(i*numbers+numbers)/time.Since(global).Seconds())

	}

	fmt.Printf("\nTime finish: %s\n", time.Now().Format(time.ANSIC))
	return nil
}

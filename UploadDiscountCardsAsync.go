package main

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"time"
)

func UploadDiscountCardsAsyncC() error {

	// declare variables
	var timer time.Time
	var global time.Time
	var loopNumbers int
	var sizes = make(chan error)
	var url = fmt.Sprintf("%s/api/loyalty_cards/loyalty_card_import", settings.DestinationHost)

	// print start time
	fmt.Printf("Start time: %s\n", time.Now().Format(time.ANSIC))

	// search file(s) to load
	files, _ := filepath.Glob(fmt.Sprintf("%s_?????%s", strings.TrimSuffix(FileDiscountCards, filepath.Ext(FileDiscountCards)), filepath.Ext(FileDiscountCards)))
	if len(files) == 0 {
		files, _ = filepath.Glob(FileDiscountCards)
	}
	if len(files) == 0 {
		return fmt.Errorf("files to loading not found")
	}

	//
	global = time.Now()
	for _, file := range files {

		// read discount cards file
		fmt.Printf("Loading cards from file %s\n", file)
		var discountCards []DiscountCard
		err := ObjectRead(&discountCards, file)
		if err != nil {
			return err
		}
		fmt.Printf("Cards loaded: %d\n", len(discountCards))

		// calculate cycles numbers
		loopNumbers = len(discountCards) / settings.PacketSize
		if len(discountCards)%settings.PacketSize != 0 {
			loopNumbers++
		}

		// run loop
		for i := 0; i < loopNumbers; i++ {
			timer = time.Now()

			// executing parallel request
			end := settings.PacketSize
			if i+1 == loopNumbers {
				end = len(discountCards) % settings.PacketSize
			}

			for j := 0; j < end; j++ {
				var card, _ = json.MarshalIndent(discountCards[i*settings.PacketSize+j], "", "\t")
				go ExecRequest2(url, string(card), sizes)
			}

			// waiting for result
			for j := 0; j < end; j++ {
				err = <-sizes
				if err != nil {
					fmt.Printf("Error: %v", err)
				}
			}

			//
			fmt.Printf("\rCyrcle %09d. Uploaded cards from %09d to %09d. Time: %05d ms, total: %05.2f min, average: %05.02f objects/second",
				i+1, i*settings.PacketSize+1, i*settings.PacketSize+settings.PacketSize, time.Since(timer).Milliseconds(), time.Since(global).Minutes(), float64(i*settings.PacketSize+settings.PacketSize)/time.Since(global).Seconds())
		}
	}

	fmt.Printf("\nTime finish: %s\n", time.Now().Format(time.ANSIC))
	return nil
}

package main

import (
	"encoding/json"
	"fmt"
)

func UploadDiscountCards() error {

	// Загружаем файл кастомеров
	var discountCards []DiscountCard
	err := ObjectRead(&discountCards, FileDiscountCards)
	if err != nil {
		return err
	}
	fmt.Printf("Loaded discount card numbers: %d", len(discountCards))

	// upload cards loop
	for i := 0; i < len(discountCards); i++ {

		fmt.Printf("\rUploading discount card: %d", i+1)

		// serialisation card
		var discountCard, err = json.MarshalIndent(discountCards[i], "", "\t")
		if err != nil {
			fmt.Printf("\tFailed to serialize card '%s': %v\n", discountCards[i].Id, err)
			continue
		}

		// uploading card
		url := fmt.Sprintf("%s/api/loyalty_cards/loyalty_card_import", settings.DestinationHost)
		err = ExecRequest(url, string(discountCard))
		if err != nil {
			fmt.Printf("\tFailed to upload card '%s': %v\n", discountCards[i].Id, err)
			continue
		}
		fmt.Printf("OK")
	}

	return nil
}

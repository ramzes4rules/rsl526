package main

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"
)

const FileDiscountCards = "discountcards.json"

func ExportDiscountCards() error {

	// getting list of customers
	fmt.Printf("Getting list of cards from db. Be patient...\n")
	cards, err := GetListOfDiscountCards()
	if err != nil {
		fmt.Printf("Error getting list of cards: %v\n", err)
		return err
	}
	fmt.Printf("Numbers of discount card got: %d\n", len(cards))

	// getting discount card group values
	fmt.Printf("Getting list of discount card groups from db. Also be patient...\n")
	groups, err := GetDiscountCardGroups()
	if err != nil {
		fmt.Printf("Error getting list of discount card groups: %v\n", err)
		return err
	}
	fmt.Printf("Numbers of discount card group got: %d\n", len(groups))

	// reading mapping
	var mappings = DiscountCardMappings{}
	err = ObjectRead(&mappings, FileMappingDiscountCards)
	if err != nil {
		return err
	}

	// set default store
	var defStore = ""
	for _, store := range mappings[StoreId] {
		defStore = store
		break
	}

	// adding discount card groups
	fmt.Printf("Setting discount card groups to cards\n")
	for _, group := range groups {
		card, ok := cards[group.DiscountCardID]
		if ok == false {
			continue
		}
		card.LoyaltyCardGroupIds = append(card.LoyaltyCardGroupIds, mappings[LoyaltyCardGroupIds][group.DiscountCardGroupsID])
		cards[group.DiscountCardID] = card
	}
	fmt.Printf("Discount card groups set\n")

	// transforming list
	var out []DiscountCard
	var i = 1
	for _, card := range cards {

		//fmt.Printf("\rTransforming card: %d", i)

		card.IssueReasonId = mappings[IssueReasonId][IssueReasonId]
		card.LoyaltyCardSeriesId = mappings[LoyaltyCardSeriesId][LoyaltyCardSeriesId]
		card.LoyaltyCardCirculationId = mappings[LoyaltyCardCirculationId][LoyaltyCardCirculationId]
		card.StoreId = mappings[StoreId][card.StoreId]
		if card.StoreId == "" {
			card.StoreId = defStore
		}
		card.ActivationDate = time.Now()

		out = append(out, card)
		i++
	}

	// write cards to file(s)
	fmt.Printf("Writing cards to file(s). Be patient\n")
	if settings.SplitNumbers == 0 {
		fmt.Printf("Creating single file '%s'\n", FileDiscountCards)
		err = WriteObject(out, FileDiscountCards)
		if err != nil {
			return err
		}
	} else {
		numbers := len(out) / settings.SplitNumbers
		if len(out)%numbers != 0 {
			numbers++
		}
		fmt.Printf("Creating %d files\n", numbers)
		for i := 0; i < numbers; i++ {
			name := fmt.Sprintf("%s_%05d%s", strings.TrimSuffix(FileDiscountCards, filepath.Ext(FileDiscountCards)), i, filepath.Ext(FileDiscountCards))
			end := i*settings.SplitNumbers + settings.SplitNumbers - 1
			if end > len(out) {
				end = len(out)
			}
			part := out[i*settings.SplitNumbers : end]
			fmt.Printf("Creating file '%s'\n", name)
			err = WriteObject(part, name)
			if err != nil {
				return err
			}
		}
	}
	fmt.Printf("File(s) created\n")
	return nil
}

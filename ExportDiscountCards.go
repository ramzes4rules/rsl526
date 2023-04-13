package main

import (
	"fmt"
)

const FileDiscountCards = "discountcards.json"

func ExportDiscountCards() error {

	// getting list of customers
	fmt.Printf("Getting list of cards from db... Be patient")
	cards, err := GetListOfDiscountCards()
	if err != nil {
		fmt.Printf("Error getting list of customers: %v", err)
		return err
	}
	fmt.Printf("Number of discount cards got: %d\n", len(cards))

	// getting discount card group values
	fmt.Printf("Getting list of discount card groups from db... Also be patient")
	groups, err := GetDiscountCardGroups()
	if err != nil {
		return err
	}
	fmt.Printf("Number of discount card groups got: %d\n", len(groups))

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
	for _, group := range groups {
		card := cards[group.DiscountCardID]
		card.LoyaltyCardGroupIds = append(card.LoyaltyCardGroupIds, mappings[LoyaltyCardGroupIds][group.DiscountCardGroupsID])
		cards[group.DiscountCardID] = card
	}

	// transforming list
	var out []DiscountCard
	var i = 1
	for _, card := range cards {

		fmt.Printf("\rTransforming card: %d", i)

		card.IssueReasonId = mappings[IssueReasonId][IssueReasonId]
		card.LoyaltyCardSeriesId = mappings[LoyaltyCardSeriesId][LoyaltyCardSeriesId]
		card.LoyaltyCardCirculationId = mappings[LoyaltyCardCirculationId][LoyaltyCardCirculationId]
		card.StoreId = mappings[StoreId][card.StoreId]
		if card.StoreId == "" {
			card.StoreId = defStore
		}

		out = append(out, card)
		i++
	}

	fmt.Printf("\nWriting cards to file... Be patient last time")
	err = WriteObject(out, FileDiscountCards)
	if err != nil {
		return err
	}
	fmt.Printf("Discount cards exported!")

	return nil
}

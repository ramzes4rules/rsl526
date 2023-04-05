package main

import (
	"fmt"
)

const FileDiscountCards = "discountcards.json"

func ExportDiscountCards() error {

	// getting list of customers
	cards, err := GetListOfDiscountCards()
	if err != nil {
		fmt.Printf("Error getting list of customers: %v", err)
		return err
	}
	fmt.Printf("Got cards amount: %d\n\r", len(cards))

	// getting discount card group values
	groups, err := GetDiscountCardGroups()
	if err != nil {
		return err
	}
	fmt.Printf("Got discount card groups amount: %d\n", len(groups))

	// reading mapping
	var mappings = DiscountCardMappings{}
	err = ObjectRead(&mappings, FileMappingDiscountCards)
	if err != nil {
		return err
	}

	// adding discount card groups
	for _, group := range groups {
		card := cards[group.DiscountCardID]
		card.LoyaltyCardGroupIds = append(card.LoyaltyCardGroupIds, mappings[LoyaltyCardGroupIds][group.DiscountCardGroupsID])
		cards[group.DiscountCardID] = card
	}

	// transforming list
	var out []DiscountCard
	for _, card := range cards {

		fmt.Printf("\rTransforming card: %s", card.Id)

		card.IssueReasonId = mappings[IssueReasonId][IssueReasonId]
		card.LoyaltyCardSeriesId = mappings[LoyaltyCardSeriesId][LoyaltyCardSeriesId]
		card.LoyaltyCardCirculationId = mappings[LoyaltyCardCirculationId][LoyaltyCardCirculationId]
		card.StoreId = mappings[StoreId][card.StoreId]

		out = append(out, card)
	}

	err = WriteObject(out, FileDiscountCards)
	if err != nil {
		return err
	}

	return nil
}

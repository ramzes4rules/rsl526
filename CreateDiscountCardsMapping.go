package main

const FileMappingDiscountCards = "mapping_discount_cards.json"
const IssueReasonId = "issueReasonId"
const LoyaltyCardSeriesId = "loyaltyCardSeriesId"
const LoyaltyCardCirculationId = "loyaltyCardCirculationId"
const LoyaltyCardGroupIds = "loyaltyCardGroupIds"
const StoreId = "storeId"

// DiscountCardMapping mapping discount cards property
type DiscountCardMapping map[string]string
type DiscountCardMappings map[string]DiscountCardMapping

func CreateDiscountCardMapping() error {
	// storeId - list of store uids
	// issueReasonId - value
	// loyaltyCardGroupIds - list of discount card group uids
	// loyaltyCardSeriesId - value
	// loyaltyCardCirculationId - value

	var mappings = DiscountCardMappings{}

	var mapping = DiscountCardMapping{}
	mapping[IssueReasonId] = ""
	mappings[IssueReasonId] = mapping

	mapping = DiscountCardMapping{}
	mapping[LoyaltyCardSeriesId] = ""
	mappings[LoyaltyCardSeriesId] = mapping

	mapping = DiscountCardMapping{}
	mapping[LoyaltyCardCirculationId] = ""
	mappings[LoyaltyCardCirculationId] = mapping

	// getting list of stores
	stores, err := GetListOfStores()
	if err != nil {
		return err
	}
	mappings[StoreId] = stores

	// getting list of loyaltyCardGroupIds
	groups, err := GetListOfDiscountCardGroups()
	if err != nil {
		return err
	}
	mappings[LoyaltyCardGroupIds] = groups

	// saving to file
	err = WriteObject(mappings, FileMappingDiscountCards)
	if err != nil {
		return err
	}

	return nil
}

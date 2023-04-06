package main

const (
	CurrencyId          = "CurrencyId"
	FileMappingAccounts = "mapping_accounts.json"
)

func CreateAccountMapping() error {

	// adding currency mapping id
	var mappings = DiscountCardMappings{}
	var mapping = DiscountCardMapping{}
	mapping[CurrencyId] = ""
	mappings[CurrencyId] = mapping

	// saving to file
	err := WriteObject(mappings, FileMappingAccounts)
	if err != nil {
		return err
	}

	return nil

}

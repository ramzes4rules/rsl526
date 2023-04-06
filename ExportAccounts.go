package main

import (
	"fmt"
	"time"
)

const FileAccounts = "accounts.json"

type Account struct {
	OperationDate          time.Time  `json:"operationDate"`
	CurrencyId             string     `json:"currencyId"`
	LoyaltyCardId          string     `json:"loyaltyCardId"`
	ExpirationDate         *time.Time `json:"expirationDate"`
	ActivationDate         time.Time  `json:"activationDate"`
	Amount                 float64    `json:"amount"`
	InteractionChannelType string     `json:"interactionChannelType"`
}

func ExportAccounts() error {

	// getting list of account
	accounts, err := GetListOfAccounts()
	if err != nil {
		return err
	}
	fmt.Printf("Got accounts numbers: %d\n", len(accounts))

	// reading mapping values
	var mappings = DiscountCardMappings{}
	err = ObjectRead(&mappings, FileMappingAccounts)
	if err != nil {
		return err
	}
	fmt.Printf("Got account mappins numbers: %d\n", len(mappings))

	// transforming list
	var out []Account
	for _, account := range accounts {
		account.CurrencyId = mappings[CurrencyId][CurrencyId]
		out = append(out, account)
	}

	fmt.Printf("Writing accounts to file: '%s'\n", FileAccounts)
	err = WriteObject(out, FileAccounts)
	if err != nil {
		return err
	}
	//fmt.Printf("Done")
	return nil
}

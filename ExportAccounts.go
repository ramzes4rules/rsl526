package main

import (
	"fmt"
	"path/filepath"
	"strings"
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
	fmt.Printf("Transforming accounts\n")
	var out []Account
	for _, account := range accounts {
		account.CurrencyId = mappings[CurrencyId][CurrencyId]
		out = append(out, account)
	}
	fmt.Printf("Accounts transformed\n")

	// write accounts to file(s)
	fmt.Printf("Writing accounts to file(s). Be patient\n")
	if settings.SplitNumbers == 0 {
		fmt.Printf("Creating single file '%s'\n", FileAccounts)
		err = WriteObject(out, FileAccounts)
		if err != nil {
			return err
		}
	} else {
		numbers := len(out) / settings.SplitNumbers
		if len(out)%numbers != 0 {
			numbers++
		}
		for i := 0; i < numbers; i++ {
			name := fmt.Sprintf("%s_%05d%s", strings.TrimSuffix(FileAccounts, filepath.Ext(FileAccounts)), i, filepath.Ext(FileAccounts))
			end := i*settings.SplitNumbers + settings.SplitNumbers
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

package main

import (
	"encoding/json"
	"fmt"
)

func UploadAccounts() error {

	// loading list of accounts
	var accounts []Account
	err := ObjectRead(&accounts, FileAccounts)
	if err != nil {
		return err
	}
	fmt.Printf("Loaded accounts numbers: %d", len(accounts))

	// uploading accounts loop
	for i := 0; i < len(accounts); i++ {

		fmt.Printf("\rUploading account: %d", i+1)

		// serialisation account
		var account, err = json.MarshalIndent(accounts[i], "", "\t")
		if err != nil {
			fmt.Printf("\tFailed to serialize account: %s\n", accounts[i].LoyaltyCardId)
			continue
		}

		// uploading account
		url := fmt.Sprintf("%s/api/accounts/accrual_to_loyalty_card", settings.DestinationHost)
		err = ExecRequest(url, string(account))
		if err != nil {
			fmt.Printf("\tFailed to upload account '%s': %v\n", accounts[i].LoyaltyCardId, err)
			continue
		}
		fmt.Printf("OK")
	}

	return nil
}

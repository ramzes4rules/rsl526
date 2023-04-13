package main

import (
	"encoding/json"
	"fmt"
)

func UploadCustomers() error {

	// Загружаем файл кастомеров
	var customers []Customer
	err := ObjectRead(&customers, FileCustomers)
	if err != nil {
		return err
	}
	fmt.Printf("Loaded customers: %d", len(customers))

	// uploading customers loop
	for i := 0; i < len(customers); i++ {

		fmt.Printf("\rUploading customer: %d", i+1)

		// serialisation customer
		var customer, err = json.MarshalIndent(customers[i], "", "\t")
		if err != nil {
			fmt.Printf("\tFailed to serialize customer '%s': %v\n", customers[i].CustomerID, err)
			continue
		}

		// uploading customer
		url := fmt.Sprintf("%s/api/customers/customer_import", settings.DestinationHost)
		err = ExecRequest(url, string(customer))
		if err != nil {
			fmt.Printf("\tFailed to upload customer '%s': %v\n", customers[i].CustomerID, err)
			continue
		}
		//fmt.Printf("OK")
	}

	return nil
}

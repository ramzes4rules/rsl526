package main

import (
	"fmt"
	Leho "github.com/ramzes4rules/leho"
	"os"
)

const ConfigFile = "rsl526.json"

type Settings struct {
	Host            string
	User            string
	Password        string
	Port            string
	Database        string
	DestinationHost string
}

var (
	settings = Settings{
		Host:            "afanasy.retailloyalty.ru",
		User:            "loy",
		Password:        "M92bv1Dv3fss",
		Port:            "1433",
		Database:        "RSL_Afanasiy",
		DestinationHost: "https://10.14.5.127:7008",
	}
)

func main() {

	var err error

	// reading app settings
	var leho = Leho.Leho{Setting: settings, FileName: ConfigFile}
	err = leho.ReadSetting()
	if err != nil {
		_ = leho.WriteSetting()
	}

	// checking number of options
	if len(os.Args) == 1 {
		fmt.Printf("Want to read help? Call: %s -h", os.Args[0])
		return
	}

	// analyzing options
	var usage = fmt.Sprintf(""+
		"Usage of %s\n"+
		"\nCreate mappings:\n"+
		"%s -m      create all mappings\n"+
		"%s -m -c   create customers mapping\n"+
		"%s -m -d   create discount cards mapping\n"+
		"%s -m -a   create accounts mapping\n"+
		"\nExport data:\n"+
		"%s -e      export all data set\n"+
		"%s -e -c   export customers to file %s\n"+
		"%s -e -d   export discount cards to file %s\n"+
		"%s -e -a   export accounts (card balance) to file %s\n"+
		"\nUpload data:\n"+
		"%s -u      upload all data from files\n"+
		"%s -u -c   upload customers from file %s\n"+
		"%s -u -d   upload discount cards from file %s\n"+
		"%s -u -a   upload accounts from file %s\n"+
		"\n%s -h      display this helo",
		os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0],
		os.Args[0], os.Args[0], FileCustomers, os.Args[0], FileDiscountCards, os.Args[0], FileAccounts, os.Args[0],
		os.Args[0], FileCustomers, os.Args[0], FileDiscountCards, os.Args[0], FileAccounts, os.Args[0])
	switch os.Args[1] {
	case "-h":
		fmt.Println(usage)
	case "-m":
		option := ""
		if len(os.Args) >= 3 {
			option = os.Args[2]
		}
		switch option {
		case "-c":
			fmt.Printf("Creating customers mapping\n")
			err = CreateCustomerMappings()
			if err != nil {
				fmt.Printf("Failed to create customers mapping: %v\n", err)
			} else {
				fmt.Printf("Customers mapping created\n")
			}
		case "-d":
			fmt.Printf("Creating discount cards mapping\n")
			err = CreateDiscountCardMapping()
			if err != nil {
				fmt.Printf("Failed to create discount cards mapping: %v\n", err)
			} else {
				fmt.Printf("Discount cards mapping created\n")
			}
		case "-a":
			fmt.Printf("Creating accounts mapping\n")
			err = CreateAccountMapping()
			if err != nil {
				fmt.Printf("Failed to create accounts mapping: %v\n", err)
			} else {
				fmt.Printf("Accounts mapping created\n")
			}
		default:
			//
			fmt.Printf("Creating mapping for all instance\n")
			err = CreateCustomerMappings()
			if err != nil {
				fmt.Printf("Failed to create customer mapping: %v\n", err)
			} else {
				fmt.Printf("Customer mapping created\n")
			}

			//
			fmt.Printf("Creating discount cards mapping\n")
			err = CreateDiscountCardMapping()
			if err != nil {
				fmt.Printf("Failed to create discount cards mapping: %v\n", err)
			} else {
				fmt.Printf("Discount cards mapping created\n")
			}

			//
			fmt.Printf("Creating accounts mapping\n")
			err = CreateAccountMapping()
			if err != nil {
				fmt.Printf("Failed to create accounts mapping: %v\n", err)
			} else {
				fmt.Printf("Accounts mapping created\n")
			}
		}
	case "-e":
		option := ""
		if len(os.Args) >= 3 {
			option = os.Args[2]
		}
		switch option {
		case "-c":
			fmt.Printf("Exporting customers...\n")
			err = ExportCustomers()
			if err != nil {
				fmt.Printf("Failed to export customers: %v\n", err)
			} else {
				fmt.Printf("Customers exported\n")
			}
		case "-d":
			fmt.Printf("Exporting discount cards...\n")
			err = ExportDiscountCards()
			if err != nil {
				fmt.Printf("Failed to export discount cards: %v\n", err)
			} else {
				fmt.Printf("Discount cards exported\n")
			}
		case "-a":
			fmt.Printf("Exporting accounts...\n")
			err = ExportAccounts()
			if err != nil {
				fmt.Printf("Failed to export accounts: %v\n", err)
			} else {
				fmt.Printf("Accounts exported\n")
			}
		default:
			//
			fmt.Printf("Exporting customers...\n")
			err = ExportCustomers()
			if err != nil {
				fmt.Printf("Failed to export customers: %v\n", err)
			} else {
				fmt.Printf("Customers exported\n")
			}

			//
			fmt.Printf("Exporting discount cards...\n")
			err = ExportDiscountCards()
			if err != nil {
				fmt.Printf("Failed to export discount cards: %v\n", err)
			} else {
				fmt.Printf("Discount cards exported\n")
			}

			//_ = ExportAccounts()
			fmt.Printf("Exporting accounts...\n")
			err = ExportAccounts()
			if err != nil {
				fmt.Printf("Failed to export accounts: %v\n", err)
			} else {
				fmt.Printf("Accounts exported\n")
			}
		}

	case "-u":
		option := ""
		if len(os.Args) >= 3 {
			option = os.Args[2]
		}

		switch option {
		case "-c":
			fmt.Printf("Uploading customers...\n")
			err = UploadCustomers()
			if err != nil {
				fmt.Printf("Failed to upload customers: %v\n", err)
			} else {
				fmt.Printf("Customers uploaded\n")
			}
		case "-d":
			fmt.Printf("Uploading discount cards...\n")
			err = UploadDiscountCards()
			if err != nil {
				fmt.Printf("Failed to upload discount cards: %v\n", err)
			} else {
				fmt.Printf("Discount cards uploaded\n")
			}
		case "-a":
			fmt.Printf("Uploading accounts...\n")
			err = UploadAccounts()
			if err != nil {
				fmt.Printf("Failed to upload accounts: %v\n", err)
			} else {
				fmt.Printf("Accounts uploaded\n")
			}
		default:
			//
			fmt.Printf("Uploading customers...\n")
			err = UploadCustomers()
			if err != nil {
				fmt.Printf("Failed to upload customers: %v\n", err)
			} else {
				fmt.Printf("Customers uploaded\n")
			}

			//
			fmt.Printf("Uploading discount cards...\n")
			err = UploadDiscountCards()
			if err != nil {
				fmt.Printf("Failed to upload discount cards: %v\n", err)
			} else {
				fmt.Printf("Discount cards uploaded\n")
			}

			//_ = ExportAccounts()
			fmt.Printf("Uploading accounts...\n")
			err = UploadAccounts()
			if err != nil {
				fmt.Printf("Failed to upload accounts: %v\n", err)
			} else {
				fmt.Printf("Accounts uploaded\n")
			}
		}

	default:
		fmt.Println(usage)
	}
}

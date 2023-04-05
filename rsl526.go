package main

import (
	_ "encoding/json"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	_ "os"
)

const configFile = "export.json"

type Settings struct {
	Host     string
	User     string
	Password string
	Port     string
	Database string
}

var (
	settings = Settings{
		Host:     "afanasy.retailloyalty.ru",
		User:     "loy",
		Password: "M92bv1Dv3fss",
		Port:     "1433",
		Database: "RSL_Afanasiy",
	}
)

func main() {

	var err error
	//err := CreateCustomerMappins()
	//if err != nil {
	//	fmt.Printf("Failed to save mapping files: %v", err)
	//}

	//ExportCustomers()

	//err = CreateDiscountCardMapping()
	//if err != nil {
	//	fmt.Printf(err.Error())
	//}

	err = ExportDiscountCards()
	if err != nil {
		fmt.Printf(err.Error())
	}

}

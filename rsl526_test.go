package main

import (
	"fmt"
	"testing"
)

/************************************ TEST FOR MAPPING ************************************/
func TestCreateCustomerMappings(t *testing.T) {

}

func TestCreateDiscountCardMapping(t *testing.T) {

}

func TestCreateAccountMapping(t *testing.T) {

}

/************************************ TEST FOR EXPORTING ************************************/
func TestExportCustomers(t *testing.T) {

}

func TestExportDiscountCards(t *testing.T) {

}

func TestExportAccounts(t *testing.T) {

}

/************************************ TEST FOR UPLOADING ************************************/
func TestUploadObjects(t *testing.T) {
	var url = fmt.Sprintf("%s/api/loyalty_cards/loyalty_card_import", settings.DestinationHost)
	err := UploadObjects(url, FileDiscountCards)
	if err != nil {
		t.Errorf(err.Error())
	}
}

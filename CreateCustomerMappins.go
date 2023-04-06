package main

import (
	"fmt"
)

const FileMappingCustomers = "mapping_customers.json"
const LocalitiesID = "LocalitiesID"
const PropertyIDs = "PropertyIDs"
const FileMappingEnumValues = "mapping_enum_values.json"

type Localities map[string]string
type EnumProperty map[string]string
type EnumProperties map[string]EnumProperty

func CreateCustomerMappings() error {

	//
	var mappings = DiscountCardMappings{}

	// getting list of enum properties value
	mappings, err := GetListOfEnumPropertyValues()
	if err != nil {
		return err
	}
	fmt.Printf("Mapping enum properties created")

	// getting list of customer properties except enum
	properties, err := GetListOfProperties()
	if err != nil {
		return err
	}
	fmt.Printf("Got customer properties amount: %d\n\r", len(properties))

	// adding customer properties
	var mapping = DiscountCardMapping{}
	for _, property := range properties {
		if property.PresetType == "" {
			mapping[property.Name] = ""
		}
	}
	mappings[PropertyIDs] = mapping

	// getting list of localities from rsl5
	localities, err := GetListOfLocalities()
	if err != nil {
		return err
	}
	mappings[LocalitiesID] = localities

	// saving to file
	err = WriteObject(mappings, FileMappingCustomers)
	if err != nil {
		return err
	}

	return nil
}

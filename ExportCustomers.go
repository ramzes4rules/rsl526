package main

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"
)

const FileCustomers = "customers.json"

type CustomerPropertyValue struct {
	PropertyUID         string     `json:"propertyId"`
	IntValue            *int32     `json:"intValue"`
	StringValue         *string    `json:"stringValue"`
	DateValue           *time.Time `json:"dateValue"`
	BooleanValue        *bool      `json:"booleanValue"`
	EnumPropertyValueId *string    `json:"enumPropertyValueId"`
}

func ExportCustomers() error {

	// getting list of customers
	customers, err := GetListOfCustomers()
	if err != nil {
		return err
	}
	fmt.Printf("Got customers numbers: %d\n", len(customers))

	// getting list of customers properties
	properties, err := GetListOfProperties()
	if err != nil {
		return err
	}
	fmt.Printf("Got customer properties numbers: %d\n\r", len(properties))

	// reading mapping
	var mappings = DiscountCardMappings{}
	err = ObjectRead(&mappings, FileMappingCustomers)
	if err != nil {
		return err
	}
	fmt.Printf("Got mappings numbers: %d\n\r", len(mappings))

	// getting list of customers property values
	// types: Integer, Enum, String, Date, Boolean
	for _, property := range properties {

		switch property.PropertyType {
		case "String":
			// getting property values
			props, err := GetStringProperty(property.PropertyId)
			if err != nil {
				return fmt.Errorf("failed to get sting customer property (id=%d): %v", property.PropertyId, err)
			}
			fmt.Printf("\rFilling PropertyID=%d values: name='%s', type='%s', preset='%s', numbers=%d\n",
				property.PropertyId, property.Name, property.PropertyType, property.PresetType, len(properties))

			// adding common string property
			if property.PresetType == "" {
				for _, prop := range props {
					customer := customers[prop.CustomerID]
					value := CustomerPropertyValue{
						PropertyUID: mappings[PropertyIDs][property.Name],
						StringValue: &prop.StringValue,
					}
					customer.CustomerPropertyValues = append(customer.CustomerPropertyValues, value)
					customers[prop.CustomerID] = customer
				}
			}

			// preset type PHONE
			if property.PresetType == "PHONE" {
				for _, prop := range props {
					customer := customers[prop.CustomerID]
					customer.Communications = append(customer.Communications, Communication{
						CommunicationChanelType: CommunicationChanelTypePhone,
						Value:                   prop.StringValue,
						Confirmed:               true,
					})
					customers[prop.CustomerID] = customer
				}
			}

			// preset type EMAIL
			if property.PresetType == "EMAIL" {
				for _, prop := range props {
					customer := customers[prop.CustomerID]
					customer.Communications = append(customer.Communications, Communication{
						CommunicationChanelType: CommunicationChanelTypeEmail,
						Value:                   prop.StringValue,
						Confirmed:               true,
					})
					customers[prop.CustomerID] = customer
				}
			}
		case "Date":
			props, err := GetDateProperties(property.PropertyId)
			if err != nil {
				return fmt.Errorf("failed to %v", err)
			}
			fmt.Printf("Filling PropertyID=%d values: name='%s', type='%s', preset='%s', numbers=%d\n",
				property.PropertyId, property.Name, property.PropertyType, property.PresetType, len(properties))

			// adding date property values
			if property.PresetType == "" {
				for _, prop := range props {
					customer := customers[prop.CustomerID]
					value := CustomerPropertyValue{
						PropertyUID: mappings[PropertyIDs][property.Name],
						DateValue:   &prop.DateValue,
					}
					customer.CustomerPropertyValues = append(customer.CustomerPropertyValues, value)
					customers[prop.CustomerID] = customer
				}
			}

			// adding preset DOB property values
			if property.PresetType == "DOB" {
				for _, prop := range props {
					cus := customers[prop.CustomerID]
					cus.DOB = prop.DateValue
					customers[prop.CustomerID] = cus
				}
			}

		case "Boolean":

			// getting boolean property values
			props, err := GetBooleanProperties(property.PropertyId)
			if err != nil {
				return fmt.Errorf("failed to get boolean customer properties: %v", err)
			}
			fmt.Printf("Filling boolean property values: id=%d, name=%s, amount=%d\n", property.PropertyId, property.Name, len(props))

			// filling properties
			for _, prop := range props {
				customer := customers[prop.CustomerID]
				value := CustomerPropertyValue{
					PropertyUID:  mappings[PropertyIDs][property.Name],
					BooleanValue: &prop.BoolValue,
				}
				customer.CustomerPropertyValues = append(customer.CustomerPropertyValues, value)
				customers[prop.CustomerID] = customer
			}

			// finished
			fmt.Printf("Property values (PropertyID=%d) filled\n", property.PropertyId)

		case "Integer":

			//// getting boolean property values
			//props, err := GetBooleanProperties(property.PropertyId)
			//if err != nil {
			//	return fmt.Errorf("failed to get boolean customer properties: %v", err)
			//}
			//fmt.Printf("Filling boolean property values: id=%d, name=%s, amount=%d\n", property.PropertyId, property.Name, len(props))
			//
			//// filling properties
			//for _, prop := range props {
			//	customer := customers[prop.CustomerID]
			//	customer.CustomerPropertyValues = append(customer.CustomerPropertyValues, CustomerPropertyValue{
			//		PropertyUID:   "",
			//		PropertyName:  property.Name,
			//		PropertyValue: fmt.Sprintf("%t", prop.BoolValue),
			//	})
			//	customers[prop.CustomerID] = customer
			//}
			//
			//// finished
			//fmt.Printf("Property values (PropertyID=%d) filled\n", property.PropertyId)

		case "Enum":

			// getting enum property values
			props, err := GetPropertiesTypeEnum(property.PropertyId)
			if err != nil {
				return fmt.Errorf("failed to get enum property (id=%d) %v\n", property.PropertyId, err)
			}
			fmt.Printf("\rFilling PropertyID=%d values: name='%s', type='%s', preset='%s', numbers=%d\n",
				property.PropertyId, property.Name, property.PropertyType, property.PresetType, len(props))

			// adding enum customer property values
			if property.PresetType == "" {
				for _, prop := range props {
					//fmt.Printf("\rAdding enum values %d", i+1)
					customer := customers[prop.CustomerID]
					//mapping := mappings[property.Name]
					//propid := prop.EnumValue
					e := mappings[property.Name][prop.EnumValue]
					//fmt.Printf("\npropid: %s\n", propid)
					//fmt.Printf("enum: %s\n", value)

					value := CustomerPropertyValue{
						PropertyUID:         mappings[PropertyIDs][property.Name],
						EnumPropertyValueId: &e,
					}
					customer.CustomerPropertyValues = append(customer.CustomerPropertyValues, value)
					customers[prop.CustomerID] = customer
				}
			}

			// filling properties
			if property.PresetType == "GENDER" {
				for _, prop := range props {
					//fmt.Printf("\rAdding gender values %d", i)
					customer := customers[prop.CustomerID]
					var gender = mappings[property.Name][prop.EnumValue]
					customer.Gender = &gender
					customers[prop.CustomerID] = customer
				}
			}
		}
	}

	// transform object & setting other properties
	var out []Customer
	var i = 0
	for _, customer := range customers {
		i++
		fmt.Printf("\rTransforming Customer: %d", i)

		//
		customer.CreatedDate = time.Now()

		//
		if customer.localityID != "" {
			var uid = mappings[LocalitiesID][customer.localityID]
			customer.TerritorialDivisionId = &uid
			//customers[customer.CustomerID] = customer
		}

		if customer.Gender == nil {
			un := "Unknown"
			customer.Gender = &un
			//customers[customer.CustomerID] = customer
		}

		out = append(out, customer)
	}
	fmt.Printf("\nList of customers transormed\n")

	// write list of customer to file(s)
	fmt.Printf("Writing list of customers to file...\n")
	if settings.SplitNumbers == 0 {
		fmt.Printf("Creating a single file '%s'\n", FileAccounts)
		err = WriteObject(out, FileCustomers)
		if err != nil {
			return err
		}
	} else {
		// calculate number of files
		numbers := len(out) / settings.SplitNumbers
		if len(out)%numbers != 0 {
			numbers++
		}

		// write files in loop
		for i := 0; i < numbers; i++ {
			name := fmt.Sprintf("%s_%05d%s", strings.TrimSuffix(FileCustomers, filepath.Ext(FileCustomers)), i, filepath.Ext(FileCustomers))
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

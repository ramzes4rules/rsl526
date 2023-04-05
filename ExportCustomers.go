package main

import (
	"fmt"
)

const FileCustomers = "customers.json"

func ExportCustomers() error {

	// getting list of customers
	customers, err := GetListOfCustomers()
	if err != nil {
		fmt.Printf("Error getting list of customers: %v", err)
		return err
	}
	fmt.Printf("Got customers amount: %d\n\r", len(customers))

	// getting list of customers properties
	properties, err := GetListOfProperties()
	if err != nil {
		fmt.Printf("Error getting list of customer properties: %v", err)
		return err
	}
	fmt.Printf("Got customer properties amount: %d\n\r", len(properties))

	// getting list of customers property values
	// types: Integer, Enum, String, Date, Boolean
	for _, property := range properties {

		switch property.PropertyType {
		case "String":
			properties, err := GetStringProperty(property.PropertyId)
			if err != nil {
				fmt.Printf("Ошибка: %v\n\r", err)
				return fmt.Errorf("failed to get sting customer property (id=%d): %v", property.PropertyId, err)
			}
			//fmt.Printf("Получено строковых свойств %s: %d\n\r", property.Name, len(properties))
			fmt.Printf("Filling string property values: id=%d, name=%s, amount=%d\n", property.PropertyId, property.Name, len(properties))

			// preset type PHONE
			if property.PresetType == "PHONE" {
				for _, prop := range properties {
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
				for _, prop := range properties {
					customer := customers[prop.CustomerID]
					customer.Communications = append(customer.Communications, Communication{
						CommunicationChanelType: CommunicationChanelTypeEmail,
						Value:                   prop.StringValue,
						Confirmed:               true,
					})
					customers[prop.CustomerID] = customer
				}
			}

			// common string property
			if property.PresetType == "" {
				for _, prop := range properties {
					customer := customers[prop.CustomerID]
					customer.CustomerPropertyValues = []CustomerPropertyValue{{
						PropertyUID:   fmt.Sprintf("%sID", property.Name),
						PropertyName:  property.Name,
						PropertyValue: prop.StringValue,
					}}
					customers[prop.CustomerID] = customer
				}
			}

		case "Date":
			props, err := GetDateProperties(property.PropertyId)
			if err != nil {
				return fmt.Errorf("failed to %v", err)
			}
			//
			if property.PresetType == "DOB" {
				for _, prop := range props {
					cus := customers[prop.CustomerID]
					cus.DOB = prop.DateValue
					customers[prop.CustomerID] = cus
				}
			} else {
				for _, prop := range props {
					customer := customers[prop.CustomerID]
					customer.CustomerPropertyValues = append(customer.CustomerPropertyValues, CustomerPropertyValue{
						PropertyUID:   fmt.Sprintf("%sID", property.Name),
						PropertyName:  property.Name,
						PropertyValue: "",
					})
					customers[prop.CustomerID] = customer
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
				customer.CustomerPropertyValues = append(customer.CustomerPropertyValues, CustomerPropertyValue{
					PropertyUID:   "",
					PropertyName:  property.Name,
					PropertyValue: fmt.Sprintf("%t", prop.BoolValue),
				})
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
			fmt.Printf("Filling enum property values: id=%d, name=%s, amount=%d\n", property.PropertyId, property.Name, len(props))

			// reading mappings
			var mappings = EnumProperties{}
			err = ObjectRead(&mappings, FileMappingEnumValues)
			if err != nil {
				return fmt.Errorf("failed to read mapping enum properties: %v\n", err)
			}

			// filling properties
			if property.PresetType == "GENDER" {
				fmt.Printf("Property is preset type GENDER\n")
				for _, prop := range props {
					customer := customers[prop.CustomerID]
					var gender = mappings[property.Name][prop.EnumValue]
					customer.Gender = &gender
					//fmt.Printf("%d -> %v\n", prop.CustomerID, &gender)
					customers[prop.CustomerID] = customer
				}
			} else {
				fmt.Printf("Property is not preset type (name=%s)\n", property.Name)
				for _, prop := range props {
					customer := customers[prop.CustomerID]
					customer.CustomerPropertyValues = append(customer.CustomerPropertyValues, CustomerPropertyValue{
						PropertyUID:   fmt.Sprintf("%sID", property.Name),
						PropertyName:  property.Name,
						PropertyValue: mappings[property.Name][prop.EnumValue],
					})
					customers[prop.CustomerID] = customer
				}
			}
			fmt.Printf("Property values (PropertyID=%d) filled\n", property.PropertyId)
		}

	}

	// filling localities uid from mapping
	var localities = Localities{}
	err = ObjectRead(&localities, FileLocalities)
	if err != nil {
		return err
	}
	for _, customer := range customers {
		if customer.localityID != "" {
			var uid = localities[customer.localityID]
			customer.TerritorialDivisionId = &uid
			customers[customer.CustomerID] = customer
		}
	}

	// transform object & setting other properties
	var out []Customer
	for _, customer := range customers {
		//
		customer.InteractionChannel = "UserInterface"
		//
		if *customer.SecretCode == "" {
			customer.SecretCode = nil
		}
		if *customer.FirstName == "" {
			customer.FirstName = nil
		}
		if *customer.SecondName == "" {
			customer.SecondName = nil
		}
		if *customer.LastName == "" {
			customer.LastName = nil
		}

		out = append(out, customer)
	}

	err = WriteObject(out, FileCustomers)
	if err != nil {
		return err
	}

	return nil
}

package main

import (
	"database/sql"
	"fmt"
)

func GetPropertiesTypeEnum(id uint16) ([]CustomerProperty, error) {

	// connecting to dbase
	connString := fmt.Sprintf("server=%s;userid=%s;password=%s;port=%s;database=%s", settings.Host, settings.User, settings.Password, settings.Port, settings.Database)
	db, connectionError := sql.Open("mssql", connString)
	if connectionError != nil {
		fmt.Println(fmt.Errorf("error opening database: %v", connectionError))
		return []CustomerProperty{}, connectionError
	}

	// executing request
	rows, err := db.Query("SELECT CPV.CustomerID, EPV.Name FROM CustomerPropertyValues CPV LEFT JOIN EnumPropertyValues EPV on CPV.EnumPropertyValueID = EPV.EnumPropertyValueID WHERE CPV.CustomerID IS NOT NULL AND EPV.NAME IS NOT NULL AND CPV.PropertyID = $1;", id)
	if err != nil {
		fmt.Printf("failed to run query %v", err)
	}
	defer rows.Close()

	//
	var properties []CustomerProperty
	for rows.Next() {
		var property = CustomerProperty{}
		err = rows.Scan(&property.CustomerID, &property.EnumValue)
		if err != nil {
			fmt.Println(err)
			continue
			//return []CustomerProperty{}, err
		}
		properties = append(properties, property)
	}

	//
	err = rows.Err()
	if err != nil {
		fmt.Println(err)
		return []CustomerProperty{}, err
	}

	return properties, nil

}

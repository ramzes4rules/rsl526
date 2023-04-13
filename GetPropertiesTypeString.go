package main

import (
	"database/sql"
	"fmt"
)

func GetStringProperty(id uint16) ([]CustomerProperty, error) {

	//
	//connString := fmt.Sprintf("server=%s;userid=%s;password=%s;port=%s;database=%s", settings.Host, settings.User, settings.Password, settings.Port, settings.Database)
	//fmt.Println(connString)
	db, connectionError := sql.Open("mssql", settings.ConnString)
	if connectionError != nil {
		fmt.Println(fmt.Errorf("error opening database: %v", connectionError))
		return []CustomerProperty{}, connectionError
	}

	// executing request
	rows, err := db.Query("SELECT CPV.CustomerID, CPV.StringValue FROM CustomerPropertyValues CPV WHERE CPV.CustomerID IS NOT NULL AND CPV.StringValue IS NOT NULL AND CPV.StringValue <> '' AND CPV.PropertyID = $1;", id)
	if err != nil {
		return []CustomerProperty{}, fmt.Errorf("failed to run query getting string properties: %v", err)
	}

	var properties []CustomerProperty
	for rows.Next() {
		var property = CustomerProperty{}
		err = rows.Scan(&property.CustomerID, &property.StringValue)
		if err != nil {
			fmt.Printf("Erorr scanning string values: %v", err)
			continue
		}
		properties = append(properties, property)
	}
	err = rows.Err()
	if err != nil {
		fmt.Println(err)
		return []CustomerProperty{}, err
	}

	return properties, nil

}

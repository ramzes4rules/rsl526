package main

import (
	"database/sql"
	"fmt"
)

func GetDateProperties(id uint16) ([]CustomerProperty, error) {

	//
	connString := fmt.Sprintf("server=%s;userid=%s;password=%s;port=%s;database=%s", settings.Host, settings.User, settings.Password, settings.Port, settings.Database)
	//fmt.Println(connString)
	db, connectionError := sql.Open("mssql", connString)
	if connectionError != nil {
		fmt.Println(fmt.Errorf("error opening database: %v", connectionError))
		return []CustomerProperty{}, connectionError
	}

	// executing request
	rows, err := db.Query("SELECT CustomerID, DateValue FROM CustomerPropertyValues WHERE CustomerID IS NOT NULL AND DateValue IS NOT NULL AND PropertyID = $1", id)
	if err != nil {
		fmt.Printf("Failed to run query: %v", err)
		return []CustomerProperty{}, err
	}
	defer rows.Close()

	var properties []CustomerProperty
	for rows.Next() {
		var property = CustomerProperty{}
		err = rows.Scan(&property.CustomerID, &property.DateValue)
		if err != nil {
			fmt.Println(err)
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

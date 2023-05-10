package main

import (
	"database/sql"
	"fmt"
	_ "log"
)

func GetBooleanProperties(id uint16) ([]CustomerProperty, error) {

	// connecting to dbase
	//connString := fmt.Sprintf("server=%s;userid=%s;password=%s;port=%s;database=%s", settings.Host, settings.User, settings.Password, settings.Port, settings.Database)
	db, connectionError := sql.Open("mssql", settings.connString)
	if connectionError != nil {
		fmt.Println(fmt.Errorf("error opening database: %v", connectionError))
		return []CustomerProperty{}, connectionError
	}

	// executing request
	rows, err := db.Query("SELECT CustomerID, BooleanValue FROM CustomerPropertyValues WHERE CustomerID IS NOT NULL AND BooleanValue IS NOT NULL AND PropertyID = $1;", id)
	if err != nil {
		fmt.Printf("failed to run query %v", err)
	}
	defer rows.Close()

	//
	var properties []CustomerProperty
	for rows.Next() {
		var property = CustomerProperty{}
		err = rows.Scan(&property.CustomerID, &property.BoolValue)
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

package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

type Property struct {
	PropertyId   uint16
	Name         string
	PropertyType string
	PresetType   string
}

type CustomerProperty struct {
	CustomerID  string
	StringValue string
	DateValue   time.Time
	IntValue    int64
	BoolValue   bool
	EnumValue   string
}

func GetListOfProperties() ([]Property, error) {

	//
	connString := fmt.Sprintf("server=%s;userid=%s;password=%s;port=%s;database=%s", settings.Host, settings.User, settings.Password, settings.Port, settings.Database)
	//fmt.Println(connString)
	db, connectionError := sql.Open("mssql", connString)
	if connectionError != nil {
		fmt.Println(fmt.Errorf("error opening database: %v", connectionError))
		return []Property{}, connectionError
	}

	// executing request
	rows, err := db.Query("SELECT PropertyID, Name, PropertyType, COALESCE(PresetType, '') AS PresetType\nFROM Properties \nWHERE ObjectTypeId = 0\n  AND IsDeleted = 0\nORDER BY Position;")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	//
	var properties []Property
	for rows.Next() {
		var property = Property{}
		err = rows.Scan(&property.PropertyId, &property.Name, &property.PropertyType, &property.PresetType)
		if err != nil {
			fmt.Println(err)
			return []Property{}, err
		}
		properties = append(properties, property)
	}
	err = rows.Err()
	if err != nil {
		fmt.Println(err)
		return []Property{}, err
	}

	return properties, nil
}

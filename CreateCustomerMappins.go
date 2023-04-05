package main

import (
	"database/sql"
	"fmt"
)

const FileLocalities = "mapping_localities.json"
const FileMappingEnumValues = "mapping_enum_values.json"

type Localities map[string]string
type EnumProperty map[string]string
type EnumProperties map[string]EnumProperty

func CreateCustomerMappings() error {

	err := CreateMappingLocalities()
	if err != nil {
		return err
	}
	fmt.Printf("Mapping localities created")

	err = CreateMappingEnumProperties()
	if err != nil {
		return err
	}
	fmt.Printf("Mapping enum properties created")

	return nil
}

func CreateMappingLocalities() error {

	connString := fmt.Sprintf("server=%s;userid=%s;password=%s;port=%s;database=%s", settings.Host, settings.User, settings.Password, settings.Port, settings.Database)
	db, connectionError := sql.Open("mssql", connString)
	if connectionError != nil {
		return connectionError
	}

	// executing request
	rows, err := db.Query("SELECT Name, LocalityID FROM Localities WHERE IsDeleted = 0;")
	if err != nil {
		return err
	}
	defer rows.Close()

	//
	var localities = Localities{}
	for rows.Next() {
		var name = ""
		var uid = ""
		err = rows.Scan(&name, &uid)
		if err != nil {
			continue
		}
		localities[name] = ""
	}
	err = rows.Err()
	if err != nil {
		return err
	}

	err = WriteObject(localities, FileLocalities)
	if err != nil {
		return err
	}

	return nil
}

func CreateMappingEnumProperties() error {

	//
	connString := fmt.Sprintf("server=%s;userid=%s;password=%s;port=%s;database=%s", settings.Host, settings.User, settings.Password, settings.Port, settings.Database)
	db, connectionError := sql.Open("mssql", connString)
	if connectionError != nil {
		return connectionError
	}

	// request all enum properties
	rows, err := db.Query("SELECT Name, PropertyID, PresetType FROM Properties WHERE ObjectTypeId = 0 AND IsDeleted = 0 AND PropertyType = 'Enum' ORDER BY Position;")
	if err != nil {
		return err
	}

	// parse all enum properties from query
	var properties []Property
	for rows.Next() {
		var property = Property{}
		err = rows.Scan(&property.Name, &property.PropertyId, &property.PresetType)
		properties = append(properties, property)
	}

	// get all enum values
	var mappings = EnumProperties{}
	for _, property := range properties {
		rows, err = db.Query("SELECT P.Name, EPV.Name FROM EnumPropertyValues EPV LEFT JOIN Properties P on EPV.PropertyID = P.PropertyID WHERE EPV.PropertyID = $1;", property.PropertyId)
		if err != nil {
			return err
		}
		var mapping = EnumProperty{}
		var name = ""
		var value = ""
		for rows.Next() {
			err = rows.Scan(&name, &value)
			mapping[value] = ""
		}
		mappings[name] = mapping
	}

	// saving to file
	err = WriteObject(mappings, FileMappingEnumValues)
	if err != nil {
		return err
	}

	return nil
}

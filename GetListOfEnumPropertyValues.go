package main

import (
	"database/sql"
)

func GetListOfEnumPropertyValues() (DiscountCardMappings, error) {

	//
	//connString := fmt.Sprintf("server=%s;userid=%s;password=%s;port=%s;database=%s", settings.Host, settings.User, settings.Password, settings.Port, settings.Database)
	//connString := fmt.Sprintf("sqlserver://%s:%s@%s/%s", settings.User, settings.Password, settings.Host, settings.Database)
	//fmt.Println(connString)
	db, connectionError := sql.Open("mssql", settings.connString)
	if connectionError != nil {
		return DiscountCardMappings{}, connectionError
	}

	// request all enum properties
	rows, err := db.Query("SELECT Name, PropertyID, PresetType FROM Properties WHERE ObjectTypeId = 0 AND IsDeleted = 0 AND PropertyType = 'Enum' ORDER BY Position;")
	if err != nil {
		return DiscountCardMappings{}, err
	}

	// parse all enum properties from query
	var properties []Property
	for rows.Next() {
		var property = Property{}
		err = rows.Scan(&property.Name, &property.PropertyId, &property.PresetType)
		properties = append(properties, property)
	}

	// get all enum values
	var mappings = DiscountCardMappings{}
	for _, property := range properties {
		rows, err = db.Query("SELECT P.Name, EPV.Name FROM EnumPropertyValues EPV LEFT JOIN Properties P on EPV.PropertyID = P.PropertyID WHERE EPV.PropertyID = $1;", property.PropertyId)
		if err != nil {
			return DiscountCardMappings{}, err
		}
		var mapping = DiscountCardMapping{}
		var name = ""
		var value = ""
		for rows.Next() {
			err = rows.Scan(&name, &value)
			mapping[value] = ""
		}
		mappings[property.Name] = mapping
	}

	//
	return mappings, nil
}

package main

import (
	"database/sql"
	"fmt"
)

func GetListOfDiscountCardGroups() (DiscountCardMapping, error) {

	// open db connection
	//connString := fmt.Sprintf("server=%s;userid=%s;password=%s;port=%s;database=%s", settings.Host, settings.User, settings.Password, settings.Port, settings.Database)
	db, connectionError := sql.Open("mssql", settings.ConnString)
	if connectionError != nil {
		return DiscountCardMapping{}, connectionError
	}

	// executing request
	rows, err := db.Query("SELECT Name FROM DiscountCardGroups;")
	if err != nil {
		return DiscountCardMapping{}, err
	}

	// parsing response
	var mapping = DiscountCardMapping{}
	for rows.Next() {
		var name = ""
		err = rows.Scan(&name)
		if err != nil {
			fmt.Println(err)
			continue
		}
		mapping[name] = ""
	}
	err = rows.Err()
	if err != nil {
		return DiscountCardMapping{}, err
	}

	// finished
	return mapping, nil
}

package main

import (
	"database/sql"
	"fmt"
)

func GetListOfLocalities() (DiscountCardMapping, error) {

	//connString := fmt.Sprintf("server=%s;userid=%s;password=%s;port=%s;database=%s", settings.Host, settings.User, settings.Password, settings.Port, settings.Database)
	db, connectionError := sql.Open("mssql", settings.ConnString)
	if connectionError != nil {
		return DiscountCardMapping{}, connectionError
	}

	// executing request
	rows, err := db.Query("SELECT Name, LocalityID FROM Localities WHERE IsDeleted = 0;")
	if err != nil {
		return DiscountCardMapping{}, err
	}

	// parsing response
	var localities = DiscountCardMapping{}
	for rows.Next() {
		var name = ""
		var id = 0
		err = rows.Scan(&name, &id)
		if err != nil {
			fmt.Println(err)
			continue
		}
		localities[name] = ""
	}
	err = rows.Err()
	if err != nil {
		return DiscountCardMapping{}, err
	}

	// finished
	return localities, nil
}

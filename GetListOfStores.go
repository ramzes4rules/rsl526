package main

import (
	"database/sql"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
)

func GetListOfStores() (DiscountCardMapping, error) {

	// open db connection
	//connString := fmt.Sprintf("server=%s;userid=%s;password=%s;port=%s;database=%s", settings.Host, settings.User, settings.Password, settings.Port, settings.Database)
	//fmt.Printf("server=%s;userid=%s;password=%s;port=%s;database=%s", settings.Host, settings.User, settings.Password, settings.Port, settings.Database)
	db, connectionError := sql.Open("mssql", settings.ConnString)
	if connectionError != nil {
		return DiscountCardMapping{}, connectionError
	}

	// executing request
	rows, err := db.Query("SELECT Name FROM Stores;")
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

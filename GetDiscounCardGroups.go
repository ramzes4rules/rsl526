package main

import (
	"database/sql"
	"fmt"
)

type DiscountCardGroup struct {
	DiscountCardID       string `json:"DiscountCardID"`
	DiscountCardGroupsID string `json:"DiscountCardGroupsID"`
}

func GetDiscountCardGroups() ([]DiscountCardGroup, error) {
	//
	connString := fmt.Sprintf("server=%s;userid=%s;password=%s;port=%s;database=%s", settings.Host, settings.User, settings.Password, settings.Port, settings.Database)
	db, connectionError := sql.Open("mssql", connString)
	if connectionError != nil {
		fmt.Println(fmt.Errorf("error opening database: %v", connectionError))
		return []DiscountCardGroup{}, connectionError
	}

	// executing request
	rows, err := db.Query("SELECT DC.DiscountCardID, DCG.Name FROM DiscountCard2DiscountCardGroup DC2DCG LEFT JOIN DiscountCards DC on DC2DCG.DiscountCardID = DC.DiscountCardID LEFT JOIN DiscountCardGroups DCG on DC2DCG.DiscountCardGroupID = DCG.DiscountCardGroupID WHERE DC.IsDeleted = 0")
	if err != nil {
		return []DiscountCardGroup{}, err
	}

	// parsing
	var discountCardGroups []DiscountCardGroup
	for rows.Next() {
		var discountCardGroup = DiscountCardGroup{}
		err = rows.Scan(&discountCardGroup.DiscountCardID, &discountCardGroup.DiscountCardGroupsID)
		if err != nil {
			fmt.Println(err)
			continue
		}
		discountCardGroups = append(discountCardGroups, discountCardGroup)
	}
	err = rows.Err()
	if err != nil {
		fmt.Println(err)
		return []DiscountCardGroup{}, err
	}

	// finishing
	return discountCardGroups, nil
}

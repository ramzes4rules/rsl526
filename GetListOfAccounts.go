package main

import (
	"database/sql"
	"time"
)

type Accounts map[string]Account

func GetListOfAccounts() (Accounts, error) {

	//
	//connString := fmt.Sprintf("server=%s;userid=%s;password=%s;port=%s;database=%s", settings.Host, settings.User, settings.Password, settings.Port, settings.Database)
	db, connectionError := sql.Open("mssql", settings.connString)
	if connectionError != nil {
		return Accounts{}, connectionError
	}

	// executing request
	rows, err := db.Query("SELECT DISTINCT DC.AccountID, I.Bonus, " +
		"(SELECT TOP 1 DiscountCardID FROM DiscountCards WHERE AccountID = DC.AccountID AND DC.IsBlocked = 0 " +
		"AND DC.IsDeleted = 0 ORDER BY DiscountCardID DESC) FROM DiscountCards DC " +
		"LEFT JOIN Indicators I on DC.AccountID = I.AccountID WHERE DC.IsBlocked = 0 AND DC.IsDeleted = 0 AND I.Bonus != 0 ORDER BY AccountID;")
	if err != nil {
		return Accounts{}, err
	}

	//
	var accounts = Accounts{}
	for rows.Next() {

		var account = Account{}
		var id = 0
		err = rows.Scan(&id, &account.Amount, &account.LoyaltyCardId)
		//fmt.Printf("scanned: id=%d, amount=%f, cardid=%s\n", id, account.Amount, account.LoyaltyCardId)
		if err != nil {
			continue
		}
		accounts[account.LoyaltyCardId] = account
	}
	err = rows.Err()
	if err != nil {
		return Accounts{}, err
	}

	for id, account := range accounts {
		//
		account.InteractionChannelType = "UserInterface"
		account.OperationDate = time.Now()
		account.ActivationDate = time.Now()
		account.ExpirationDate = nil
		accounts[id] = account
	}

	//
	return accounts, nil
}

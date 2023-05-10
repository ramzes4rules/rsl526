package main

import (
	"database/sql"
	"fmt"
	"time"
)

type DiscountCard struct {
	Id                          string     `json:"id"`
	Now                         time.Time  `json:"now"`
	CreateDate                  time.Time  `json:"createDate"`
	CustomerId                  *string    `json:"customerId"`
	Number                      string     `json:"number"`
	NumberChangedDate           time.Time  `json:"numberChangedDate"`
	Barcode                     string     `json:"barcode"`
	BarcodeChangedDate          time.Time  `json:"barcodeChangedDate"`
	PinCode                     *string    `json:"pinCode"`
	PinCodeChangedDate          time.Time  `json:"pinCodeChangedDate"`
	ReadyToIssueDate            time.Time  `json:"readyToIssueDate"`
	StoreId                     string     `json:"storeId"`
	IssueReasonId               string     `json:"issueReasonId"`
	IssueDate                   time.Time  `json:"issueDate"`
	ActivationDate              time.Time  `json:"activationDate"`
	BlockDate                   *time.Time `json:"blockDate"`
	UnblockDate                 *time.Time `json:"unblockDate"`
	BlockReasonId               *string    `json:"blockReasonId"`
	ExpirationDateChangedDate   time.Time  `json:"expirationDateChangedDate"`
	ExpirationDate              time.Time  `json:"expirationDate"`
	DiscountCardGroupsBoundDate time.Time  `json:"discountCardGroupsBoundDate"`
	LoyaltyCardGroupIds         []string   `json:"loyaltyCardGroupIds"`
	LoyaltyCardSeriesId         string     `json:"loyaltyCardSeriesId"`
	LoyaltyCardCirculationId    string     `json:"loyaltyCardCirculationId"`
	InteractionChannel          string     `json:"interactionChannel"`
}
type DiscountCards map[string]DiscountCard

func GetListOfDiscountCards() (DiscountCards, error) {
	//
	//connString := fmt.Sprintf("server=%s;userid=%s;password=%s;port=%s;database=%s", settings.Host, settings.User, settings.Password, settings.Port, settings.Database)
	db, connectionError := sql.Open("mssql", settings.connString)
	if connectionError != nil {
		fmt.Println(fmt.Errorf("error opening database: %v", connectionError))
		return DiscountCards{}, connectionError
	}

	// executing request
	rows, err := db.Query("SELECT DC.DiscountCardID, COALESCE(DC.CreatedDate, ''), COALESCE(CU.CustomerID, ''), DC.Number, DC.Barcode, COALESCE(DC.PinCode, ''), COALESCE(S.Name, ''), COALESCE(DC.ActivationDate, ''), DC.ExpirationDate FROM DiscountCards DC LEFT JOIN Accounts AC on AC.AccountID = DC.AccountID LEFT JOIN Customers CU on CU.CustomerID = AC.CustomerID LEFT JOIN Stores S on DC.StoreID = S.StoreID WHERE DC.IsBlocked = 0 AND DC.IsDeleted = 0")
	if err != nil {
		return DiscountCards{}, err
	}

	// parsing
	var dicountCards = DiscountCards{}
	for rows.Next() {
		var discountCard = DiscountCard{}
		err = rows.Scan(&discountCard.Id, &discountCard.CreateDate, &discountCard.CustomerId, &discountCard.Number,
			&discountCard.Barcode, &discountCard.PinCode, &discountCard.StoreId, &discountCard.ActivationDate, &discountCard.ExpirationDate)
		if err != nil {
			fmt.Println(err)
			continue
		}
		dicountCards[discountCard.Id] = discountCard
	}
	err = rows.Err()
	if err != nil {
		fmt.Println(err)
		return DiscountCards{}, err
	}

	// filling default values
	for id, discountcard := range dicountCards {
		discountcard.Now = time.Now()
		discountcard.NumberChangedDate = time.Now()
		discountcard.BarcodeChangedDate = time.Now()
		discountcard.PinCodeChangedDate = time.Now()
		discountcard.ReadyToIssueDate = time.Now()
		discountcard.IssueDate = time.Now()
		discountcard.ExpirationDateChangedDate = time.Now()
		discountcard.DiscountCardGroupsBoundDate = time.Now()
		discountcard.InteractionChannel = "UserInterface"
		discountcard.BlockDate = nil
		discountcard.UnblockDate = nil
		discountcard.BlockReasonId = nil
		if *discountcard.CustomerId == "0" {
			discountcard.CustomerId = nil
		}
		if *discountcard.PinCode == "" {
			discountcard.PinCode = nil
		}
		//fmt.Printf("---> %s\n", discountcard.ActivationDate.Format("2006-01-02T15:04:05"))
		//if discountcard.ActivationDate.Format("2006-01-02T15:04:05") == "1900-01-01T00:00:00" {
		discountcard.ActivationDate = time.Now()
		//}

		dicountCards[id] = discountcard

	}

	return dicountCards, nil

}

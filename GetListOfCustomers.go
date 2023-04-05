package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

type Customer struct {
	CustomerID            string    `json:"id"`
	InteractionChannel    string    `json:"interactionChannel"`
	FirstName             *string   `json:"firstName"`
	SecondName            *string   `json:"secondName"`
	LastName              *string   `json:"lastName"`
	TerritorialDivisionId *string   `json:"territorialDivisionId"`
	localityID            string    `json:"localityID"`
	SecretCode            *string   `json:"password"`
	CreatedDate           time.Time `json:"operationDate"`
	DOB                   time.Time `json:"birthday"`
	phone                 string
	email                 string
	//gender                   string
	Gender                   *string                 `json:"gender"`
	Communications           []Communication         `json:"communications"`
	CustomerPropertyValues   []CustomerPropertyValue `json:"customerPropertyValues"`
	SendingVirtualCopyCheque bool                    `json:"sendingVirtualCopyCheque"`
}
type Customers map[string]Customer
type CustomerPropertyValue struct {
	PropertyUID   string
	PropertyName  string
	PropertyValue string
}

type Communication struct {
	CommunicationChanelType CommunicationChanelType `json:"communicationChanelType"`
	Value                   string                  `json:"value"`
	Confirmed               bool                    `json:"confirmed"`
}

type CommunicationChanelType string

const (
	CommunicationChanelTypePhone CommunicationChanelType = "Phone"
	CommunicationChanelTypeEmail CommunicationChanelType = "Email"
)

func GetListOfCustomers() (Customers, error) {
	//
	connString := fmt.Sprintf("server=%s;userid=%s;password=%s;port=%s;database=%s", settings.Host, settings.User, settings.Password, settings.Port, settings.Database)
	db, connectionError := sql.Open("mssql", connString)
	if connectionError != nil {
		fmt.Println(fmt.Errorf("error opening database: %v", connectionError))
		return Customers{}, connectionError
	}

	// executing request
	rows, err := db.Query("SELECT CU.CustomerID, COALESCE(CU.FirstName, ''), COALESCE(CU.SecondName, ''), COALESCE(CU.LastName, ''), COALESCE(L.Name, ''), COALESCE(CU.SecretCode, ''), CU.CreatedDate FROM Customers CU LEFT JOIN Localities L on CU.LocalityID = L.LocalityID WHERE CU.IsDeleted = 0")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	//
	var customers = Customers{}
	for rows.Next() {
		var customer = Customer{}
		err = rows.Scan(&customer.CustomerID, &customer.FirstName, &customer.SecondName, &customer.LastName, &customer.localityID, &customer.SecretCode, &customer.CreatedDate)
		if err != nil {
			fmt.Println(err)
			continue
		}
		//customers = append(customers, customer)
		customers[customer.CustomerID] = customer
	}
	err = rows.Err()
	if err != nil {
		fmt.Println(err)
		return Customers{}, err
	}

	return customers, nil

}

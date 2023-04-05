package main

import (
	_ "encoding/json"
	_ "github.com/denisenkom/go-mssqldb"
	_ "os"
)

const configFile = "export.json"

//
//type Server struct {
//	Server   string
//	Userid   string
//	Password string
//	Port     string
//	Database string
//}
//type TimeTable struct {
//	Time string
//	Done bool
//}
//
//type Config struct {
//	Server     Server
//	SQL        string
//	TimeTables []TimeTable
//	TimeOut    int
//}
//
//func (config *Config) load() {
//
//	data, err := ioutil.ReadFile(configFile)
//	if err != nil {
//		config.Server.Server = "sl098000LoyLSR1"
//		config.Server.Userid = "svc_service_1c_uzb_sql"
//		config.Server.Password = "9AF5YRGoHwdQ9IZqbIDr"
//		config.Server.Port = "1433"
//		config.Server.Database = "RS_LOYALTY"
//
//		config.SQL = "SELECT s.Name [Магазин], s.Address [Адрес магазина], cp.Number [Номер купона], " +
//			"i.Name [Наименование товара], c.Time [Дата чека], c.ChequeNo [Номер чека], cl.Amount, cl.Quantity " +
//			"FROM ChequeLineCoupons clc " +
//			"LEFT JOIN ChequeLines cl on cl.ChequeLineID = clc.ChequeLineID " +
//			"LEFT JOIN ChequeDiscounts cd on cd.ChequeLineID = clc.ChequeLineID " +
//			"LEFT JOIN Cheques c on c.ChequeID = cl.ChequeID " +
//			"LEFT JOIN Stores s on s.StoreID = c.StoreID " +
//			"LEFT JOIN Coupons cp on cp.CouponID = clc.CouponID " +
//			"LEFT JOIN Items i on i.ItemID = cl.ItemID " +
//			"WHERE cp.TemplateCouponID in (13, 14, 15) AND cd.AccrualID <> 0 " +
//			"ORDER BY c.Time"
//		config.TimeOut = 300
//
//		config.TimeTables = append(config.TimeTables, TimeTable{Time: time.Now().Format(time.RFC3339), Done: false})
//		config.TimeTables = append(config.TimeTables, TimeTable{Time: time.Now().AddDate(0, 0, 7).Format(time.RFC3339), Done: false})
//
//		config.save()
//		return
//	}
//
//	err = json.Unmarshal(data, config)
//	if err != nil {
//		log.Fatal("Ошибка парсинга:", err)
//	}
//
//	return
//
//}
//func (config *Config) save() {
//
//	var response, _ = json.MarshalIndent(config, "", "\t")
//
//	file, err := os.Create(configFile)
//	if err != nil {
//		log.Println(err)
//		return
//	}
//
//	_, err = file.WriteString(string(response))
//	if err != nil {
//		log.Println(err)
//		return
//	}
//
//	err = file.Close()
//	if err != nil {
//		log.Println(err)
//		return
//	}
//}
//
//var config = Config{}

type Settings struct {
	Host     string
	User     string
	Password string
	Port     string
	Database string
}

var (
	settings = Settings{
		Host:     "afanasy.retailloyalty.ru",
		User:     "loy",
		Password: "M92bv1Dv3fss",
		Port:     "1433",
		Database: "RSL_Afanasiy",
	}
)

func main() {

	//err := CreateCustomerMappins()
	//if err != nil {
	//	fmt.Printf("Failed to save mapping files: %v", err)
	//}

	ExportCustomers()

}

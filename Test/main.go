package main

import (
	rabbitMQ "Test/RabbitMQ"
	"Test/sql"
	"fmt"
	"strconv"
	"strings"
)

func main() {
	var db = sql.SqlOpen()
	defer db.Close()

	//GetSqlContent
	companyNames, companyCodes, companyUrls, apiUrls, err := sql.GetSqlContent(db)
	if err != nil {
		fmt.Println("(sqltest) Error getting content: " + err.Error())
	}
	fmt.Println(strings.Repeat("-", 100))
	// Now read the contents
	for i := range companyNames {
		fmt.Println("Company " + strconv.Itoa(i) + ": " + companyNames[i] + ", CompanyCode: " + strconv.FormatFloat(companyCodes[i], 'f', 2, 64) + ", CopmanyUrl: " + companyUrls[i] + ", ApiUrl: " + apiUrls[i])
	}

	//RabbitMQ Consumer
	rabbitMQ.Consumer(db)
}

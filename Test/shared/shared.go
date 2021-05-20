package shared

import "time"

type Configuration struct{
	AMQPURL string
	SQLURL string
}

var Config = Configuration{
	AMQPURL: "amqp://guest:guest@localhost:5672/",
	SQLURL: "sqlserver://accountinguser:testUser@192.168.1.1?database=CompanyBlog&connection+timeout=30",
}

type AddCompany struct {
	CompanyName string
	CompanyCode float64
	CompanyUrl string
	ApiUrl string
	CreatedDate time.Time
	UpdateDate time.Time
	IsDeleted bool
}
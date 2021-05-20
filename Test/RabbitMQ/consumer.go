package rabbitMQ

import (
	"Test/parser"
	"Test/shared"
	sql2 "Test/sql"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"os"
	"strconv"
	"strings"
)

func Consumer(db *sql.DB) {
	conn, err := amqp.Dial(shared.Config.AMQPURL)
	handleError(err, "Can't connect to AMQP")
	defer conn.Close()
	amqpChannel, err := conn.Channel()
	handleError(err, "Can't create a amqpChannel")

	defer amqpChannel.Close()

	queue, err := amqpChannel.QueueDeclare("company", false, false, false, false, nil)
	handleError(err, "Could not declare `add` queue")

	err = amqpChannel.Qos(1, 0, false)
	handleError(err, "Could not configure QoS")

	messageChannel, err := amqpChannel.Consume(
		queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	handleError(err, "Could not register consumer")
	stopChan := make(chan bool)

	go func() {
		log.Printf("Consumer ready, PID: %d", os.Getpid())
		for d := range messageChannel {
			fmt.Println(strings.Repeat("-", 100))
			log.Printf("Received a message: %s", d.Body)

			addCompany := &shared.AddCompany{}

			err := json.Unmarshal(d.Body, addCompany)
			if err != nil {
				log.Printf("Error decoding JSON: %s", err)
			}

			fmt.Println(strings.Repeat("-", 100))
			fmt.Printf("Company :%s - ApiUrl: %s\n", addCompany.CompanyName, addCompany.ApiUrl)

			log.Printf("Company %s of %f. ApiUrl : %s", addCompany.CompanyName, addCompany.CompanyCode, addCompany.ApiUrl)
			res, err2 := sql2.InsertSqlContent(db, addCompany)
			handleError(err2, "Could not Insert Product to Sql")
			log.Printf("Inserted Product ID : %d", res)

			if err := d.Ack(false); err != nil {
				log.Printf("Error acknowledging message : %s", err)
			} else {
				log.Printf("Acknowledged message")
			}

			// SQL List All Product
			companyName, companyCode, companyUrl, apiUrl, err := sql2.GetSqlContent(db)
			if err != nil {
				fmt.Println("(sqltest) Error getting content: " + err.Error())
			}
			fmt.Println(strings.Repeat("-", 100))
			// Now read the contents
			for i := range companyName {
				fmt.Println("Company Name " + strconv.Itoa(i) + ": " + companyName[i] + ", Company Code: " + strconv.FormatFloat(companyCode[i] , 'f', 2, 64) + ", Company url: " + companyUrl[i] + ", Api Url: " + apiUrl[i])
			}

			//Parse borakasmer.com
			articleList := parser.ParseWeb()
			for title,url :=range articleList{
				fmt.Println("Blog:", title, "=>", "Url:", url)
			}
		}
	}()

	// Stop for program termination
	<-stopChan
}

func handleError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

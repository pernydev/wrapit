package datapack

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

var handlers = map[string]func(filename string, content string, dp *DataPackage){
	`^account/user\.json$`:          accountHandler,
	`messages/c(\d+)/messages\.csv`: channelMessageHandler,
}

func accountHandler(filename string, content string, dp *DataPackage) {
	fmt.Println(content)
	var data map[string]interface{}
	err := json.Unmarshal([]byte(content), &data)
	if err != nil {
		panic(err)
	}

	dp.Account = Account{
		Email:                  data["email"].(string),
		Verifited:              data["verified"].(bool),
		HasMobile:              data["has_mobile"].(bool),
		NeedsEmailVerification: data["needs_email_verification"].(bool),
		PremiumUntil:           data["premium_until"].(string),
		Flags:                  int(data["flags"].(float64)),
		Phone:                  data["phone"].(string),
		IP:                     data["ip"].(string),
		User:                   User{},
	}

	dp.Account.User = User{
		ID:            data["id"].(string),
		Username:      data["username"].(string),
		Discriminator: int(data["discriminator"].(float64)),
		GlobalName:    data["global_name"].(string),
		Bot:           false,
	}

	dp.Account.Relationships = make([]Relationship, 0)
	for _, relationship := range data["relationships"].([]interface{}) {
		dp.Account.Relationships = append(dp.Account.Relationships, Relationship{
			ID:   relationship.(map[string]interface{})["id"].(string),
			Type: int(relationship.(map[string]interface{})["type"].(float64)),
			User: User{},
		})

		// if type is 1 add since
		if int(relationship.(map[string]interface{})["type"].(float64)) == 1 {
			dp.Account.Relationships[len(dp.Account.Relationships)-1].Since = relationship.(map[string]interface{})["since"].(string)
		}
	}

	dp.Account.Payments = make([]Payment, 0)
	for _, payment := range data["payments"].([]interface{}) {
		dp.Account.Payments = append(dp.Account.Payments, Payment{
			ID:             payment.(map[string]interface{})["id"].(string),
			Timestamp:      payment.(map[string]interface{})["created_at"].(string),
			Amount:         int(payment.(map[string]interface{})["amount"].(float64)),
			Currency:       payment.(map[string]interface{})["currency"].(string),
			AmountRefunded: int(payment.(map[string]interface{})["amount_refunded"].(float64)),
		})
	}

	dp.Account.Applications = make([]Application, 0)
}

func channelMessageHandler(filename string, content string, dp *DataPackage) {
	// Messages are in CSV format, so we need to parse them
	/*
		Example:
		ID,Timestamp,Contents,Attachments
		1150773626195410945,2023-09-11 12:43:32.825000+00:00,ICANN more like ICANNOT get my application accepted,
	*/

	// get the channel id from the filename
	re := regexp.MustCompile(`^messages/c(\d+)/messages\.csv$`)
	matches := re.FindStringSubmatch(filename)

	channel := Channel{
		ID:       matches[1],
		Name:     "Unknown",
		Messages: make([]Message, 0),
	}

	// Parse the CSV
	r := csv.NewReader(strings.NewReader(content))
	records, err := r.ReadAll()
	if err != nil {
		panic(err)
	}

	records = records[1:]

	for _, record := range records {
		message := Message{
			ID:        record[0],
			Timestamp: record[1],
			Content:   record[2],
		}

		if len(record) > 3 {
			message.Attachments = record[3]
		}

		channel.Messages = append(channel.Messages, message)
	}

}

func Parse(files map[string]string) *DataPackage {
	dp := &DataPackage{}

	for path, content := range files {
		for regex, handler := range handlers {
			if match, _ := regexp.MatchString(regex, path); match {
				handler(path, content, dp)
			}
		}
	}

	return dp
}

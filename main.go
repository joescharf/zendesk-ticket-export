package main

import (
	"encoding/csv"
	"github.com/adamar/zego/zego"
	"github.com/golang/glog"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strconv"
)

type ZendeskCfg struct {
	UserName  string `yaml:"username"`
	Token     string `yaml:"token"`
	Subdomain string `yaml:"subdomain"`
}

func main() {

	var zendeskCfg ZendeskCfg
	record := make([]string, 9)
	heading := []string{"id", "created_at", "type", "status", "subject", "requester_id", "submitter_id", "recipient", "description"}

	zendeskSettings, err := ioutil.ReadFile("config.yaml")
	if err == nil {
		err = yaml.Unmarshal(zendeskSettings, &zendeskCfg)
	}
	// If either error is nil, we bail...
	if err != nil {
		glog.Fatalf("Error - Failed to read config file. Exiting: ", err)
	}

	auth := zego.Auth{zendeskCfg.UserName + "/token", zendeskCfg.Token, zendeskCfg.Subdomain}

	file, err := os.Create("zendesk.csv")
	if err != nil {
		glog.Fatalln("Error creating output file: ", err)
	}
	defer file.Close()

	w := csv.NewWriter(file)
	if err = w.Write(heading); err != nil {
		glog.Fatalln("Error writing header to csv: ", err)
	}

	tickets, _ := auth.ListTickets()
	for _, t := range tickets.Tickets {
		record[0] = strconv.Itoa(int(t.Id))
		record[1] = t.CreatedAt
		record[2] = t.Type
		record[3] = t.Status
		record[4] = t.Subject
		record[5] = strconv.Itoa(int(t.RequesterId))
		record[6] = strconv.Itoa(int(t.SubmitterId))
		record[7] = t.Recipient
		record[8] = t.Description
		if err := w.Write(record); err != nil {
			glog.Fatalln("Error writing record to csv: ", err)
		}
		// fmt.Print(t.Subject + " ")
	}
	w.Flush()
}

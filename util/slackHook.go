package util

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/Grama-Check/Grama-Check-App/models"
)

var config Config
var slackErrID string
var slackIssueID string

func init() {
	var err error
	config, err = LoadConfig(".")
	if err != nil {
		log.Fatal("Error loading config:", err)
	}
	// slackIssueID = "https://hooks.slack.com/services/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"
	// slackErrID = "https://hooks.slack.com/services/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"

	slackIssueID = config.SlackIssueID
	slackErrID = config.SlackErrorID

}

func SendIssue(p models.Person, issue string) {
	var err error

	reqstr := fmt.Sprintf(`{"text":"%v code with error: %v"}`, p.NIC, issue)

	jsonBody := []byte(reqstr)

	bodyReader := bytes.NewReader(jsonBody)

	req, err := http.NewRequest("POST", slackIssueID, bodyReader)

	if err != nil {
		log.Println("err: couldn't create failed check slack request.: ", err)
	}

	_, err = http.DefaultClient.Do(req)

	if err != nil {
		log.Println("err: couldn't send issue to slack : ", err)

	}

}

func SendError(errCode int, errString string) {
	var err error
	datastr := strings.ReplaceAll(errString, "\"", " ")
	reqstr := fmt.Sprintf(`{"text":"%v code with error: %v"}`, errCode, datastr)

	jsonBody := []byte(reqstr)

	bodyReader := bytes.NewReader(jsonBody)

	req, err := http.NewRequest("POST", slackErrID, bodyReader)

	if err != nil {
		log.Println("err: couldn't create failed check slack request.: ", err)
	}

	_, err = http.DefaultClient.Do(req)

	if err != nil {
		log.Println("err: couldn't send issue to slack : ", err)

	}
}

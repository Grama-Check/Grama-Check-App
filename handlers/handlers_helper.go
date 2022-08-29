package handlers

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/Grama-Check/Grama-Check-App/util"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

var srv *gmail.Service

func init() {
	// Configuration
	ctx := context.Background()
	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, gmail.GmailSendScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err = gmail.NewService(ctx, option.WithHTTPClient(client))

	if err != nil {
		log.Fatalln("Unable to retrieve Gmail client:", err)
	}

}

func SendStatus(nic string) {

	from := "jhivan@gmail.com"
	person, err := queries.GetCheck(context.Background(), nic)
	if err != nil {
		util.SendError(http.StatusInternalServerError, nic+" "+err.Error())

		missingError := "sql: no rows in result set"

		if strings.EqualFold(err.Error(), missingError) {
			log.Println(http.StatusInternalServerError, "Please send a grama check request first")
			return
		}

		log.Println(http.StatusInternalServerError, err.Error())
		return
	}
	to := person.Email
	var plainTextContent string
	if person.Failed {
		plainTextContent = fmt.Sprintf("Status: One or more tests failed.\nIdentity Check Passed?%v\nAddress Check Passed?%v\nPolice Check Passed?%v\n", person.Idcheck, person.Addresscheck, person.Policecheck)
	} else {
		plainTextContent = fmt.Sprintf("Status: One or more tests  pending completion.\nIdentity Check Passed?%v\nAddress Check Passed?%v\nPolice Check Passed?%v\n", person.Idcheck, person.Addresscheck, person.Policecheck)

	}

	// Create message
	ret, _ := sendMail(from, to, "GramaCheck | Police check status", plainTextContent, srv)
	fmt.Println("Email sent:", ret)

}

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func sendMail(from string, to string, title string, message string, srv *gmail.Service) (bool, error) {
	// Create the message
	msgStr := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s", from, to, title, message)
	msg := []byte(msgStr)
	// Get raw
	gMessage := &gmail.Message{Raw: base64.URLEncoding.EncodeToString(msg)}

	// Send the message
	_, err := srv.Users.Messages.Send("me", gMessage).Do()
	if err != nil {
		fmt.Println("Could not send mail>", err)
		return false, err
	}
	return true, nil
}

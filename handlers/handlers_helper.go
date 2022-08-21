package handlers

import (
	"fmt"

	db "github.com/Grama-Check/Grama-Check-App/db/sqlc"
)

func SendStatus(person db.Check) {
	fmt.Println("here")
	// from := mail.NewEmail("jhivan", "jhivan@wso2.com")
	// subject := "Status of Police Check"
	// to := mail.NewEmail(person.Name, person.Email)
	// var plainTextContent string
	// if person.Failed {
	// 	plainTextContent = fmt.Sprintf("Status: One or more tests failed.\nIdentity Check Passed?%v\nAddress Check Passed?%v\nPolice Check Passed?%v,\n ", person.Idcheck, person.Addresscheck, person.Policecheck)
	// } else {
	// 	plainTextContent = fmt.Sprintf("Status: One or more tests  pending completion.\nIdentity Check Passed?%v\nAddress Check Passed?%v\nPolice Check Passed?%v,\n ", person.Idcheck, person.Addresscheck, person.Policecheck)

	// }

	// htmlContent := "<strong>GRAMACHECK</strong>"
	// message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	// client := sendgrid.NewSendClient(config.SendGridKey)
	// response, err := client.Send(message)
	// if err != nil {
	// 	log.Println("err", err)
	// } else {
	// 	fmt.Println(response.StatusCode)
	// 	fmt.Println(response.Body)
	// 	fmt.Println(response.Headers)
	// }
}

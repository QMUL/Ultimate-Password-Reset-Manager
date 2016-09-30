package prm

import (
	"log"
	"net/smtp"
	"strings"
)

// SendEmail uses smtp to post an email to a successful user at the end
// of the password change
func SendEmail(given_name string, email_address string, config *PRMConfig) {
	// Set up authentication information.
	auth := smtp.PlainAuth("", "user@example.com", "password", "localhost")

	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.

	body := config.EmailMsg
	body = strings.Replace(body, "%NAME%", given_name, -1)

	to := []string{email_address}
	msg := []byte("To: " + email_address + "\r\n" +
		"Subject: " + config.EmailSub + "\r\n" +
		"\r\n" +
		body)
	err := smtp.SendMail("localhost:25", auth, "its-research-support@qmul.ac.uk", to, msg)

	if err != nil {
		log.Print(err)
	}
}

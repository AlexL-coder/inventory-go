package handler

import (
	"awesomeProject1/services/grpc_auth"
	"encoding/json"
	"log"
	"net/smtp"
	"os"
	"strings"
)

// EmailDetails represents the email task details
type EmailDetails struct {
	Recipient string `json:"recipient"`
	Subject   string `json:"subject"`
	Body      string `json:"body"`
}

// HandleEmail Handle email-related tasks
func HandleEmail(msgBody []byte, queueName string, authService *grpc_auth.AuthServiceServer) error {
	var emailDetails EmailDetails
	// Deserialize the msgBody into the req struct
	if err := json.Unmarshal(msgBody, &emailDetails); err != nil {
		log.Printf("Failed to unmarshal message body: %v", err)
		return err
	}

	// Gmail SMTP server details
	smtpHost := os.Getenv("SMTPHOST")
	smtpPort := os.Getenv("SMTPPORT")
	username := os.Getenv("EMAIL_USERNAME")
	password := os.Getenv("EMAIL_PASSWORD")

	// Set up authentication
	auth := smtp.PlainAuth("", username, password, smtpHost)

	// Email content
	to := []string{emailDetails.Recipient}
	subject := "Subject: " + emailDetails.Subject + "\r\n"
	body := "From: " + username + "\r\n" +
		"To: " + strings.Join(to, ", ") + "\r\n" +
		subject +
		"\r\n" + emailDetails.Body

	log.Println(smtpHost+":"+smtpPort, auth, username)
	// Send the email
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, username, to, []byte(body))
	if err != nil {
		log.Printf("Failed to send email: %v", err)
		return err
	}

	log.Printf("Email sent successfully to: %s", emailDetails.Recipient)
	return nil
}

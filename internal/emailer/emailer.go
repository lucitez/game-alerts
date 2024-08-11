package emailer

import (
	"log/slog"
	"net/smtp"
	"os"
)

type Emailer struct {
	fromEmail string
	password  string
}

func New() Emailer {
	fromEmail := os.Getenv("SENDER_EMAIL")
	testEmail := os.Getenv("TEST_EMAIL")
	password := os.Getenv("SFTP_PASS")

	slog.Debug("Printing email", "from email", fromEmail)
	slog.Debug("Printing email", "test email", testEmail)

	return Emailer{
		fromEmail: fromEmail,
		password:  password,
	}
}

// https://pkg.go.dev/net/smtp#example-SendMail
func (e Emailer) SendEmail(toEmail, subject, body string) error {
	auth := smtp.PlainAuth("", e.fromEmail, e.password, "smtp.gmail.com")

	message := []byte("Subject: " + subject + "\r\n\r\n" + body)

	err := smtp.SendMail("smtp.gmail.com:587", auth, e.fromEmail, []string{toEmail}, message)
	return err
}

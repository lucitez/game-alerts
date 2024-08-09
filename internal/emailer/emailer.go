package emailer

import "net/smtp"

type Emailer struct {
	fromEmail string
	toEmail   string
	password  string
}

func New(fromEmail, toEmail, password string) Emailer {
	return Emailer{
		fromEmail: fromEmail,
		toEmail:   toEmail,
		password:  password,
	}
}

// https://pkg.go.dev/net/smtp#example-SendMail
func (e Emailer) SendEmail(subject, body string) error {
	auth := smtp.PlainAuth("", e.fromEmail, e.password, "smtp.gmail.com")

	message := []byte("Subject: " + subject + "\r\n\r\n" + body)

	err := smtp.SendMail("smtp.gmail.com:587", auth, e.fromEmail, []string{e.toEmail}, message)
	return err
}

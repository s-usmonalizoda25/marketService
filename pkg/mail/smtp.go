package mail

import (
	"context"
	"fmt"
	"net/smtp"
)

type SMTP struct {
	host, port, username, password, from string
}

func NewSMTP(host, port, username, password, from string) *SMTP {
	return &SMTP{host, port, username, password, from}
}

func (s *SMTP) SendOTP(ctx context.Context, email string, code string) error {
	auth := smtp.PlainAuth("", s.username, s.password, s.host)
	subject := "Subject: Email Verification\r\n"
	body := fmt.Sprintf("Your verification code is: %s\r\n\nThe code is valid for 10 minutes.", code)
	message := []byte(subject + "\r\n" + body)

	return smtp.SendMail(s.host+":"+s.port, auth, s.from, []string{email}, message)
}

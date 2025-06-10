package services

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/smtp"

	"financeAppAPI/internal/config"

	"github.com/go-mail/mail/v2"
)

type EmailService struct {
	dialer *mail.Dialer
	from   string
}

func NewEmailService() *EmailService {
	cfg := config.LoadConfig()
	port := 1025
	fmt.Sscanf(cfg.SMTPPort, "%d", &port)
	d := mail.NewDialer(cfg.SMTPHost, port, cfg.SMTPFrom, cfg.SMTPPass)
	d.TLSConfig = &tls.Config{
		ServerName:         cfg.SMTPHost,
		InsecureSkipVerify: false,
	}
	return &EmailService{dialer: d, from: cfg.SMTPFrom}
}

func (s *EmailService) Send(to, subject, body string) error {
	m := mail.NewMessage()
	m.SetHeader("From", s.from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)
	if err := s.dialer.DialAndSend(m); err != nil {
		log.Printf("SMTP error: %v", err)
		return fmt.Errorf("email sending failed")
	}
	return nil
}

func SendEmail(to, subject, body string) error {
	cfg := config.LoadConfig()
	from := cfg.SMTPFrom
	pass := cfg.SMTPPass
	host := cfg.SMTPHost
	port := cfg.SMTPPort
	addr := host + ":" + port
	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n\n" +
		body
	auth := smtp.PlainAuth("", from, pass, host)
	return smtp.SendMail(addr, auth, from, []string{to}, []byte(msg))
}

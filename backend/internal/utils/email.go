package utils

import (
	"fmt"
	"log"
	"net/smtp"
	"strings"
)

type emailConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

var emailCfg *emailConfig

func InitEmail(host string, port int, username, password, from string) {
	if host == "" {
		log.Println("SMTP not configured, email notifications disabled")
		return
	}
	emailCfg = &emailConfig{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
		From:     from,
	}
	log.Printf("SMTP configured: %s:%d", host, port)
}

func SendEmail(to, subject, htmlBody string) error {
	if emailCfg == nil {
		return fmt.Errorf("SMTP not configured")
	}

	addr := fmt.Sprintf("%s:%d", emailCfg.Host, emailCfg.Port)
	auth := smtp.PlainAuth("", emailCfg.Username, emailCfg.Password, emailCfg.Host)

	headers := map[string]string{
		"From":         emailCfg.From,
		"To":           to,
		"Subject":      subject,
		"MIME-Version": "1.0",
		"Content-Type": "text/html; charset=UTF-8",
	}

	var msg strings.Builder
	for k, v := range headers {
		msg.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	msg.WriteString("\r\n")
	msg.WriteString(htmlBody)

	err := smtp.SendMail(addr, auth, emailCfg.From, []string{to}, []byte(msg.String()))
	if err != nil {
		log.Printf("Failed to send email to %s: %v", to, err)
		return err
	}
	log.Printf("Email sent to %s: %s", to, subject)
	return nil
}

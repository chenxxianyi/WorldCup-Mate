package utils

import (
	"fmt"
	"log"
	"mime"
	"net/mail"
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
	if from == "" {
		from = username
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

	fromAddr, fromHeader, err := normalizeEmailAddress(emailCfg.From)
	if err != nil {
		return fmt.Errorf("invalid sender email address: %w", err)
	}
	toAddr, toHeader, err := normalizeEmailAddress(to)
	if err != nil {
		return fmt.Errorf("invalid recipient email address: %w", err)
	}

	addr := fmt.Sprintf("%s:%d", emailCfg.Host, emailCfg.Port)
	auth := smtp.PlainAuth("", emailCfg.Username, emailCfg.Password, emailCfg.Host)

	headers := map[string]string{
		"From":         fromHeader,
		"To":           toHeader,
		"Subject":      mime.QEncoding.Encode("UTF-8", subject),
		"MIME-Version": "1.0",
		"Content-Type": "text/html; charset=UTF-8",
	}

	var msg strings.Builder
	for k, v := range headers {
		msg.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	msg.WriteString("\r\n")
	msg.WriteString(htmlBody)

	err = smtp.SendMail(addr, auth, fromAddr, []string{toAddr}, []byte(msg.String()))
	if err != nil {
		log.Printf("Failed to send email to %s: %v", toAddr, err)
		return err
	}
	log.Printf("Email sent to %s: %s", toAddr, subject)
	return nil
}

func normalizeEmailAddress(addr string) (string, string, error) {
	addr = strings.TrimSpace(addr)
	if addr == "" || strings.ContainsAny(addr, "\r\n") {
		return "", "", fmt.Errorf("invalid email address")
	}
	parsed, err := mail.ParseAddress(addr)
	if err != nil {
		return "", "", err
	}
	return parsed.Address, parsed.String(), nil
}

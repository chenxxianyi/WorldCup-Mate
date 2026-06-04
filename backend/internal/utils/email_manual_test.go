package utils

import (
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/joho/godotenv"
)

func TestManualEmailSend(t *testing.T) {
	if os.Getenv("SMTP_MANUAL_TEST") != "1" {
		t.Skip("set SMTP_MANUAL_TEST=1 to send a real email")
	}

	_ = godotenv.Load("../../.env")

	host := os.Getenv("SMTP_HOST")
	username := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")
	from := os.Getenv("SMTP_FROM")
	to := os.Getenv("SMTP_TEST_TO")
	if to == "" {
		to = username
	}

	if host == "" || username == "" || password == "" || from == "" || to == "" {
		t.Fatal("SMTP_HOST, SMTP_USERNAME, SMTP_PASSWORD, SMTP_FROM, and recipient email are required")
	}
	port, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil || port <= 0 {
		port = 587
	}

	InitEmail(host, port, username, password, from)
	subject := fmt.Sprintf("WorldCup Mate email test %s", time.Now().Format(time.RFC3339))
	body := `<div style="font-family:sans-serif;padding:16px;">
<h2>WorldCup Mate email test</h2>
<p>If you received this, SMTP notification delivery is working.</p>
</div>`
	if err := SendEmail(to, subject, body); err != nil {
		t.Fatalf("send email: %v", err)
	}
}

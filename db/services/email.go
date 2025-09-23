package services

import (
	"fmt"
	"os"
	"strconv"

	gomail "gopkg.in/gomail.v2"
)

func AdminEmail() string {
	admin := os.Getenv("ADMIN_EMAIL")
	if admin == "" {
		return "tayousman17@example.com"
	}
	return admin
}

func SendEmail(to, subject, body string) error {
	host := os.Getenv("SMTP_HOST")
	portStr := os.Getenv("SMTP_PORT")
	user := os.Getenv("SMTP_USERNAME")
	pass := os.Getenv("SMTP_PASSWORD")
	if host == "" || portStr == "" || user == "" || pass == "" {
		return fmt.Errorf("smtp not configured")
	}
	port, _ := strconv.Atoi(portStr)

	m := gomail.NewMessage()
	from := user
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	d := gomail.NewDialer(host, port, user, pass)
	return d.DialAndSend(m)
}

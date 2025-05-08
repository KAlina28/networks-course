package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"os"
	"strings"
	"time"

	"gopkg.in/gomail.v2"
)

const (
	senderName = "client"
)

var (
	user     string
	password string
)

type EmailData struct {
	Email      string
	Date       string
	SenderName string
}

func sendEmail(recipient, body, contentType string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", fmt.Sprintf("%s <%s>", senderName, "krylovaalina2004@gmail.com"))
	m.SetHeader("To", recipient)

	if contentType == "html" {
		m.SetBody("text/html", body)
	} else {
		m.SetBody("text/plain", body)
	}

	d := gomail.NewDialer("smtp.gmail.com", 587, user, password)
	d.SSL = false

	return d.DialAndSend(m)
}

func main() {
	user = os.Getenv("SMTP_USER")
	password = os.Getenv("SMTP_PASSWORD")

	if user == "" || password == "" {
		log.Fatal("SMTP_USER и SMTP_PASSWORD")
	}

	recipient := flag.String("to", "", "Почта")
	format := flag.String("format", "txt", "txt или html")
	flag.Parse()

	if *recipient == "" {
		fmt.Println("адрес -to")
		os.Exit(1)
	}

	var body string
	var err error

	if *format == "html" {
		tmpl, err := template.ParseFiles("example.html")
		if err != nil {
			log.Fatalf("%v", err)
		}

		data := EmailData{
			Email:      *recipient,
			Date:       time.Now().Format("02.01.2006 15:04:05"),
			SenderName: senderName,
		}

		buffer := &strings.Builder{}
		if err := tmpl.Execute(buffer, data); err != nil {
			log.Fatalf("%v", err)
		}
		body = buffer.String()
	} else {
		body = fmt.Sprintf(`Привет. Как дела ?(txt) С уважением, %s`, senderName)
	}

	err = sendEmail(*recipient, body, *format)
	if err != nil {
		log.Fatalf("ошибка - %v", err)
	}

	fmt.Printf("отправлено на %s - %s\n", *recipient, *format)
}

package models

import (
	"bytes"
	"fmt"
	"log"
	"net/smtp"
	"text/template"

	"github.com/liquiloans/sftp/config"
)

var auth smtp.Auth

type Request struct {
	subject string
	body    string
}

func NewRequest(subject string, body string) *Request {
	return &Request{
		subject: subject,
		body:    body,
	}
}

func (r *Request) SendEmail() (bool, error) {
	auth = smtp.PlainAuth("", config.EMAIL_USERNAME, config.EMAIL_PASSWORD, config.EMAIL_HOST)
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

	msg := []byte("Subject:" + r.subject + "\n" + mime + "\n\n" + "From:" + config.FromEmails + "\n\n" + r.body)
	fmt.Println(msg)
	addr := fmt.Sprintf("%s:%d", config.EMAIL_HOST, 586)
	if err := smtp.SendMail(addr, auth, config.FromEmails, config.ToEmails, msg); err != nil {
		log.Fatal("SendMail Error Message", err)
	}
	return true, nil
}

func (r *Request) ParseTemplate(templateFileName string, data interface{}) error {
	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		log.Fatal("Template parse error", err)
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return err
	}
	r.body = buf.String()
	return nil
}

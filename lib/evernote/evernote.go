package evernote

import (
	"log"
	"net/smtp"
)

type Sender struct {
	email        string
	gmailAccount string
	gmailPass    string
}

func NewSenderFromMap(config map[interface{}]interface{}) (sender *Sender) {
	evernote := config["evernote_mail"].(string)
	gmail := config["gmail_account"].(string)
	pass := config["gmail_pass"].(string)
	return &Sender{
		email:        evernote,
		gmailAccount: gmail,
		gmailPass:    pass,
	}
}

func (sender *Sender) SendNote(title string, text string) {
	body := "To: " + sender.email + "\r\nSubject: " +
	title + "\r\n\r\n" + text
	auth := smtp.PlainAuth("",sender.gmailAccount,sender.gmailPass,"smtp.gmail.com")
	err := smtp.SendMail("smtp.gmail.com:587",auth,sender.gmailAccount,
	[]string{sender.email},[]byte(body))
	if err != nil {
		log.Fatal(err)
	}
}

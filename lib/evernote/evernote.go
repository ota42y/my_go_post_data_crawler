package evernote

import (
	"log"
	"net/smtp"
	"gopkg.in/yaml.v2"
)

type Sender struct {
	EvernoteMail string
	GmailAccount string
	GmailPass    string
}

func NewSenderFromData(data []byte) (s *Sender) {
	s = &Sender{}
	err := yaml.Unmarshal(data, s)
	if err != nil {
		s = nil
	}
	return
}

func (sender *Sender) SendNote(title string, text string) {
	body := "To: " + sender.EvernoteMail + "\r\nSubject: " +
	title + "\r\n\r\n" + text
	auth := smtp.PlainAuth("",sender.GmailAccount,sender.GmailPass,"smtp.gmail.com")
	err := smtp.SendMail("smtp.gmail.com:587",auth,sender.GmailAccount,
	[]string{sender.EvernoteMail},[]byte(body))
	if err != nil {
		log.Fatal(err)
	}
}

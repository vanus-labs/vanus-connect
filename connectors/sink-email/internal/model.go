package internal

import (
	stdMail "net/mail"
	"strings"

	"github.com/nikoksr/notify/service/mail"
	"github.com/pkg/errors"
)

type EmailMessage struct {
	Subject  string   `json:"subject"`
	Body     string   `json:"body"`
	Receiver string   `json:"receiver"`
	BodyType string   `json:"body_type"`
	To       []string `json:"-"`
}

func (e *EmailMessage) GetBodyType() mail.BodyType {
	if e.BodyType == "" {
		return mail.PlainText
	}
	bodyType := strings.ToLower(e.BodyType)
	if bodyType == "html" {
		return mail.HTML
	}
	return mail.PlainText
}

func (e *EmailMessage) Validate() error {
	if e.Subject == "" {
		return errors.New("email subject is empty")
	}
	if e.Body == "" {
		return errors.New("email body is empty")
	}
	if e.Receiver == "" {
		return errors.New("email receiver is empty")
	}
	adders, err := stdMail.ParseAddressList(e.Receiver)
	if err != nil {
		return errors.Wrapf(err, "failed to parse receiver address %s", e.Receiver)
	}
	if len(adders) == 0 {
		return errors.New("receiver address is empty")
	}
	for i := range adders {
		e.To = append(e.To, adders[i].Address)
	}
	return nil
}

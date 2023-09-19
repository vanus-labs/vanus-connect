package internal

import (
	stdMail "net/mail"
	"strings"

	"github.com/nikoksr/notify/service/mail"
	"github.com/pkg/errors"
)

type EmailMessage struct {
	Subject  string   `json:"subject"`
	To       string   `json:"to"`
	Body     string   `json:"body"`
	BodyType string   `json:"body_type"`
	ToAdders []string `json:"-"`
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
	if e.To == "" {
		return errors.New("email to is empty")
	}
	adders, err := stdMail.ParseAddressList(e.To)
	if err != nil {
		return errors.Wrapf(err, "failed to parse to address %s", e.To)
	}
	if len(adders) == 0 {
		return errors.New("to address is empty")
	}
	for i := range adders {
		e.ToAdders = append(e.ToAdders, adders[i].Address)
	}
	return nil
}

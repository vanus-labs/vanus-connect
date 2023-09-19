package internal

import (
	"github.com/pkg/errors"
)

type EmailMessage struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

func (e *EmailMessage) Validate() error {
	if e.To == "" {
		return errors.New("email to is empty")
	}
	if e.Subject == "" {
		return errors.New("email subject is empty")
	}
	if e.Body == "" {
		return errors.New("email body is empty")
	}
	return nil
}

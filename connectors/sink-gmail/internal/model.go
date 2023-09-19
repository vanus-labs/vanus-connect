package internal

import (
	"github.com/pkg/errors"
)

type EmailMessage struct {
	Subject  string `json:"subject"`
	Body     string `json:"body"`
	Receiver string `json:"receiver"`
	Sender   string `json:"-"`
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
	return nil
}

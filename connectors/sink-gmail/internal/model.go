package internal

import (
	"github.com/pkg/errors"
)

type EmailMessage struct {
	Subject    string `json:"subject"`
	Body       string `json:"body"`
	Recipients string `json:"recipients"`
}

func (e *EmailMessage) Validate() error {
	if e.Subject == "" {
		return errors.New("email subject is empty")
	}
	if e.Body == "" {
		return errors.New("email body is empty")
	}
	if e.Recipients == "" {
		return errors.New("email recipients is empty")
	}
	return nil
}

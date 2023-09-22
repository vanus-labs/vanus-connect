package internal

import (
	"strings"

	"github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/pkg/errors"
)

type EmailMessage struct {
	To       string `json:"to"`
	Subject  string `json:"subject"`
	Body     string `json:"body"`
	BodyType string `json:"body_type"`
}

func (e *EmailMessage) GetBodyType() models.BodyType {
	if e.BodyType == "" {
		return models.TEXT_BODYTYPE
	}
	bodyType := strings.ToLower(e.BodyType)
	if bodyType == "html" {
		return models.HTML_BODYTYPE
	}
	return models.TEXT_BODYTYPE
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

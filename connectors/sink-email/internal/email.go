package internal

import (
	"context"
	"fmt"
	"strings"
	"time"

	ce "github.com/cloudevents/sdk-go/v2"
	"github.com/nikoksr/notify/service/mail"
)

func getEmailBodyType(bodyType string) mail.BodyType {
	bodyType = strings.ToLower(bodyType)
	if bodyType == "html" {
		return mail.HTML
	}
	return mail.PlainText
}

func (e *emailSink) getEmailBodyType(event *ce.Event) mail.BodyType {
	format, exist := event.Extensions()[xvEmailFormat].(string)
	if exist && format != "" {
		return getEmailBodyType(format)
	}
	return e.defaultEmailFormat
}

func (e *emailSink) getEmailAddress(event *ce.Event) string {
	address, exist := event.Extensions()[xvEmailFrom].(string)
	if exist && address != "" {
		return address
	}
	return e.defaultAccount
}

func (e *emailSink) send(ctx context.Context, em *EmailMessage) error {
	_ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	cfg := e.mails[em.Sender]
	port := cfg.Port
	if port == 0 {
		port = 25
	}
	m := mail.New(cfg.Account, fmt.Sprintf("%s:%d", cfg.Host, port))
	m.AuthenticateSMTP(cfg.Identity, cfg.Account, cfg.Password, cfg.Host)
	m.AddReceivers(em.To...)
	m.BodyFormat(em.Type)
	return m.Send(_ctx, em.Subject, em.Body)
}

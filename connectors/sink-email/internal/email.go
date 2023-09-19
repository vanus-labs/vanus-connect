package internal

import (
	"context"
	"fmt"
	"time"

	"github.com/nikoksr/notify/service/mail"
)

func (e *emailSink) send(ctx context.Context, em *EmailMessage) error {
	_ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	cfg := e.emailCfg
	port := cfg.Port
	if port == 0 {
		port = 25
	}
	m := mail.New(cfg.Account, fmt.Sprintf("%s:%d", cfg.Host, port))
	m.AuthenticateSMTP(cfg.Identity, cfg.Account, cfg.Password, cfg.Host)
	m.AddReceivers(em.To...)
	m.BodyFormat(em.GetBodyType())
	return m.Send(_ctx, em.Subject, em.Body)
}

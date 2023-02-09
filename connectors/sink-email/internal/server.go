// Copyright 2022 Linkall Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package internal

import (
	"context"
	"fmt"
	"net/http"
	stdMail "net/mail"
	"strings"
	"time"

	v2 "github.com/cloudevents/sdk-go/v2"
	cdkgo "github.com/linkall-labs/cdk-go"
	"github.com/linkall-labs/cdk-go/log"
	"github.com/nikoksr/notify/service/mail"
	"github.com/pkg/errors"
)

const (
	name = "Email Sink"
	// must-have for incoming event
	xvEmailSubject = "xvemailsubject"
	// must-have for incoming event
	xvEmailRecipients = "xvemailrecipients"
	xvEmailFrom       = "xvemailfrom"
	xvEmailFormat     = "xvemailformat"
)

var (
	errInvalidFromAddress = cdkgo.NewResult(http.StatusBadRequest,
		"email: invalid or empty email from address")
	errInvalidSubject = cdkgo.NewResult(http.StatusBadRequest,
		"email: subject can't be empty")
	errInvalidRecipients = cdkgo.NewResult(http.StatusBadRequest,
		"email: invalid or empty recipients")
	errFromAddressNotInConfiguration = cdkgo.NewResult(http.StatusBadRequest,
		"email: the email from address not found in configuration")
	errFailedToSend = cdkgo.NewResult(http.StatusInternalServerError,
		"email: failed to sent email, please view logs")
)
var _ cdkgo.SinkConfigAccessor = &emailConfig{}

type EmailConfig struct {
	Account  string `json:"account" yaml:"account" validate:"required,email"`
	Password string `json:"password" yaml:"password" validate:"required"`
	Host     string `json:"host" yaml:"host" validate:"required"`
	Port     int    `json:"port" yaml:"port"`
	Format   string `json:"format" yaml:"format"`
	Identity string `json:"identity" yaml:"identity"`
}

type emailConfig struct {
	cdkgo.SinkConfig `json:",inline" yaml:",inline"`
	DefaultAddress   string        `json:"default" yaml:"default"`
	Emails           []EmailConfig `json:"email" yaml:"email" validate:"required,gt=0,dive"`
}

func (c *emailConfig) Validate() error {
	if c.DefaultAddress != "" {
		var exist bool
		for _, email := range c.Emails {
			if email.Account == c.DefaultAddress {
				exist = true
				break
			}
		}
		if !exist {
			return errors.New("email: the default email config isn't exist")
		}
	}
	return c.SinkConfig.Validate()
}

func NewConfig() cdkgo.SinkConfigAccessor {
	return &emailConfig{}
}

func NewEmailSink() cdkgo.Sink {
	return &emailSink{
		mails: map[string]EmailConfig{},
	}
}

var _ cdkgo.Sink = &emailSink{}

type emailSink struct {
	count          int64
	mails          map[string]EmailConfig
	defaultAccount string
}

func (e *emailSink) Arrived(ctx context.Context, events ...*v2.Event) cdkgo.Result {
	for idx := range events {
		event := events[idx]
		var address string
		val, exist := event.Extensions()[xvEmailFrom]
		if exist {
			_address, ok := val.(string)
			if !ok {
				return errInvalidFromAddress
			}
			address = _address
		} else {
			address = e.defaultAccount
		}
		if address == "" {
			return errInvalidFromAddress
		}

		cfg, exist := e.mails[address]
		if !exist {
			return errFromAddressNotInConfiguration
		}

		val, exist = event.Extensions()[xvEmailSubject]
		if !exist {
			return errInvalidSubject
		}
		subject := val.(string)
		if subject == "" {
			return errInvalidSubject
		}

		val, exist = event.Extensions()[xvEmailRecipients]
		if !exist {
			return errInvalidRecipients
		}
		recipients := val.(string)
		if recipients == "" {
			return errInvalidRecipients
		}
		addrs, err := stdMail.ParseAddressList(recipients)
		if err != nil {
			log.Error("failed to parse receiver address", map[string]interface{}{
				log.KeyError: err,
				"recipients": recipients,
				"event_id":   event.ID(),
			})
			return errInvalidRecipients
		}

		var to = make([]string, len(addrs))
		for idx, v := range addrs {
			if v.Address != "" {
				to[idx] = v.Address
			}
		}
		if len(to) == 0 {
			return errInvalidRecipients
		}

		var format string
		val, exist = event.Extensions()[xvEmailFormat]
		if exist {
			format = val.(string)
			if format == "" {
				format = cfg.Format
			}
		}

		var f mail.BodyType
		switch strings.ToLower(format) {
		case "html":
			f = mail.HTML
		default:
			f = mail.PlainText
		}

		start := time.Now()
		if err := e.send(ctx, cfg, subject, event.Data(), f, to...); err != nil {
			log.Error("failed to send email", map[string]interface{}{
				log.KeyError: err,
				"address":    address,
				"event_id":   event.ID(),
			})
			return errFailedToSend
		} else if time.Now().Sub(start) > time.Second {
			log.Debug("success to send email, but takes too long", map[string]interface{}{
				"address":   address,
				"event_id":  event.ID(),
				"used_time": time.Now().Sub(start),
			})
		} else {
			log.Debug("success to send email", map[string]interface{}{
				"address":  address,
				"event_id": event.ID(),
			})
		}
	}
	return cdkgo.SuccessResult
}

func (e *emailSink) Initialize(_ context.Context, cfg cdkgo.ConfigAccessor) error {
	config, ok := cfg.(*emailConfig)
	if !ok {
		return errors.New("email: invalid configuration type")
	}

	for _, m := range config.Emails {
		e.mails[m.Account] = m
	}
	e.defaultAccount = config.DefaultAddress
	return nil
}

func (e *emailSink) Name() string {
	return name
}

func (e *emailSink) Destroy() error {
	// nothing to do
	return nil
}

func (e *emailSink) send(ctx context.Context, cfg EmailConfig,
	subject string, message []byte, f mail.BodyType, to ...string) error {
	_ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	port := cfg.Port
	if port == 0 {
		port = 25
	}
	m := mail.New(cfg.Account, fmt.Sprintf("%s:%d", cfg.Host, port))
	m.AuthenticateSMTP(cfg.Identity, cfg.Account, cfg.Password, cfg.Host)
	m.AddReceivers(to...)
	m.BodyFormat(f)
	return m.Send(_ctx, subject, string(message))
}

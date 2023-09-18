package internal

import (
	"github.com/pkg/errors"

	cdkgo "github.com/vanus-labs/cdk-go"
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
	DefaultAccount   string        `json:"default" yaml:"default"`
	Emails           []EmailConfig `json:"email" yaml:"email" validate:"dive"`
}

func (c *emailConfig) Validate() error {
	if len(c.Emails) == 0 {
		return errors.New("email length is 0")
	}
	if c.DefaultAccount == "" {
		c.DefaultAccount = c.Emails[0].Account
	} else {
		var exist bool
		for _, email := range c.Emails {
			if email.Account == c.DefaultAccount {
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

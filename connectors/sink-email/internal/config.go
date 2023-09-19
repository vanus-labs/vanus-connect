package internal

import (
	cdkgo "github.com/vanus-labs/cdk-go"
)

var _ cdkgo.SinkConfigAccessor = &emailConfig{}

type EmailConfig struct {
	Account  string `json:"account" yaml:"account" validate:"required,email"`
	Password string `json:"password" yaml:"password" validate:"required"`
	Host     string `json:"host" yaml:"host" validate:"required"`
	Port     int    `json:"port" yaml:"port"`
	Identity string `json:"identity" yaml:"identity"`
}

type emailConfig struct {
	cdkgo.SinkConfig `json:",inline" yaml:",inline"`
	Email            EmailConfig `json:"email" yaml:"email" validate:"required"`
}

func NewConfig() cdkgo.SinkConfigAccessor {
	return &emailConfig{}
}

package internal

import (
	cdk "github.com/vanus-labs/cdk-go"
)

var _ cdk.SourceConfigAccessor = &Config{}

type Config struct {
	cdk.SourceConfig `json:",inline" yaml:",inline"`

	Port              int           `json:"port" yaml:"port"`
	AppID             string        `json:"feishu_app_id" yaml:"feishu_app_id" validate:"required"`
	AppSecret         string        `json:"feishu_app_secret" yaml:"feishu_app_secret" validate:"required"`
	VerificationToken string        `json:"verification_token" yaml:"verification_token" validate:"required"`
	EventEncryptKey   string        `json:"event_encrypt_key" yaml:"event_encrypt_key"`
	EventType         []MessageType `json:"event_type" yaml:"event_type"`
}

func NewConfig() cdk.SourceConfigAccessor {
	return &Config{}
}

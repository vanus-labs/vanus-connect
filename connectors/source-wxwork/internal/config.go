package internal

import (
	cdkgo "github.com/vanus-labs/cdk-go"
)

var _ cdkgo.SourceConfigAccessor = &Config{}

type Config struct {
	cdkgo.SourceConfig `json:",inline" yaml:",inline"`

	VanusAIURL   string `json:"vanus_ai_url" yaml:"vanus_ai_url"`
	VanusAIAppId string `json:"vanus_ai_app_id" yaml:"vanus_ai_app_id" validate:"required"`

	WeworkCorpId         string `json:"wework_corp_id" yaml:"wework_corp_id" validate:"required"`
	WeworkAppId          int    `json:"wework_app_id" yaml:"wework_app_id" validate:"required"`
	WeworkAppSecret      string `json:"wework_app_secret" yaml:"wework_app_secret" validate:"required"`
	WeworkToken          string `json:"wework_token" yaml:"wework_token" validate:"required"`
	WeworkEncodingAESKey string `json:"wework_encoding_aes_key" yaml:"wework_encoding_aes_key" validate:"required"`
}

func NewConfig() cdkgo.SourceConfigAccessor {
	return &Config{}
}

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
	WeworkAgentId        int64  `json:"wework_agent_id" yaml:"wework_agent_id" validate:"required"`
	WeworkAgentSecret    string `json:"wework_agent_secret" yaml:"wework_agent_secret" validate:"required"`
	WeworkToken          string `json:"wework_token" yaml:"wework_token" validate:"required"`
	WeworkEncodingAESKey string `json:"wework_encoding_aes_key" yaml:"wework_encoding_aes_key" validate:"required"`
}

func NewConfig() cdkgo.SourceConfigAccessor {
	return &Config{}
}

func (c *Config) Init() {

}

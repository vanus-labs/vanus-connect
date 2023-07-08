package internal

import (
	cdkgo "github.com/vanus-labs/cdk-go"
)

var _ cdkgo.SourceConfigAccessor = &Config{}

type Config struct {
	cdkgo.SourceConfig `json:",inline" yaml:",inline"`

	VanusAIURL   string `json:"vanus_ai_url" yaml:"vanus_ai_url" validate:"required"`
	VanusAIAppId string `json:"vanus_ai_app_id" yaml:"vanus_ai_app_id" validate:"required"`

	DingtalkAppKey    string `json:"dingtalk_app_key" yaml:"dingtalk_app_key" validate:"required"`
	DingtalkAppSecret string `json:"dingtalk_app_secret" yaml:"dingtalk_app_secret" validate:"required"`
}

func NewConfig() cdkgo.SourceConfigAccessor {
	return &Config{}
}

func (c *Config) Init() {

}

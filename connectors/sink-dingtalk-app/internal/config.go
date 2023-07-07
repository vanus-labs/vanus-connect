package internal

import (
	cdkgo "github.com/vanus-labs/cdk-go"
)

var _ cdkgo.SinkConfigAccessor = &Config{}

type Config struct {
	cdkgo.SinkConfig `json:",inline" yaml:",inline"`

	DingtalkAppKey    string `json:"dingtalk_app_key" yaml:"dingtalk_app_key"`
	DingtalkAppSecret string `json:"dingtalk_app_secret" yaml:"dingtalk_app_secret"`
}

func NewConfig() cdkgo.SinkConfigAccessor {
	return &Config{}
}

func (c *Config) Init() {

}

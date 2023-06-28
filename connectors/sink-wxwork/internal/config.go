package internal

import (
	cdkgo "github.com/vanus-labs/cdk-go"
)

var _ cdkgo.SinkConfigAccessor = &Config{}

type Config struct {
	cdkgo.SinkConfig `json:",inline" yaml:",inline"`

	WeworkCorpId      string `json:"wework_corp_id" yaml:"wework_corp_id" validate:"required"`
	WeworkAgentId     int64  `json:"wework_agent_id" yaml:"wework_agent_id" validate:"required"`
	WeworkAgentSecret string `json:"wework_agent_secret" yaml:"wework_agent_secret" validate:"required"`
}

func NewConfig() cdkgo.SinkConfigAccessor {
	return &Config{}
}

func (c *Config) Init() {

}

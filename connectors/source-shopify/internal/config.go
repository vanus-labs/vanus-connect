package internal

import cdkgo "github.com/vanus-labs/cdk-go"

var _ cdkgo.SourceConfigAccessor = &shopifySourceConfig{}

type shopifySourceConfig struct {
	cdkgo.SourceConfig `json:",inline" yaml:",inline"`
	Port               int    `json:"port" yaml:"port"`
	ClientSecret       string `json:"client_secret" yaml:"client_secret"`
}

func NewConfig() cdkgo.SourceConfigAccessor {
	return &shopifySourceConfig{}
}

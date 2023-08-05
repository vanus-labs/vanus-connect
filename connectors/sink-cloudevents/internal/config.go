package internal

import (
	cdkgo "github.com/vanus-labs/cdk-go"
)

func NewConfig() cdkgo.SinkConfigAccessor {
	return &config{}
}

type config struct {
	cdkgo.SinkConfig `json:",inline" yaml:",inline"`

	Target string `json:"target" yaml:"target"`
}

package main

import (
	cdkgo "github.com/vanus-labs/cdk-go"
	"github.com/vanus-labs/sink-dingtalk-app/internal"
)

func main() {
	cdkgo.RunSink(internal.NewConfig, internal.NewSink)
}

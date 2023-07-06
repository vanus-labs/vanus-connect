package main

import (
	cdkgo "github.com/vanus-labs/cdk-go"
	"github.com/vanus-labs/source-dingtalk-app/internal"
)

func main() {
	cdkgo.RunSource(internal.NewConfig, internal.Source)
}

package main

import (
	cdk "github.com/vanus-labs/cdk-go"
	"github.com/vanus-labs/source-feishu-app/internal"
)

func main() {
	cdk.RunHttpSource(internal.NewConfig, internal.NewSource)
}

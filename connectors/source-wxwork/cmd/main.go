package main

import (
	cdkgo "github.com/vanus-labs/cdk-go"
	"github.com/vanus-labs/source-wxwork/internal"
)

func main() {
	cdkgo.RunHttpSource(internal.NewConfig, internal.NewSource)
}

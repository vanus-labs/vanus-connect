package main

import (
	cdkgo "github.com/vanus-labs/cdk-go"
	"github.com/vanus-labs/sink-wxwork/internal"
)

func main() {
	cdkgo.RunSink(internal.NewConfig, internal.NewSink)
}

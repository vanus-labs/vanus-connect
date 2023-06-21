package internal

import (
	"context"
	cdkgo "github.com/vanus-labs/cdk-go"
	"os"
	"testing"
)

func TestDouyin(t *testing.T) {
	var err error

	s := Source()
	c := &DouyinConfig{
		SourceConfig: cdkgo.SourceConfig{
			Target: "http://localhost:9191",
		},
		AuthCode:     os.Getenv("AuthCode"),
		ClientKey:    os.Getenv("ClientKey"),
		ClientSecret: os.Getenv("ClientSecret"),
	}
	err = c.Validate()
	if err != nil {
		println("c.Validate error, ", err.Error())
		return
	}
	err = s.Initialize(context.Background(), c)
	if err != nil {
		println("s.Initialize error, ", err.Error())
		return
	}
}

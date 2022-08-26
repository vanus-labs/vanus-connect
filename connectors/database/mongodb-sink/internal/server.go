// Copyright 2022 Linkall Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package internal

import (
	"bytes"
	"context"
	"fmt"
	v2 "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/protocol"
	"github.com/golang/protobuf/jsonpb"
	proto "github.com/linkall-labs/connector/mongodb-sink/database"
	"net"

	"github.com/cloudevents/sdk-go/v2/client"
	cehttp "github.com/cloudevents/sdk-go/v2/protocol/http"
)

type Config struct {
	Port int
}

type Secret struct {
}

type sink struct {
	cfg *Config
}

func NewMongoSink(config Config) *sink {
	return &sink{
		cfg: &config,
	}
}

func (s *sink) StartReceive(ctx context.Context) error {
	ls, err := net.Listen("tcp", fmt.Sprintf(":%d", s.cfg.Port))
	if err != nil {
		return err
	}

	c, err := client.NewHTTP(cehttp.WithListener(ls), cehttp.WithRequestDataAtContextMiddleware())
	if err != nil {
		return err
	}
	return c.StartReceiver(ctx, s.receive)
}

func (s *sink) receive(ctx context.Context, event v2.Event) protocol.Result {
	e := &proto.Event{}
	err := jsonpb.Unmarshal(bytes.NewReader(event.Data()), e)
	println(err)
	d, _ := event.MarshalJSON()
	println(string(d))
	return cehttp.NewResult(200, "")
}

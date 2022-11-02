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

package main

import (
	"context"
	"fmt"
	ce "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
	"github.com/linkall-labs/connector/sink/tencent-cloud/cos/internal"
	"github.com/tencentyun/scf-go-lib/cloudfunction"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	eventType  = "tencent-cloud-cos-event"
	httpPrefix = "http://"
)

type Records struct {
	Records []map[string]interface{} `json:"Records"`
}

type wrapper struct {
	target   string
	ceClient ce.Client
	funcName string
	debug    bool
	eventbus string
}

func (w *wrapper) receive(ctx context.Context, records Records) error {
	for _, v := range records.Records {
		event := ce.NewEvent()
		event.SetID(uuid.New().String())
		event.SetType(eventType)
		event.SetSource(w.funcName)
		event.SetTime(time.Now())
		err := event.SetData(ce.ApplicationJSON, v)
		if err != nil {
			fmt.Printf("failed to set event data: %v, raw: %v\n", err, v)
			continue
		}

		res := w.ceClient.Send(ctx, event)
		if !ce.IsACK(res) {
			fmt.Printf("failed to send event: %s\n", res.Error())
		} else if w.debug {
			fmt.Printf("success to send event: %v\n", v)
		}
	}
	return nil
}

func main() {
	endpoint := os.Getenv(internal.EnvEventGateway)
	if endpoint == "" {
		panic("event gateway can't be empty")
	}
	fName := os.Getenv(internal.EnvFuncName)
	if endpoint == "" {
		panic("function name can't be empty")
	}
	eventbus := os.Getenv(internal.EnvVanusEventbus)
	if eventbus == "" {
		panic("eventbus can't be empty")
	}

	var target string
	if strings.HasPrefix(endpoint, httpPrefix) {
		target = fmt.Sprintf("%s/gateway/%s", endpoint, eventbus)
	} else {
		target = fmt.Sprintf("%s%s/gateway/%s", httpPrefix, endpoint, eventbus)
	}
	dg, _ := strconv.ParseBool(os.Getenv(internal.EnvDebugMode))
	cli, err := ce.NewClientHTTP(ce.WithTarget(target))
	if err != nil {
		panic("failed to init cloudevents client")
	}

	w := &wrapper{
		eventbus: eventbus,
		target:   endpoint,
		ceClient: cli,
		funcName: fName,
		debug:    dg,
	}
	// Make the handler available for Remote Procedure Call by Cloud Function
	cloudfunction.Start(w.receive)
}

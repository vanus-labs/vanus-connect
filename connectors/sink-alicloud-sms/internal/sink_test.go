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
	"context"
	"fmt"
	"os"
	"testing"

	v2 "github.com/cloudevents/sdk-go/v2"
)

var dataFmt = "{\"specversion\":\"1.0\",\"id\":\"55545b87-049d-48f2-a88c-1630915d6dd8\",\"source\":\"/debezium/mysql/test\"" +
	",\"type\":\"io.debezium.mysql.datachangeevent\",\"datacontenttype\":\"application/json\",\"time\":\"" +
	"2022-12-28T01:23:43Z\",\"data\":{\"phones\":\"%s\"},\"iodebeziumrow\":\"0\",\"xvanusstime\":\"2022-12-28T0" +
	"1:23:43.714Z\",\"iodebeziumpos\":\"1073\",\"iodebeziumsnapshot\":\"false\",\"iodebeziumname\":\"test\",\"iodeb" +
	"eziumtsms\":\"1672190623000\",\"iodebeziumconnector\":\"mysql\",\"xvanuslogoffset\":\"AAAAAAAAAAQ=\",\"iodebez" +
	"iumthread\":\"696\",\"iodebeziumversion\":\"2.0.1.Final\",\"iodebeziumfile\":\"mysql-bin-changelog.000392\",\"i" +
	"odebeziumtable\":\"test\",\"iodebeziumdb\":\"linkall\",\"iodebeziumserverid\":\"1220090578\",\"xvanuseventbus\"" +
	":\"demo\",\"iodebeziumop\":\"c\"}"

func TestSMSSink(t *testing.T) {
	s := NewSink()
	c := &smsConfig{
		aliConfig: aliConfig{
			AccessKeyId:     os.Getenv("AccessKeyId"),
			AccessKeySecret: os.Getenv("AccessKeySecret"),
			SignName:        "Vanus",
			TemplateCode:    "SMS_280041016",
			TemplateParam:   "{\"code\":\"919191\"}",
		},
	}
	_ = c.Validate()
	_ = s.Initialize(context.Background(), c)

	phones := os.Getenv("Phones")
	data := fmt.Sprintf(dataFmt, phones)

	fmt.Println(data)

	var events []*v2.Event
	e := v2.NewEvent()
	_ = e.UnmarshalJSON([]byte(data))
	events = append(events, &e)

	s.Arrived(context.Background(), events...)
}

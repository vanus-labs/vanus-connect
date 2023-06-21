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
	"encoding/json"
	"errors"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	v2 "github.com/cloudevents/sdk-go/v2"
	"github.com/tidwall/gjson"
)

type TemplateKV struct {
	Key   string `json:"key" yaml:"key" validate:"required"`
	Value string `json:"value" yaml:"value" validate:"required"`
}

type aliConfig struct {
	AccessKeyId     string       `json:"access_key_id" yaml:"access_key_id" validate:"required"`
	AccessKeySecret string       `json:"access_key_secret" yaml:"access_key_secret" validate:"required"`
	SignName        string       `json:"sign_name" yaml:"sign_name" validate:"required"`
	PhoneNumbers    string       `json:"phone_numbers" yaml:"phone_numbers" validate:"required"`
	TemplateCode    string       `json:"template_code" yaml:"template_code" validate:"required"`
	TemplateParam   []TemplateKV `json:"template_param" yaml:"template_param"`
}

const (
	UnfixedKeyPrefix = "$."
)

type aliSMS struct {
	cfg    *aliConfig
	client *dysmsapi.Client
}

func (sms *aliSMS) init(cfg aliConfig) (err error) {
	sms.cfg = &cfg
	sms.client, err = dysmsapi.NewClientWithAccessKey("", sms.cfg.AccessKeyId, sms.cfg.AccessKeySecret)
	if err != nil {
		return err
	}

	return nil
}

func (sms *aliSMS) sendMsg(e *v2.Event) (err error) {
	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = "https"
	request.SignName = sms.cfg.SignName
	request.TemplateCode = sms.cfg.TemplateCode
	request.PhoneNumbers = sms.getPhones(e)
	if param := sms.getTemplateParam(e); param != "" {
		request.TemplateParam = param
	}

	resp, err := sms.client.SendSms(request)
	if err != nil {
		return err
	}
	if resp.Code != "OK" {
		return errors.New(resp.Message)
	}
	return nil
}

func (sms *aliSMS) getPhones(e *v2.Event) string {
	if !strings.HasPrefix(sms.cfg.PhoneNumbers, UnfixedKeyPrefix) {
		return sms.cfg.PhoneNumbers
	}

	keyField, _ := cutPrefix(sms.cfg.PhoneNumbers, UnfixedKeyPrefix)
	eStr, _ := e.MarshalJSON()
	return gjson.Get(string(eStr), keyField).String()
}

func (sms *aliSMS) getTemplateParam(e *v2.Event) string {
	if len(sms.cfg.TemplateParam) == 0 {
		return ""
	}

	m := make(map[string]string)
	eStr, _ := e.MarshalJSON()

	for idx := range sms.cfg.TemplateParam {
		k, v := sms.cfg.TemplateParam[idx].Key, sms.cfg.TemplateParam[idx].Value
		if !strings.HasPrefix(v, UnfixedKeyPrefix) {
			m[k] = v
		} else {
			keyField, _ := cutPrefix(v, UnfixedKeyPrefix)
			m[k] = gjson.Get(string(eStr), keyField).String()
		}
	}

	jsonStr, _ := json.Marshal(m)
	return string(jsonStr)
}

func cutPrefix(s, prefix string) (after string, found bool) {
	if !strings.HasPrefix(s, prefix) {
		return s, false
	}
	return s[len(prefix):], true
}

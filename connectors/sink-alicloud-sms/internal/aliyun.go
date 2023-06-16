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
	"errors"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
)

type aliConfig struct {
	RegionId        string `json:"region_id" yaml:"region_id" validate:"required"`
	AccessKeyId     string `json:"access_key_id" yaml:"access_key_id" validate:"required"`
	AccessKeySecret string `json:"access_key_secret" yaml:"access_key_secret" validate:"required"`
	SignName        string `json:"sign_name" yaml:"sign_name" validate:"required"`
	TemplateCode    string `json:"template_code" yaml:"template_code" validate:"required"`
	TemplateParam   string `json:"template_param" yaml:"template_param" validate:"required"`
	Phones          string `json:"phones" yaml:"phones" validate:"required"`
}

type aliSMS struct {
	cfg    *aliConfig
	client *dysmsapi.Client
}

func (sms *aliSMS) init(cfg aliConfig) (err error) {
	sms.cfg = &cfg
	sms.client, err = dysmsapi.NewClientWithAccessKey(sms.cfg.RegionId, sms.cfg.AccessKeyId, sms.cfg.AccessKeySecret)
	if err != nil {
		return err
	}
	return nil
}

func (sms *aliSMS) sendMsg() (err error) {
	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = "https"
	request.PhoneNumbers = sms.cfg.Phones
	request.SignName = sms.cfg.SignName
	request.TemplateCode = sms.cfg.TemplateCode
	request.TemplateParam = sms.cfg.TemplateParam

	resp, err := sms.client.SendSms(request)
	if err != nil {
		return err
	}
	if resp.Code != "OK" {
		return errors.New(resp.Message)
	}
	return nil
}

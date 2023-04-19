// Copyright 2023 Linkall Inc.
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
	"time"

	"golang.org/x/oauth2"
)

type OAuth struct {
	ClientID     string    `json:"client_id" yaml:"client_id" validate:"required"`
	ClientSecret string    `json:"client_secret" yaml:"client_secret" validate:"required"`
	RefreshToken string    `json:"refresh_token" yaml:"refresh_token" validate:"required"`
	AccessToken  string    `json:"access_token" yaml:"access_token"`
	TokenType    string    `json:"token_type" yaml:"token_type"`
	Expiry       time.Time `json:"expiry,omitempty" yaml:"expiry"`
}

func (a *OAuth) GetToken() *oauth2.Token {
	return &oauth2.Token{
		AccessToken:  a.AccessToken,
		RefreshToken: a.RefreshToken,
		TokenType:    a.TokenType,
		Expiry:       a.Expiry,
	}

}

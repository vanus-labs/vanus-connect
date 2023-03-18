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
	"context"
	"encoding/json"
	"time"

	"golang.org/x/oauth2"

	cdkgo "github.com/vanus-labs/cdk-go"
	"github.com/vanus-labs/cdk-go/log"
	"github.com/vanus-labs/cdk-go/store"
)

type OAuth struct {
	ClientID     string    `json:"client_id" yaml:"client_id" validate:"required"`
	ClientSecret string    `json:"client_secret" yaml:"client_secret" validate:"required"`
	RefreshToken string    `json:"refresh_token" yaml:"refresh_token" validate:"required"`
	AccessToken  string    `json:"access_token" yaml:"access_token"`
	TokenType    string    `json:"token_type" yaml:"token_type"`
	Expiry       time.Time `json:"expiry,omitempty" yaml:"expiry"`
}

const storeKey = "google_sheet_oauth_token"

func (a *OAuth) TokenChange(token *oauth2.Token) {
	log.Info("receive a new oauth token", map[string]interface{}{
		"refresh_token": token.RefreshToken,
	})
	kvStore := cdkgo.GetKVStore()
	if kvStore == nil {
		log.Info("receive a new oauth token, but no store config", map[string]interface{}{
			"refresh_token": token.RefreshToken,
		})
		return
	}
	v, err := json.Marshal(token)
	if err != nil {
		log.Error("kv store save oauth token marshal error", map[string]interface{}{
			log.KeyError:    err,
			"refresh_token": token.RefreshToken,
		})
		return
	}
	err = kvStore.Set(context.Background(), storeKey, v)
	if err != nil {
		log.Error("kv store save oauth token error", map[string]interface{}{
			log.KeyError: err,
			storeKey:     string(v),
		})
		return
	}
	log.Info("save oauth token to store success", map[string]interface{}{
		"refresh_token": token.RefreshToken,
	})
}

func (a *OAuth) GetToken() *oauth2.Token {
	ctx := context.Background()
	kvStore := cdkgo.GetKVStore()
	if kvStore != nil {
		tokenValue, err := kvStore.Get(ctx, storeKey)
		if err == nil && len(tokenValue) > 0 {
			var token oauth2.Token
			err2 := json.Unmarshal(tokenValue, &token)
			if err2 != nil {
				log.Error("get token from store unmarshal error", map[string]interface{}{
					log.KeyError: err2,
					"token":      string(tokenValue),
				})
			}
			return &token
		}
		if err != nil && err != store.ErrKeyNotExist {
			log.Warning("get refresh token error", map[string]interface{}{
				log.KeyError: err,
			})
		}
	} else {
		log.Warning("no store config, it will lost token", nil)
	}
	return &oauth2.Token{
		AccessToken:  a.AccessToken,
		RefreshToken: a.RefreshToken,
		TokenType:    a.TokenType,
		Expiry:       a.Expiry,
	}

}

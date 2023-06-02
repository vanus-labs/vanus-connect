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

package client

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
	"golang.org/x/oauth2"

	"github.com/vanus-labs/cdk-go/log"
)

const (
	tokenURL = "https://aip.baidubce.com/oauth/2.0/token"
)

type tokenSource struct {
	ClientID     string
	ClientSecret string
	httpClient   *resty.Client
}

func NewTokenSource(clientID, clientSecret string) oauth2.TokenSource {
	t := &tokenSource{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		httpClient:   resty.New(),
	}
	return &reuseTokenSource{
		new: t,
	}
}

type tokenJSON struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int32  `json:"expires_in"` // at least PayPal returns string, while most return number
}

func (e *tokenJSON) expiry() (t time.Time) {
	if v := e.ExpiresIn; v != 0 {
		return time.Now().Add(time.Duration(v) * time.Second)
	}
	return
}

func (t *tokenSource) Token() (*oauth2.Token, error) {
	res, err := t.httpClient.R().SetQueryParam("grant_type", "client_credentials").
		SetQueryParam("client_id", t.ClientID).
		SetQueryParam("client_secret", t.ClientSecret).
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").Post(tokenURL)
	if err != nil {
		return nil, err
	}
	var token tokenJSON
	err = json.Unmarshal(res.Body(), &token)
	if err != nil {
		log.Info().Str("body", res.String()).Msg("get token unmarshal failed")
		return nil, err
	}
	return &oauth2.Token{
		AccessToken: token.AccessToken,
		Expiry:      token.expiry(),
	}, nil
}

type reuseTokenSource struct {
	new oauth2.TokenSource // called when t is expired.

	mu sync.Mutex
	t  *oauth2.Token
}

func (s *reuseTokenSource) Token() (*oauth2.Token, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.t.Valid() {
		return s.t, nil
	}
	t, err := s.new.Token()
	if err != nil {
		return nil, err
	}
	s.t = t
	return t, nil
}

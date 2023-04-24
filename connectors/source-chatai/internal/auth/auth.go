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

package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"hash"
	"net/http"
	"strings"
)

var (
	ErrBasic            = errors.New("header authorization invalid")
	ErrSignatureMissing = errors.New("header signature missing")
	ErrSignatureInvalid = errors.New("header signature invalid")
)

type Auth interface {
	Auth(req http.Header, body []byte) bool
	Write(w http.ResponseWriter)
}

type basicAuth struct {
	config   BasicAuth
	username string
	password string
}

func (a *basicAuth) Write(w http.ResponseWriter) {
	w.Header().Set("WWW-Authenticate", `Basic realm= "vans connect", charset="UTF-8"`)
	http.Error(w, "Unauthorized", http.StatusForbidden)
}

func (a *basicAuth) Auth(_ http.Header, _ []byte) bool {
	return a.config.Username == a.username && a.config.Password == a.password
}

type hmacAuth struct {
	signature []byte
	hmac      hash.Hash
}

func (a *hmacAuth) Write(w http.ResponseWriter) {
	http.Error(w, "Unauthorized", http.StatusForbidden)
}

func (a *hmacAuth) Auth(_ http.Header, body []byte) bool {
	_, _ = a.hmac.Write(body)
	computed := a.hmac.Sum(nil)
	return hmac.Equal(computed, a.signature)
}

func NewAuth(config *Config, req *http.Request) (Auth, error) {
	if config == nil {
		return nil, nil
	}
	switch config.Type {
	case Basic:
		username, password, ok := req.BasicAuth()
		if !ok {
			return nil, ErrBasic
		}
		return &basicAuth{
			config:   config.Basic,
			username: username,
			password: password,
		}, nil
	case Hmac:
		signature := req.Header.Get(config.HMAC.Header)
		if signature == "" {
			return nil, ErrSignatureMissing
		}
		bsignature, err := hex.DecodeString(strings.TrimPrefix(signature, "sha256="))
		if err != nil {
			return nil, ErrSignatureInvalid
		}
		header := config.HMAC.Header
		if header == "" {
			header = DefaultHeaderSignature
		}
		hash := hmac.New(sha256.New, []byte(config.HMAC.Secret))
		return &hmacAuth{
			signature: bsignature,
			hmac:      hash,
		}, nil
	}
	return nil, errors.New("unknown auth type")
}

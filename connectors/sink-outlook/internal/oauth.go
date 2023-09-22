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

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"golang.org/x/oauth2"
)

type OAuth struct {
	TenantID     string `json:"tenant_id" yaml:"tenant_id"`
	ClientID     string `json:"client_id" yaml:"client_id" validate:"required"`
	ClientSecret string `json:"client_secret" yaml:"client_secret" validate:"required"`
	RefreshToken string `json:"refresh_token" yaml:"refresh_token"`
	TokenType    string `json:"token_type" yaml:"token_type"`
}

func (a *OAuth) GetToken() *oauth2.Token {
	return &oauth2.Token{
		RefreshToken: a.RefreshToken,
		TokenType:    a.TokenType,
	}
}

type AzureToken struct {
	Token oauth2.TokenSource
}

func (t *AzureToken) GetToken(ctx context.Context, opts policy.TokenRequestOptions) (azcore.AccessToken, error) {
	token, err := t.Token.Token()
	if err != nil {
		return azcore.AccessToken{}, err
	}
	return azcore.AccessToken{
		Token:     token.AccessToken,
		ExpiresOn: token.Expiry,
	}, nil
}

type LocalAzureToken struct {
	Token *oauth2.Token
}

func (t *LocalAzureToken) GetToken(ctx context.Context, opts policy.TokenRequestOptions) (azcore.AccessToken, error) {
	return azcore.AccessToken{
		Token:     t.Token.AccessToken,
		ExpiresOn: t.Token.Expiry,
	}, nil
}

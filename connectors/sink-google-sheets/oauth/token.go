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

package oauth

import (
	"context"
	"net/url"

	"github.com/vanus-labs/connector/sink/googlesheets/oauth/internal"
	"golang.org/x/oauth2"
)

// copy from golang.org/x/oauth2

// tokenFromInternal maps an *internal.Token struct into
// a *Token struct.
func tokenFromInternal(t *internal.Token) *oauth2.Token {
	if t == nil {
		return nil
	}
	return &oauth2.Token{
		AccessToken:  t.AccessToken,
		TokenType:    t.TokenType,
		RefreshToken: t.RefreshToken,
		Expiry:       t.Expiry,
	}
}

// retrieveToken takes a *Config and uses that to retrieve an *internal.Token.
// This token is then mapped from *internal.Token into an *oauth2.Token which is returned along
// with an error..
func retrieveToken(ctx context.Context, c *Config, v url.Values) (*oauth2.Token, error) {
	tk, err := internal.RetrieveToken(ctx, c.ClientID, c.ClientSecret, c.Endpoint.TokenURL, v, internal.AuthStyle(c.Endpoint.AuthStyle))
	if err != nil {
		if rErr, ok := err.(*internal.RetrieveError); ok {
			return nil, (*oauth2.RetrieveError)(rErr)
		}
		return nil, err
	}
	return tokenFromInternal(tk), nil
}

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
	"sync"

	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

// copy from golang.org/x/oauth2, only add for tokenChange

type Config struct {
	oauth2.Config
	TokenChange func(token *oauth2.Token)
}

// TokenSource returns a TokenSource that returns t until t expires,
// automatically refreshing it as necessary using the provided context.
//
// Most users will use Config.Client instead.
func (c *Config) TokenSource(ctx context.Context, t *oauth2.Token) oauth2.TokenSource {
	tkr := &tokenRefresher{
		ctx:  ctx,
		conf: c,
	}
	if t != nil {
		tkr.refreshToken = t.RefreshToken
	}
	return &reuseTokenSource{
		t:   t,
		new: tkr,
	}
}

// tokenRefresher is a TokenSource that makes "grant_type"=="refresh_token"
// HTTP requests to renew a token using a RefreshToken.
type tokenRefresher struct {
	ctx          context.Context // used to get HTTP requests
	conf         *Config
	refreshToken string
}

// WARNING: Token is not safe for concurrent access, as it
// updates the tokenRefresher's refreshToken field.
// Within this package, it is used by reuseTokenSource which
// synchronizes calls to this method with its own mutex.
func (tf *tokenRefresher) Token() (*oauth2.Token, error) {
	if tf.refreshToken == "" {
		return nil, errors.New("oauth2: token expired and refresh token is not set")
	}

	tk, err := retrieveToken(tf.ctx, tf.conf, url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {tf.refreshToken},
	})

	if err != nil {
		return nil, err
	}
	if tf.refreshToken != tk.RefreshToken {
		tf.refreshToken = tk.RefreshToken

		// token change call
		tf.conf.TokenChange(tk)
	}
	return tk, err
}

// reuseTokenSource is a TokenSource that holds a single token in memory
// and validates its expiry before each call to retrieve it with
// Token. If it's expired, it will be auto-refreshed using the
// new TokenSource.
type reuseTokenSource struct {
	new oauth2.TokenSource // called when t is expired.

	mu sync.Mutex // guards t
	t  *oauth2.Token
}

// Token returns the current token if it's still valid, else will
// refresh the current token (using r.Context for HTTP client
// information) and return the new one.
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

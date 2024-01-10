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
	"errors"
	"fmt"
	"time"

	"github.com/vanus-labs/cdk-go/store"
)

type apiType string

const (
	OrderApi   apiType = "orders"
	ProductApi apiType = "products"
)

func syncBeginDateKey() string {
	return fmt.Sprintf("sync_begin_date")
}

func (s *shopifySource) getSyncBeginDate(ctx context.Context) (string, error) {
	v, err := s.store.Get(ctx, syncBeginDateKey())
	if err != nil {
		if errors.Is(err, store.ErrKeyNotExist) {
			return "", nil
		}
		return "", err
	}
	return string(v), nil
}

func (s *shopifySource) setSyncBeginDate(ctx context.Context, t string) error {
	return s.store.Set(ctx, syncBeginDateKey(), []byte(t))
}

func syncTimeKey(apiType apiType) string {
	return fmt.Sprintf("sync_time_%s", apiType)
}

func (s *shopifySource) getSyncTime(ctx context.Context, apiType apiType) (time.Time, error) {
	v, err := s.store.Get(ctx, syncTimeKey(apiType))
	if err != nil {
		return time.Time{}, err
	}
	t, err := time.Parse(time.RFC3339, string(v))
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

func (s *shopifySource) setSyncTime(ctx context.Context, apiType apiType, t time.Time) error {
	return s.store.Set(ctx, syncTimeKey(apiType), []byte(t.Format(time.RFC3339)))
}

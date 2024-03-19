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

type ApiType string

const (
	OrderApi   ApiType = "orders"
	ProductApi ApiType = "products"
)

func (s *shopifySource) initSyncTime(ctx context.Context) error {
	syncBeginDate, err := s.getSyncBeginDate(ctx)
	if err != nil {
		return errors.Join(err, errors.New("get sync begin date error"))
	}
	if syncBeginDate == s.config.SyncBeginDate {
		s.logger.Info().
			Str("sync_begin_date", s.config.SyncBeginDate).
			Msg("sync begin date no change")
		return nil
	}
	for _, t := range syncApiArr {
		err = s.setSyncTime(ctx, t, s.syncBeginTime)
		if err != nil {
			return errors.Join(err, errors.New(string(t)+"set sync time error"))
		}
	}
	err = s.setSyncBeginDate(ctx, s.config.SyncBeginDate)
	if err != nil {
		return errors.Join(err, errors.New("set sync begin date error"))
	}
	s.logger.Info().
		Str("sync_begin_date", s.config.SyncBeginDate).
		Msg("init sync time success")
	return nil
}

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

func syncTimeKey(apiType ApiType) string {
	return fmt.Sprintf("sync_time_%s", apiType)
}

func (s *shopifySource) getSyncTime(ctx context.Context, apiType ApiType) (time.Time, error) {
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

func (s *shopifySource) setSyncTime(ctx context.Context, apiType ApiType, t time.Time) error {
	return s.store.Set(ctx, syncTimeKey(apiType), []byte(t.Format(time.RFC3339)))
}

func syncApiKey(apiType ApiType) string {
	return fmt.Sprintf("sync_date_%s", apiType)
}

func (s *shopifySource) isApiNeedSync(ctx context.Context, apiType ApiType) (bool, error) {
	v, err := s.store.Get(ctx, syncApiKey(apiType))
	if err != nil {
		if errors.Is(err, store.ErrKeyNotExist) {
			return true, nil
		}
		return false, err
	}
	return s.config.SyncBeginDate != string(v), nil
}

func (s *shopifySource) setApiSync(ctx context.Context, apiType ApiType) error {
	return s.store.Set(ctx, syncApiKey(apiType), []byte(s.config.SyncBeginDate))
}

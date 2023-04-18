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
	"fmt"
	"time"

	cdkgo "github.com/vanus-labs/cdk-go"
	"github.com/vanus-labs/cdk-go/store"
)

const (
	storeKeySyncTime = "sync_time"
)

type apiType string

const (
	OrderApi   apiType = "orders"
	ProductApi apiType = "products"
)

func syncTimeKey(apiType apiType) string {
	return fmt.Sprintf("sync_time_%s", apiType)
}

func (s *shopifySource) getSyncBeginTime(ctx context.Context, apiType apiType) (time.Time, error) {
	kvStore := cdkgo.GetKVStore()
	v, err := kvStore.Get(ctx, syncTimeKey(apiType))
	if err != nil {
		if err == store.ErrKeyNotExist {
			return s.syncBeginTime, nil
		}
		return time.Time{}, err
	}
	t, err := time.Parse(time.RFC3339, string(v))
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

func (s *shopifySource) setSyncTime(ctx context.Context, apiType apiType, t time.Time) error {
	kvStore := cdkgo.GetKVStore()
	return kvStore.Set(ctx, syncTimeKey(apiType), []byte(t.Format(time.RFC3339)))
}

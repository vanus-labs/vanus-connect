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
	"time"

	ce "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"

	cdkgo "github.com/vanus-labs/cdk-go"
	"github.com/vanus-labs/cdk-go/log"
)

const (
	EventTypeUserInfo   = "userInfo"
	EventTypeFans       = "fans"
	EventTypeFollowings = "followings"
)

func (s *DouyinSource) syncUser(ctx context.Context) {
	s.getUser()
	tk := time.NewTicker(6 * time.Hour)
	defer tk.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-tk.C:
			s.getUser()
		}
	}
}

func (s *DouyinSource) getUser() {
	s.getUserInfo()
	s.getFans()
	s.getFollowings()
}

func (s *DouyinSource) getUserInfo() {
	s.Limiter.Take()
	userInfo, err := s.douyin.GetUser().GetUserInfo(s.config.DouyinToken.OpenID)
	if err != nil {
		log.Warning("GetUserInfo", map[string]interface{}{
			"error": err,
		})
		return
	}

	event := ce.NewEvent()
	event.SetSource(EventSource)
	event.SetType(EventTypeUserInfo)
	event.SetTime(time.Now())
	event.SetID(uuid.New().String())
	_ = event.SetData(ce.ApplicationJSON, userInfo)
	s.events <- &cdkgo.Tuple{
		Event: &event,
	}
}

func (s *DouyinSource) getFans() {
	hasMore := true
	cursor, pageSize := int64(0), int64(30)
	for hasMore {
		s.Limiter.Take()
		items, err := s.douyin.GetUser().ListFans(s.config.DouyinToken.OpenID, cursor, pageSize)
		if err != nil {
			log.Warning("ListFans", map[string]interface{}{
				"error": err,
			})
			return
		}

		for i := range items.List {
			item := items.List[i]

			event := ce.NewEvent()
			event.SetSource(EventSource)
			event.SetType(EventTypeFans)
			event.SetTime(time.Now())
			event.SetID(uuid.New().String())
			_ = event.SetData(ce.ApplicationJSON, item)
			s.events <- &cdkgo.Tuple{
				Event: &event,
			}
		}

		cursor = items.Cursor
		hasMore = items.HasMore
	}
}

func (s *DouyinSource) getFollowings() {
	hasMore := true
	cursor, pageSize := int64(0), int64(30)
	for hasMore {
		s.Limiter.Take()
		items, err := s.douyin.GetUser().ListFollowing(s.config.DouyinToken.OpenID, cursor, pageSize)
		if err != nil {
			log.Warning("ListFollowing", map[string]interface{}{
				"error": err,
			})
			return
		}

		for i := range items.List {
			item := items.List[i]

			event := ce.NewEvent()
			event.SetSource(EventSource)
			event.SetType(EventTypeFollowings)
			event.SetTime(time.Now())
			event.SetID(uuid.New().String())
			_ = event.SetData(ce.ApplicationJSON, item)
			s.events <- &cdkgo.Tuple{
				Event: &event,
			}
		}

		cursor = items.Cursor
		hasMore = items.HasMore
	}
}

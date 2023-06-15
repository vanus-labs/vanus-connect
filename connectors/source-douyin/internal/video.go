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
	EventSource = "douyin"
	EventType   = "video"
)

func (s *DouyinSource) syncVideo(ctx context.Context) {
	s.getVideo()
	tk := time.NewTicker(12 * time.Hour)
	defer tk.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-tk.C:
			s.getVideo()
		}
	}
}

func (s *DouyinSource) getVideo() {
	hasMore := true
	cursor, pageSize := int64(0), int64(30)
	for hasMore {
		s.Limiter.Take()
		info, err := s.openAPI.GetVideo().List(s.openID, cursor, pageSize)
		if err != nil {
			log.Warning("getVideo", map[string]interface{}{
				"error": err,
			})
			return
		}

		for i := range info.List {
			video := info.List[i]

			event := ce.NewEvent()
			event.SetSource(EventSource)
			event.SetType(EventType)
			event.SetTime(time.Now())
			event.SetID(uuid.New().String())
			_ = event.SetData(ce.ApplicationJSON, video)
			s.events <- &cdkgo.Tuple{
				Event: &event,
			}
		}

		cursor = info.Cursor
		hasMore = info.HasMore
	}
}

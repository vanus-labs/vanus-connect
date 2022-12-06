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
	"net/http"

	ce "github.com/cloudevents/sdk-go/v2"
	cdkgo "github.com/linkall-labs/cdk-go"
)

type DorisSink struct {
	streamLoad *StreamLoad
}

func (s *DorisSink) Initialize(_ context.Context, config cdkgo.ConfigAccessor) error {
	cfg := config.(*Config)
	// init stream load
	s.streamLoad = NewStreamLoad(cfg)
	return s.streamLoad.Start()
}

func (s *DorisSink) Name() string {
	return "DorisSink"
}

func (s *DorisSink) Destroy() error {
	if s.streamLoad != nil {
		s.streamLoad.Stop()
	}
	return nil
}

func (s *DorisSink) Arrived(ctx context.Context, events ...*ce.Event) cdkgo.Result {
	for _, event := range events {
		err := s.streamLoad.WriteEvent(ctx, event)
		if err != nil {
			return cdkgo.NewResult(http.StatusInternalServerError, err.Error())
		}
	}
	return cdkgo.SuccessResult
}

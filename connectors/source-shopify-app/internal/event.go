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
	goshopify "github.com/bold-commerce/go-shopify/v3"
	ce "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"

	cdkgo "github.com/vanus-labs/cdk-go"
)

func (s *shopifySource) newEvent() ce.Event {
	event := ce.NewEvent()
	event.SetID(uuid.NewString())
	event.SetSource("shopify-source-" + s.config.ShopName)
	return event
}

func (s *shopifySource) orderEvent(orders []goshopify.Order) {
	for _, order := range orders {
		event := s.newEvent()
		event.SetType("orders")
		event.SetData(ce.ApplicationJSON, order)
		s.events <- &cdkgo.Tuple{
			Event: &event,
		}
	}
}

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
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/vanus-labs/cdk-go/log"

	cdkgo "github.com/vanus-labs/cdk-go"

	ce "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
)

const (
	eventSource = "vanus.facebook"
)

var (
	errMissingHubSignatureHeader = errors.New("missing X-Hub-Signature-256 Header")
	errInvalidContentTypeHeader  = errors.New("only support application/json Content-Type Header")
	errReadPayload               = errors.New("error read payload")
	errVerificationFailed        = errors.New("signature verification failed")
	errInvalidPayload            = errors.New("payload invalid")
)

func (s *facebookSource) verifyRequestSignature(req *http.Request, body []byte) error {
	if s.config.AppSecret == "" {
		return nil
	}
	signature := req.Header.Get("X-Hub-Signature-256")
	if signature == "" {
		return errMissingHubSignatureHeader
	}
	mac := hmac.New(sha256.New, []byte(s.config.AppSecret))
	_, _ = mac.Write(body)
	expectedMAC := hex.EncodeToString(mac.Sum(nil))
	// sha256=signature
	if !hmac.Equal([]byte(signature[7:]), []byte(expectedMAC)) {
		return errVerificationFailed
	}
	return nil
}

func (s *facebookSource) event(req *http.Request) error {
	contentType := req.Header.Get("Content-Type")
	if contentType != "application/json" {
		return errInvalidContentTypeHeader
	}
	body, err := ioutil.ReadAll(req.Body)
	if err != nil || len(body) == 0 {
		return errReadPayload
	}
	err = s.verifyRequestSignature(req, body)
	if err != nil {
		return err
	}
	var payload []PageEvent
	err = json.Unmarshal(body, &payload)
	if err != nil {
		return err
	}
	for _, e := range payload {
		switch e.Object {
		case "page":
			s.pageEvent(e.Object, e)
		default:
			log.Info("not support object type", map[string]interface{}{
				"object": e.Object,
			})
			return nil
		}
	}
	return nil
}

func (s *facebookSource) pageEvent(object string, e PageEvent) error {
	for _, entry := range e.Entry {
		pageID := entry.ID
		if pageID == "" {
			return errors.New("page id is empty")
		}
		for _, change := range entry.Changes {
			field := change.Field
			if field == "" {
				return errors.New("change filed is empty")
			}
			event := ce.NewEvent()
			event.SetID(uuid.New().String())
			event.SetSource(eventSource)
			event.SetType(object)
			event.SetExtension("pageid", pageID)
			event.SetExtension("fields", field)
			event.SetData(ce.ApplicationJSON, change)
			if field == "feed" {
				value, ok := change.Value.(map[string]interface{})
				if !ok {
					return errors.New("value is invalid")
				}
				item, ok := value["item"].(string)
				if !ok {
					return errors.New("change value item is invalid")
				}
				event.SetExtension("changeitem", item)
			}
			s.ch <- &cdkgo.Tuple{
				Event: &event,
			}
		}
		for _, message := range entry.Messaging {
			event := ce.NewEvent()
			event.SetID(uuid.New().String())
			event.SetSource(eventSource)
			event.SetType(object)
			event.SetExtension("pageid", pageID)
			event.SetExtension("fields", "messages")
			event.SetData(ce.ApplicationJSON, message)
		}
	}
	return nil
}

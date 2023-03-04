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
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	ce "github.com/cloudevents/sdk-go/v2"
	cdkgo "github.com/linkall-labs/cdk-go"
	"github.com/linkall-labs/cdk-go/log"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

var (
	errInvalidHTTPMethod         = errors.New("invalid HTTP Method")
	errInvalidContentTypeHeader  = errors.New("only support application/json Content-Type Header")
	errMissingGithubEventHeader  = errors.New("missing X-GitHub-Event Header")
	errMissingHubSignatureHeader = errors.New("missing X-Hub-Signature-256 Header")
	errMissingHubDeliveryHeader  = errors.New("missing X-GitHub-Delivery Header")
	errReadPayload               = errors.New("error read payload")
	errVerificationFailed        = errors.New("signature verification failed")
	errPingEvent                 = errors.New("receive ping event")
	errInvalidBody               = errors.New("invalid body")
)

type handler struct {
	config GitHubCfg
	events chan *cdkgo.Tuple
	client *http.Client
}

func newHandler(events chan *cdkgo.Tuple, config GitHubCfg) *handler {
	h := &handler{
		config: config,
		events: events,
	}
	if config.AccessToken != "" {
		h.client = oauth2.NewClient(context.Background(),
			oauth2.StaticTokenSource(
				&oauth2.Token{
					AccessToken: config.AccessToken,
				}))
		h.client.Timeout = 5 * time.Second
	}
	return h
}

func (h *handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	err := h.handle(req)
	if err == nil {
		w.WriteHeader(http.StatusAccepted)
		_, _ = w.Write([]byte("accepted"))
		return
	}
	var code int
	switch err {
	case errPingEvent:
		code = http.StatusOK
	case errInvalidHTTPMethod, errInvalidContentTypeHeader, errMissingGithubEventHeader, errMissingHubDeliveryHeader, errMissingHubSignatureHeader, errReadPayload, errInvalidBody:
		code = http.StatusBadRequest
	case errVerificationFailed:
		code = http.StatusForbidden
	default:
		code = http.StatusInternalServerError
	}
	w.WriteHeader(code)
	_, _ = w.Write([]byte(err.Error()))
}

func (h *handler) handle(req *http.Request) error {
	if req.Method != http.MethodPost {
		return errInvalidHTTPMethod
	}
	contentType := req.Header.Get(HeaderContentType)
	if contentType != "application/json" {
		return errMissingGithubEventHeader
	}
	eventType := req.Header.Get(GHHeaderEvent)
	if eventType == "" {
		return errMissingGithubEventHeader
	}
	log.Info("receive event", map[string]interface{}{
		"eventType": eventType,
	})
	if eventType == "ping" {
		return errPingEvent
	}
	eventID := req.Header.Get(GHHeaderDelivery)
	if eventID == "" {
		return errMissingHubDeliveryHeader
	}
	body, err := ioutil.ReadAll(req.Body)
	if err != nil || len(body) == 0 {
		return errReadPayload
	}
	if h.config.WebHookSecret != "" {
		signature := req.Header.Get(GHHeaderSignature256)
		if signature == "" {
			return errMissingHubSignatureHeader
		}
		mac := hmac.New(sha256.New, []byte(h.config.WebHookSecret))
		_, _ = mac.Write(body)
		expectedMAC := hex.EncodeToString(mac.Sum(nil))
		// sha256=signature
		if !hmac.Equal([]byte(signature[7:]), []byte(expectedMAC)) {
			return errVerificationFailed
		}
	}

	event := ce.NewEvent()
	event.SetID(eventID)
	err = h.setEvent(&event, eventType, body)
	if err != nil {
		return err
	}
	h.events <- &cdkgo.Tuple{Event: &event}
	return nil
}

func (h *handler) setEvent(event *ce.Event, eventType string, body []byte) error {
	var payload map[string]interface{}
	err := json.Unmarshal(body, &payload)
	if err != nil {
		return err
	}
	repo, ok := payload["repository"].(map[string]interface{})
	if !ok {
		return errInvalidBody
	}
	event.SetSource(getString(repo["url"]))
	t := "com.github." + eventType
	action, ok := payload["action"].(string)
	if ok {
		event.SetType(t + "." + action)
	}
	switch eventType {
	case "star":
		event.SetTime(getTime(repo["starred_at"]))
	case "push":
		event.SetType(t)
		event.SetSubject(getString(payload["ref"]))
		event.SetTime(getTime(repo["updated_at"]))
	case "issues":
		issue, ok := payload["issue"].(map[string]interface{})
		if ok {
			event.SetSubject(getString(issue["number"]))
			event.SetTime(getTime(issue["updated_at"]))
		}
	case "check_run":
		checkRun, ok := payload["check_run"].(map[string]interface{})
		if ok {
			event.SetSubject(getString(checkRun["id"]))
			time, ok := checkRun["completed_at"].(string)
			if !ok {
				time, _ = checkRun["started_at"].(string)
			}
			event.SetTime(getTime(time))
		}
	case "check_suite":
		checkSuit, ok := payload["check_suite"].(map[string]interface{})
		if ok {
			event.SetSubject(getString(checkSuit["id"]))
			event.SetTime(getTime(checkSuit["updated_at"]))
		}
	case "commit_comment":
		comment, ok := payload["comment"].(map[string]interface{})
		if ok {
			event.SetSource(getString(comment["url"]) + "/" + getString(comment["comment_id"]))
			event.SetSubject(getString(comment["id"]))
			event.SetTime(getTime(comment["updated_at"]))
		}
	case "content_reference":
		reference, ok := payload["content_reference"].(map[string]interface{})
		if ok {
			event.SetSubject(getString(reference["id"]))
		}
		event.SetTime(time.Now())
	case "create", "delete":
		event.SetType(t + "." + getString(payload["ref_type"]))
		event.SetSubject(getString(payload["ref"]))
		event.SetTime(time.Now())
	case "deploy_key":
		key, ok := payload["key"].(map[string]interface{})
		if ok {
			event.SetSubject(getString(key["id"]))
			time, ok := key["deleted_at"].(string)
			if !ok {
				time, _ = key["created_at"].(string)
			}
			event.SetTime(getTime(time))
		}
	case "deployment":
		event.SetType(t)
		deployment, ok := payload["deployment"].(map[string]interface{})
		if ok {
			event.SetSubject(getString(deployment["id"]))
			event.SetTime(getTime(deployment["updated_at"]))
		}
	case "deployment_status":
		deployment, ok := payload["deployment"].(map[string]interface{})
		if ok {
			event.SetSource(getString(deployment["url"]))
		}
		deploymentStatus, ok := payload["deployment_status"].(map[string]interface{})
		if ok {
			event.SetType(t + "." + getString(deploymentStatus["state"]))
			event.SetSubject(getString(deploymentStatus["url"]))
			event.SetTime(getTime(deploymentStatus["updated_at"]))
		}
	case "fork":
		event.SetType(t)
		forkee, ok := payload["forkee"].(map[string]interface{})
		if ok {
			event.SetSubject(getString(forkee["url"]))
			event.SetTime(getTime(forkee["created_at"]))
		}
	case "github_app_authorization":
		event.SetType(t)
		event.SetTime(time.Now())
		sender, ok := payload["sender"].(map[string]interface{})
		if ok {
			event.SetSource(getString(sender["url"]))
		}
	case "gollum":
		event.SetTime(time.Now())
		pages, ok := payload["pages"].(map[string]interface{})
		if ok {
			event.SetType(t + "." + getString(pages["action"]))
			event.SetSubject(getString(pages["page_name"]))
		}
	case "installation", "installation_repositories":
		installation, ok := payload["installation"].(map[string]interface{})
		if ok {
			event.SetSubject(getString(installation["id"]))
			event.SetTime(getTimeByTimestamp(installation["updated_at"]))
			account, ok := installation["account"].(map[string]interface{})
			if ok {
				event.SetSource(getString(account["url"]))
			}
		}
	case "issue_comment":
		issue, ok := payload["issue"].(map[string]interface{})
		if ok {
			event.SetSource(getString(issue["url"]))
		}
		comment, ok := payload["comment"].(map[string]interface{})
		if ok {
			event.SetSubject(getString(comment["id"]))
			event.SetTime(getTime(comment["updated_at"]))
		}
	case "label":
		label, ok := payload["label"].(map[string]interface{})
		if ok {
			event.SetSubject(getString(label["name"]))
		}
		event.SetTime(time.Now())
	case "marketplace_purchase":
		sender, ok := payload["sender"].(map[string]interface{})
		if ok {
			event.SetSource(strings.ReplaceAll(getString(sender["url"]), "/username", ""))
		}
		label, ok := payload["label"].(map[string]interface{})
		if ok {
			event.SetSubject(getString(label["name"]))
		}
		event.SetTime(getTime(payload["effective_date"]))
	case "member":
		member, ok := payload["member"].(map[string]interface{})
		if ok {
			event.SetSubject(getString(member["login"]))
		}
		event.SetTime(time.Now())
	case "membership":
		member, ok := payload["member"].(map[string]interface{})
		if ok {
			event.SetSubject(getString(member["login"]))
		}
		event.SetTime(time.Now())
		event.SetType(t + "." + getString(payload["scope"]) + action)
	case "meta":
		hook, ok := payload["hook"].(map[string]interface{})
		if ok {
			event.SetSubject(getString(hook["hook_id"]))
			event.SetTime(getTime(hook["updated_at"]))
		}
	case "milestone":
		milestone, ok := payload["milestone"].(map[string]interface{})
		if ok {
			event.SetSubject(getString(milestone["number"]))
			event.SetTime(getTime(milestone["updated_at"]))
		}
	case "organization":
		organization, ok := payload["organization"].(map[string]interface{})
		if ok {
			event.SetSource(getString(organization["url"]))
		}
		membership, ok := payload["membership"].(map[string]interface{})
		if ok {
			user, ok := membership["user"].(map[string]interface{})
			if ok {
				event.SetSubject(getString(user["login"]))
			}
		}
		event.SetTime(time.Now())
	case "org_block":
		organization, ok := payload["organization"].(map[string]interface{})
		if ok {
			event.SetSource(getString(organization["url"]))
		}
		blockedUser, ok := payload["blocked_user"].(map[string]interface{})
		if ok {
			event.SetSubject(getString(blockedUser["login"]))
		}
		event.SetTime(time.Now())
	case "page_build":
		event.SetType(t)
		build, ok := payload["build"].(map[string]interface{})
		if ok {
			event.SetSubject(getString(build["url"]))
		}
		pusher, ok := payload["pusher"].(map[string]interface{})
		if ok {
			event.SetTime(getTime(pusher["updated_at"]))
		}
	case "project_card", "project_column", "project":
		project, ok := payload[eventType].(map[string]interface{})
		if ok {
			event.SetSubject(getString(project["id"]))
			event.SetTime(getTime(project["updated_at"]))
		}
	case "repository":
		owner, ok := repo["owner"].(map[string]interface{})
		if ok {
			event.SetSource(getString(owner["url"]))
		}
		event.SetSubject(getString(repo["name"]))
		event.SetTime(getTime(repo["updated_at"]))
	case "public", "repository_import":
		event.SetType(t)
		owner, ok := repo["owner"].(map[string]interface{})
		if ok {
			event.SetSource(getString(owner["url"]))
		}
		event.SetSubject(getString(repo["name"]))
		event.SetTime(getTime(repo["updated_at"]))
	case "pull_request":
		event.SetSubject(getString(payload["number"]))
		event.SetTime(getTime(repo["updated_at"]))
	case "pull_request_review":
		pull, ok := payload["pull_request"].(map[string]interface{})
		if ok {
			event.SetSource(getString(pull["url"]))
		}
		review, ok := repo["review"].(map[string]interface{})
		if ok {
			event.SetSubject(getString(review["id"]))
			event.SetTime(getTime(review["submitted_at"]))
		}
	case "pull_request_review_comment":
		pull, ok := payload["pull_request"].(map[string]interface{})
		if ok {
			event.SetSource(getString(pull["url"]))
			event.SetTime(getTime(pull["updated_at"]))
		}
		comment, ok := payload["comment"].(map[string]interface{})
		if ok {
			event.SetSubject(getString(comment["id"]))
		}
	case "registry_package":
		registry, ok := payload["registry_package"].(map[string]interface{})
		if ok {
			event.SetSubject(getString(registry["html_url"]))
			event.SetTime(getTime(registry["updated_at"]))
		}
	case "release":
		release, ok := payload["release"].(map[string]interface{})
		if ok {
			event.SetSubject(getString(release["id"]))
			time, ok := release["published_at"].(string)
			if !ok {
				time, _ = release["created_at"].(string)
			}
			event.SetTime(getTime(time))
		}
	case "repository_vulnerability_alert":
		alert, ok := payload["alert"].(map[string]interface{})
		if ok {
			event.SetSubject(getString(alert["id"]))
		}
		event.SetTime(time.Now())
	case "security_advisory":
		advisory, ok := payload["security_advisory"].(map[string]interface{})
		if ok {
			event.SetSubject(getString(advisory["ghsa_id"]))
			event.SetTime(getTime(advisory["updated_at"]))
		}
	case "status":
		event.SetSubject(getString(payload["sha"]))
		event.SetTime(getTime(payload["updated_at"]))
	case "team", "team_add":
		event.SetTime(getTime(payload["updated_at"]))
		team, ok := payload["team"].(map[string]interface{})
		if ok {
			event.SetSubject(getString(team["id"]))
		}
	case "watch":
		event.SetTime(time.Now())
	default:
		log.Info("unknown event type", map[string]interface{}{
			"eventType": eventType,
		})
		event.SetTime(time.Now())
	}
	if h.client != nil {
		sender, ok := payload["sender"].(map[string]interface{})
		if ok {
			url, ok := sender["url"].(string)
			if ok && url != "" {
				resp, err := h.client.Get(url)
				if err != nil {
					log.Error("get url error", map[string]interface{}{
						log.KeyError: err,
						"url":        url,
					})
					return err
				}
				defer resp.Body.Close()
				var res map[string]interface{}
				if err = json.NewDecoder(resp.Body).Decode(&res); err != nil {
					log.Error("parse response error", map[string]interface{}{
						log.KeyError: err,
					})
					return err
				}
				sender["url_response"] = res
				_ = event.SetData(ce.ApplicationJSON, payload)
				return nil
			}
		}
	}
	_ = event.SetData(ce.ApplicationJSON, body)
	return nil
}

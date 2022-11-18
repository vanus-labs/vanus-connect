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
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	ce "github.com/cloudevents/sdk-go/v2"
	"github.com/linkall-labs/cdk-go/log"
	"github.com/pkg/errors"
	"io"
	"net/http"
	pkgurl "net/url"
	"sync"
	"time"
)

type StreamLoad struct {
	logger       log.Logger
	config       *Config
	loadUrl      *pkgurl.URL
	authEncoding string

	client  *http.Client
	timeout time.Duration

	eventCh      chan ce.Event
	lock         sync.Mutex
	buffer       *bytes.Buffer
	lastLoadTime time.Time

	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

const (
	maxSize      = 10 * 2 << 20
	loadSize     = maxSize - 4*2<<10
	loadInterval = 5 * time.Second
)

func NewStreamLoad(config *Config, timeout time.Duration, logger log.Logger) *StreamLoad {
	l := &StreamLoad{
		config:  config,
		timeout: timeout,
		logger:  logger,
		eventCh: make(chan ce.Event, 100),
		buffer:  bytes.NewBuffer(make([]byte, 0, maxSize)),
	}
	l.ctx, l.cancel = context.WithCancel(context.Background())
	return l
}

func (l *StreamLoad) WriteEvent(ctx context.Context, event ce.Event) error {
	select {
	case l.eventCh <- event:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (l *StreamLoad) Start() error {
	loadUrlStr := fmt.Sprintf("http://%s/api/%s/%s/_stream_load", l.config.Fenodes, l.config.DbName, l.config.TableName)
	u, err := pkgurl.Parse(loadUrlStr)
	if err != nil {
		return err
	}
	l.loadUrl = u
	l.authEncoding = base64.StdEncoding.EncodeToString([]byte(l.config.Username + ":" + l.config.Password))
	l.lastLoadTime = time.Now()
	l.client = &http.Client{
		Timeout: l.timeout,
	}
	l.wg.Add(1)
	go func() {
		defer l.wg.Done()
		t := time.NewTicker(loadInterval)
		defer t.Stop()
		for {
			select {
			case <-t.C:
				l.checkAndLoad()
			case <-l.ctx.Done():
				return
			}
		}
	}()
	l.wg.Add(1)
	go func() {
		defer l.wg.Done()
		for {
			select {
			case event, ok := <-l.eventCh:
				if !ok {
					return
				}
				l.event2Buffer(event)
			}
		}
	}()
	return nil
}

func (l *StreamLoad) Stop() {
	l.cancel()
	close(l.eventCh)
	l.wg.Wait()
	l.checkAndLoad()
}
func (l *StreamLoad) event2Buffer(event ce.Event) {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.buffer.Write(event.Data())
	l.buffer.WriteString("\n")
	if l.buffer.Len() >= loadSize {
		l.loadAndReset()
	}
}

func (l *StreamLoad) checkAndLoad() {
	l.lock.Lock()
	defer l.lock.Unlock()
	size := l.buffer.Len()
	if size == 0 {
		return
	}
	l.loadAndReset()
}

func (l *StreamLoad) loadAndReset() {
	err := l.load()
	if err != nil {
		l.logger.Warning(l.ctx, "load has error,will retry", map[string]interface{}{
			log.KeyError: err,
		})
		return
	}
	l.buffer.Reset()
	l.lastLoadTime = time.Now()
}

func (l *StreamLoad) load() error {
	label := fmt.Sprintf("vance_sink_%s_%d", l.config.TableName, time.Now().UnixMilli())
	ctx, cancel := context.WithTimeout(l.ctx, l.timeout)
	defer cancel()
	req := l.makeRequest(label)
	req = req.WithContext(ctx)
	resp, err := l.client.Do(req)
	if err != nil {
		return errors.Wrapf(err, "%s load client do error", label)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%s load response not ok", label)
	}
	var res map[string]interface{}
	if err = json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return errors.Wrapf(err, "%s resp body json decode error", label)
	}
	status := res["Status"]
	if status != "Success" && status != "Publish Timeout" {
		return fmt.Errorf("%s resp status is %s not success", label, status)
	}
	l.logger.Info(l.ctx, fmt.Sprintf("load success %s", label), nil)
	return nil
}

func (l *StreamLoad) makeRequest(label string) *http.Request {
	req := &http.Request{
		Method:        http.MethodPut,
		URL:           l.loadUrl,
		Header:        make(http.Header),
		ContentLength: int64(l.buffer.Len()),
		GetBody: func() (io.ReadCloser, error) {
			r := bytes.NewReader(l.buffer.Bytes())
			return io.NopCloser(r), nil
		},
	}
	for k, v := range l.config.StreamLoad {
		req.Header.Set(k, v)
	}
	req.Header.Set("Expect", "100-continue")
	req.Header.Set("Authorization", "Basic "+l.authEncoding)
	req.Header.Set("format", "json")
	req.Header.Set("read_json_by_line", "true")
	req.Header.Set("label", label)
	return req
}

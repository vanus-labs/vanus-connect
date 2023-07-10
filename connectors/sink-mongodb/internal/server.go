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
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	ce "github.com/cloudevents/sdk-go/v2"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	cdkgo "github.com/vanus-labs/cdk-go"
	"github.com/vanus-labs/cdk-go/config"
	"github.com/vanus-labs/cdk-go/connector"
	"github.com/vanus-labs/cdk-go/log"
)

const (
	mongoCollection = "xvdcoll"

	name = "Sink MongoDB"

	debeziumConnector = "iodebeziumconnector"
)

var _ cdkgo.SinkConfigAccessor = &Config{}

type Config struct {
	cdkgo.SinkConfig `json:",inline" yaml:",inline"`
	ConnectionURI    string     `json:"connection_uri" yaml:"connection_uri" validate:"required"`
	Database         string     `json:"database" yaml:"database" validate:"required"`
	Collection       string     `json:"collection" yaml:"collection" validate:"required"`
	Credential       Credential `json:"credential" yaml:"credential"`
	BulkSize         int        `json:"bulk_size" yaml:"bulk_size"`
	FlushInterval    int        `json:"flush_interval" yaml:"flush_interval"`
}

type Credential struct {
	Username                string            `json:"username" yaml:"username"`
	Password                string            `json:"password" yaml:"password"`
	AuthSource              string            `json:"auth_source" yaml:"auth_source"`
	AuthMechanism           string            `json:"auth_mechanism" yaml:"auth_mechanism"`
	AuthMechanismProperties map[string]string `json:"auth_mechanism_properties" yaml:"auth_mechanism_properties"`
}

func (c *Config) GetSecret() cdkgo.SecretAccessor {
	return &c.Credential
}

func (c *Config) Validate() error {
	if c.BulkSize == 0 {
		c.BulkSize = 100
	}
	if c.FlushInterval == 0 {
		c.FlushInterval = 2000
	}
	return c.SinkConfig.Validate()
}

func (c *Credential) IsSet() bool {
	return c.Username != "" || c.Password != "" || c.AuthSource != "" ||
		c.AuthMechanism != "" || len(c.AuthMechanismProperties) > 0
}

func (c Credential) GetMongoDBCredential() *options.Credential {
	return &options.Credential{
		AuthMechanism:           c.AuthMechanism,
		AuthMechanismProperties: c.AuthMechanismProperties,
		AuthSource:              c.AuthSource,
		Username:                c.Username,
		Password:                c.Password,
		PasswordSet:             true,
	}
}

func NewConfig() cdkgo.SinkConfigAccessor {
	return &Config{}
}

var _ cdkgo.Sink = &mongoSink{}

type mongoSink struct {
	cfg      *Config
	writer   map[string]*InsertWriter
	dbClient *mongo.Client
	logger   zerolog.Logger
	lock     sync.Mutex
	stop     chan bool
}

func NewMongoSink() cdkgo.Sink {
	return &mongoSink{
		writer: map[string]*InsertWriter{},
		stop:   make(chan bool),
	}
}

func (s *mongoSink) Initialize(ctx context.Context, cfg config.ConfigAccessor) error {
	c, _ := cfg.(*Config)
	s.logger = log.FromContext(ctx)
	s.cfg = c
	clientOptions := options.Client().ApplyURI(s.cfg.ConnectionURI)
	if s.cfg.Credential.IsSet() {
		clientOptions.Auth = s.cfg.Credential.GetMongoDBCredential()
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	mongoClient, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return fmt.Errorf("connect to mongodb error: %s", err.Error())
	}
	if err = mongoClient.Ping(ctx, readpref.Primary()); err != nil {
		return fmt.Errorf("failed to connect mongodb: %s", err.Error())
	} else {
		s.logger.Info().Msg("mongodb is connected")
	}
	s.dbClient = mongoClient
	go s.start()
	return nil
}

func (s *mongoSink) start() {
	ticker := time.NewTicker(time.Duration(s.cfg.FlushInterval) * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-s.stop:
			s.logger.Info().Msg("flush task stop")
			return
		case <-ticker.C:
			s.flush()
		}
	}
}

func (s *mongoSink) flush() {
	s.lock.Lock()
	defer s.lock.Unlock()
	for collName, writer := range s.writer {
		err := writer.Flush()
		if err != nil {
			s.logger.Warn().Str("collection", collName).Msg("flush error")
		}
	}
}

func (s *mongoSink) Destroy() error {
	s.logger.Info().Msg("destroy mongodb sink")
	s.stop <- true
	s.flush()
	s.logger.Info().Msg("destroy mongodb sink flush complete")
	return s.dbClient.Disconnect(context.Background())
}

func (s *mongoSink) Name() string {
	return name
}

func (s *mongoSink) Arrived(ctx context.Context, events ...*ce.Event) connector.Result {
	for idx := range events {
		e := events[idx]
		coll, err := getAttr(e, mongoCollection, s.cfg.Collection)
		if err != nil {
			return cdkgo.NewResult(http.StatusBadRequest, err.Error())
		}
		collName := coll
		var data interface{}
		if err := json.Unmarshal(e.Data(), &data); err != nil {
			return cdkgo.NewResult(http.StatusBadRequest, err.Error())
		}
		writer := s.getWriter(collName)
		writer.Write(data)
	}
	return cdkgo.SuccessResult
}

func (s *mongoSink) getWriter(collName string) *InsertWriter {
	s.lock.Lock()
	defer s.lock.Unlock()
	writer, ok := s.writer[collName]
	if !ok {
		writer = NewInsertWriter(s.dbClient, s.logger, s.cfg.Database, collName, s.cfg.BulkSize)
		s.writer[collName] = writer
	}
	return writer
}

func getAttr(e *ce.Event, key, defaultValue string) (string, error) {
	val, exist := e.Extensions()[key]
	if val == nil && !exist {
		return defaultValue, nil
	}

	str, ok := val.(string)
	if !ok {
		return "", fmt.Errorf("mongodb: invalid attribute %s=%v", key, val)
	}
	return str, nil
}

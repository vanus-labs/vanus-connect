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
	"encoding/json"
	"fmt"
	"github.com/linkall-labs/cdk-go/log"
	cdkutil "github.com/linkall-labs/cdk-go/utils"
	"net/http"
	"strings"
	"time"

	v2 "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/protocol"
	cehttp "github.com/cloudevents/sdk-go/v2/protocol/http"
	"github.com/golang/protobuf/jsonpb"
	"github.com/linkall-labs/cdk-go/runtime"
	proto "github.com/linkall-labs/connector/mongodb-sink/proto/database"
	"github.com/linkall-labs/connector/proto/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	mongoSinkDatabase   = "vancemongosinkdatabase"
	mongoSinkCollection = "vancemongosinkcollection"
)

type Config struct {
	Port    int      `json:"port" yaml:"port"`
	DBHosts []string `json:"db_hosts" yaml:"db_hosts"`
	Secret  *Secret  `json:"-" yaml:"-"`
}

type Secret struct {
	Username   string `json:"username" yaml:"username"`
	Password   string `json:"password" yaml:"password"`
	AuthSource string `json:"authSource" yaml:"authSource"`
}

func (sc Secret) isSet() bool {
	return sc.Username != "" || sc.Password != "" || sc.AuthSource != ""
}

type sink struct {
	cfg      *Config
	dbClient *mongo.Client
	logger   log.Logger
}

func NewMongoSink() runtime.Sink {
	return &sink{}
}

func (s *sink) SetLogger(logger log.Logger) {
	s.logger = logger
}

func (s *sink) Init(cfgPath, secretPath string) error {
	cfg := &Config{}
	if err := cdkutil.ParseConfig(cfgPath, cfg); err != nil {
		return err
	}

	if runtime.IsSecretEnable() {
		secret := &Secret{}
		if err := cdkutil.ParseConfig(secretPath, secret); err != nil {
			return err
		}
		cfg.Secret = secret
	}

	s.cfg = cfg
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	opts := options.Client()
	opts.ApplyURI(fmt.Sprintf("mongodb://%s", strings.Join(s.cfg.DBHosts, ",")))
	opts.SetReplicaSet("replicaset-01")
	if s.cfg.Secret != nil && s.cfg.Secret.isSet() {
		opts.Auth = &options.Credential{
			AuthSource:  s.cfg.Secret.AuthSource,
			Username:    s.cfg.Secret.Username,
			Password:    s.cfg.Secret.Password,
			PasswordSet: false,
		}
	}
	mongoClient, err := mongo.Connect(ctx, opts)
	if err != nil {
		panic("connect to mongodb error: " + err.Error())
	}
	s.dbClient = mongoClient
	return nil
}

func (s *sink) Destroy() error {
	return nil
}

func (s *sink) Name() string {
	return "mongodb-sink"
}

func (s *sink) Port() int {
	return s.cfg.Port
}

func (s *sink) Handle(ctx context.Context, event v2.Event) protocol.Result {
	extensions := event.Extensions()
	dbName, exist := extensions[mongoSinkDatabase]
	if !exist {
		return cehttp.NewResult(http.StatusBadRequest, "vancemongosinkdatabase is empty")
	}

	collName, exist := extensions[mongoSinkCollection]
	if !exist {
		return cehttp.NewResult(http.StatusBadRequest, "vancemongosinkcollection is empty")
	}
	id, err := primitive.ObjectIDFromHex(event.ID())
	if err != nil {
		return cehttp.NewResult(http.StatusBadRequest,
			fmt.Sprintf("invalid id %s, hex mongo id required", err))
	}

	e := &proto.Event{}
	if err := jsonpb.Unmarshal(bytes.NewReader(event.Data()), e); err != nil {
		return cehttp.NewResult(http.StatusBadRequest, err.Error())
	}

	if err := s.validate(e); err != nil {
		return cehttp.NewResult(http.StatusBadRequest, err.Error())
	}

	_dbName := fmt.Sprintf("%s", dbName)
	_collName := fmt.Sprintf("%s", collName)
	switch e.Op {
	case database.Operation_INSERT:
		err = s.insert(ctx, id, _dbName, _collName, e)
	case database.Operation_UPDATE:
		err = s.update(ctx, id, _dbName, _collName, e)
	case database.Operation_DELETE:
		err = s.delete(ctx, id, _dbName, _collName)
	default:
		return cehttp.NewResult(http.StatusBadRequest, fmt.Sprintf("unsupported event operation: %s", e.Op))
	}

	if err != nil {
		return cehttp.NewResult(http.StatusInternalServerError, err.Error())
	}
	return cehttp.NewResult(http.StatusOK, "")
}

func (s *sink) insert(ctx context.Context, id primitive.ObjectID, dbName, collName string, e *proto.Event) error {
	m := make(map[string]interface{})
	data, err := e.Insert.Document.MarshalJSON()
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	delete(m, "_id")
	m["_id"] = id
	_, err = s.dbClient.Database(dbName).Collection(collName).InsertOne(ctx, m)
	return err
}

func (s *sink) update(ctx context.Context, id primitive.ObjectID, dbName, collName string, e *proto.Event) error {
	data, err := e.Update.UpdateDescription.UpdatedFields.MarshalJSON()
	if err != nil {
		return fmt.Errorf("try to marhsall UpdatedFields error: %s", err)
	}
	updates := make(map[string]interface{})
	if err = json.Unmarshal(data, &updates); err != nil {
		return fmt.Errorf("try to unmarhsall UpdatedFields data to map error: %s", err)
	}
	ctx, cancel := context.WithTimeout(ctx, 3*time.Minute)
	defer cancel()
	_, err = s.dbClient.Database(dbName).Collection(collName).UpdateOne(ctx,
		bson.M{"_id": id},
		bson.M{
			"$set": updates,
		})
	return err
}

func (s *sink) delete(ctx context.Context, id primitive.ObjectID, dbName, collName string) error {
	_, err := s.dbClient.Database(dbName).Collection(collName).DeleteOne(ctx, bson.M{
		"_id": id,
	})
	return err
}

func (s *sink) validate(event *proto.Event) error {
	if event.Op == database.Operation_INSERT && event.Insert == nil {
		return fmt.Errorf("invalid insert event")
	}

	if event.Op == database.Operation_UPDATE && event.Update == nil {
		return fmt.Errorf("invalid update event")
	}

	return nil
}

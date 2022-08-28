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
	"net/http"
	"strings"
	"time"

	v2 "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/protocol"
	cehttp "github.com/cloudevents/sdk-go/v2/protocol/http"
	"github.com/golang/protobuf/jsonpb"
	"github.com/linkall-labs/connector/cdk-go/runtime"
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
	Port    int      `json:"port"`
	DBHosts []string `json:"db_hosts"`
	Secret  Secret   `json:"-"`
}

type Secret struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	AuthSource string `json:"auth_source"`
}

func (sc Secret) isSet() bool {
	return sc.Username != "" || sc.Password != "" || sc.AuthSource != ""
}

type sink struct {
	cfg      *Config
	dbClient *mongo.Client
}

func NewMongoSink() runtime.Sink {
	return &sink{}
}

func (s *sink) Init(cfgPath, secretPath string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	opts := options.Client()
	opts.ApplyURI(fmt.Sprintf("mongodb://%s", strings.Join(s.cfg.DBHosts, ",")))
	if s.cfg.Secret.isSet() {
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

func (s *sink) Handle(ctx context.Context, event *v2.Event) protocol.Result {
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

	switch e.Op {
	case database.Operation_INSERT:
		err = s.insert(ctx, fmt.Sprintf("%s", dbName), fmt.Sprintf("%s", collName), e)
	case database.Operation_UPDATE:
		err = s.update(ctx, id, fmt.Sprintf("%s", dbName), fmt.Sprintf("%s", collName), e)
	case database.Operation_DELETE:
		err = s.delete(ctx, id, fmt.Sprintf("%s", dbName), fmt.Sprintf("%s", collName))
	default:
		return cehttp.NewResult(http.StatusBadRequest, fmt.Sprintf("unsupported event operation: %s", e.Op))
	}

	if err != nil {
		return cehttp.NewResult(http.StatusInternalServerError, err.Error())
	}
	return cehttp.NewResult(http.StatusOK, "")
}

func (s *sink) insert(ctx context.Context, dbName, collName string, e *proto.Event) error {
	_, err := s.dbClient.Database(dbName).Collection(collName).InsertOne(ctx, e.Insert.Document)
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

	_, err = s.dbClient.Database(dbName).Collection(collName).UpdateOne(ctx,
		bson.M{"_id": id},
		updates)
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

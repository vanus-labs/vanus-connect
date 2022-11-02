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
	"github.com/linkall-labs/cdk-go/connector"
	"github.com/linkall-labs/cdk-go/connector/sink"
	"net/http"
	"strings"
	"time"

	cehttp "github.com/cloudevents/sdk-go/v2/protocol/http"
	stdlib "github.com/golang/protobuf/proto"
	"github.com/linkall-labs/cdk-go/log"
	cdkutil "github.com/linkall-labs/cdk-go/utils"
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
	Port       int      `json:"port" yaml:"port"`
	DBHosts    []string `json:"db_hosts" yaml:"db_hosts"`
	ReplicaSet string   `json:"replica_set" yaml:"replica_set"`
	Secret     *Secret  `json:"-" yaml:"-"`
}

type Secret struct {
	Username   string `json:"username" yaml:"username"`
	Password   string `json:"password" yaml:"password"`
	AuthSource string `json:"authSource" yaml:"authSource"`
}

func (sc Secret) isSet() bool {
	return sc.Username != "" || sc.Password != "" || sc.AuthSource != ""
}

type mongoSink struct {
	cfg      *Config
	dbClient *mongo.Client
	logger   log.Logger
}

func NewMongoSink() sink.ProtobufSink {
	return &mongoSink{}
}

func (s *mongoSink) SetLogger(logger log.Logger) {
	s.logger = logger
}

func (s *mongoSink) Init(cfgPath, secretPath string) error {
	cfg := &Config{}
	if err := cdkutil.ParseConfig(cfgPath, cfg); err != nil {
		return err
	}

	if connector.IsSecretEnable() {
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
	opts.SetReplicaSet(s.cfg.ReplicaSet)
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

func (s *mongoSink) NewEvent() stdlib.Message {
	return &proto.Event{}
}

func (s *mongoSink) Validate(msg stdlib.Message) error {
	event, ok := msg.(*proto.Event)
	if !ok {
		return fmt.Errorf("invalid event type")
	}
	if event.Op == database.Operation_INSERT && event.Insert == nil {
		return fmt.Errorf("invalid insert event")
	}

	if event.Op == database.Operation_UPDATE && event.Update == nil {
		return fmt.Errorf("invalid update event")
	}

	_, exist := event.Metadata.Extension.GetFields()[mongoSinkDatabase]
	if !exist {
		return fmt.Errorf("vancemongosinkdatabase is empty")
	}

	_, exist = event.Metadata.Extension.GetFields()[mongoSinkCollection]
	if !exist {
		return fmt.Errorf("vancemongosinkcollection is empty")
	}

	_, err := primitive.ObjectIDFromHex(event.Metadata.Id)
	if err != nil {
		return fmt.Errorf("invalid id %s, hex mongo id required", err)
	}

	return nil
}

func (s *mongoSink) Destroy() error {
	return s.dbClient.Disconnect(context.TODO())
}

func (s *mongoSink) Name() string {
	return "mongodb-mongoSink"
}

func (s *mongoSink) Port() int {
	return s.cfg.Port
}

func (s *mongoSink) Handle(ctx context.Context, msg stdlib.Message) error {
	event, _ := msg.(*proto.Event)

	id, _ := primitive.ObjectIDFromHex(event.Metadata.Id)
	dbName := event.Metadata.Extension.GetFields()[mongoSinkDatabase].GetStringValue()
	collName := event.Metadata.Extension.GetFields()[mongoSinkCollection].GetStringValue()
	var err error
	switch event.Op {
	case database.Operation_INSERT:
		err = s.insert(ctx, id, dbName, collName, event)
	case database.Operation_UPDATE:
		err = s.update(ctx, id, dbName, collName, event)
	case database.Operation_DELETE:
		err = s.delete(ctx, id, dbName, collName)
	default:
		return cehttp.NewResult(http.StatusBadRequest, fmt.Sprintf("unsupported event operation: %s", event.Op))
	}

	return err
}

func (s *mongoSink) insert(ctx context.Context, id primitive.ObjectID, dbName, collName string, e *proto.Event) error {
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

func (s *mongoSink) update(ctx context.Context, id primitive.ObjectID, dbName, collName string, e *proto.Event) error {
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

func (s *mongoSink) delete(ctx context.Context, id primitive.ObjectID, dbName, collName string) error {
	_, err := s.dbClient.Database(dbName).Collection(collName).DeleteOne(ctx, bson.M{
		"_id": id,
	})
	return err
}

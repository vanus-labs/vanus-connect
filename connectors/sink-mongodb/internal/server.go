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
	"time"

	ce "github.com/cloudevents/sdk-go/v2"
	cdkgo "github.com/linkall-labs/cdk-go"
	"github.com/linkall-labs/cdk-go/config"
	"github.com/linkall-labs/cdk-go/connector"
	"github.com/linkall-labs/cdk-go/log"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	mongoSinkDatabase   = "xvdatabasedb"
	mongoSinkCollection = "xvdatabasecoll"

	name = "Sink MongoDB"
)

var _ cdkgo.SinkConfigAccessor = &Config{}

type Config struct {
	cdkgo.SinkConfig      `json:",inline" yaml:",inline"`
	ConnectionURI         string         `json:"connection_uri" yaml:"connection_uri"`
	Credential            Credential     `json:"credential" yaml:"credential"`
	ConvertConfig         *ConvertConfig `json:"convert" yaml:"convert"`
	IgnoreDuplicatedError bool           `json:"ignore_duplicated_error" yaml:"ignore_duplicated_error"`
}

type Credential struct {
	Username                string            `json:"username" yaml:"username"`
	Password                string            `json:"password" yaml:"password"`
	AuthSource              string            `json:"auth_source" yaml:"auth_source"`
	AuthMechanism           string            `json:"auth_mechanism" yaml:"auth_mechanism"`
	AuthMechanismProperties map[string]string `json:"auth_mechanism_properties" yaml:"auth_mechanism_properties"`
}

func (c *Config) Validate() error {
	if c.ConvertConfig != nil {
		err := c.ConvertConfig.Validate()
		if err != nil {
			return err
		}
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

func (c *Config) GetSecret() cdkgo.SecretAccessor {
	return &c.Credential
}

var _ cdkgo.Sink = &mongoSink{}

type mongoSink struct {
	cfg           *Config
	dbClient      *mongo.Client
	logger        log.Logger
	convertStruct *convertStruct
}

func NewMongoSink() cdkgo.Sink {
	return &mongoSink{}
}

func (s *mongoSink) Initialize(ctx context.Context, cfg config.ConfigAccessor) error {
	c, ok := cfg.(*Config)
	if !ok {
		return errors.New("unexpected config type")
	}

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
		log.Info("mongodb is connected", map[string]interface{}{
			"url": s.cfg.ConnectionURI,
		})
	}
	s.dbClient = mongoClient
	if s.cfg.ConvertConfig != nil {
		s.convertStruct = newConvert(s.cfg.ConvertConfig)
	}
	return nil
}

func (s *mongoSink) Destroy() error {
	log.Info("mongodb connection is closing", nil)
	return s.dbClient.Disconnect(context.Background())
}

func (s *mongoSink) Name() string {
	return name
}

type Command struct {
	Inserts []interface{} `json:"inserts"`
	Updates []Update      `json:"updates"`
	Deletes []Delete      `json:"deletes"`
}

type Update struct {
	Filter     bson.M `json:"filter"`
	Update     bson.M `json:"update"`
	UpdateMany bool   `json:"update_many"`
}

type Delete struct {
	Filter     bson.M `json:"filter"`
	DeleteMany bool   `json:"delete_many"`
}

func (s *mongoSink) Arrived(ctx context.Context, events ...*ce.Event) connector.Result {
	var err error
	events, err = s.convertEvents(events...)
	if err != nil {
		return cdkgo.NewResult(http.StatusBadRequest, err.Error())
	}
	var db string
	var coll string
	wms := make([]mongo.WriteModel, 0)
	var inserts []interface{}
	for idx := range events {
		e := events[idx]

		if db == "" {
			_db, err := getAttr(e, mongoSinkDatabase)
			if err != nil {
				return cdkgo.NewResult(http.StatusBadRequest, err.Error())
			}
			db = _db
		}

		if coll == "" {
			_coll, err := getAttr(e, mongoSinkCollection)
			if err != nil {
				return cdkgo.NewResult(http.StatusBadRequest, err.Error())
			}
			coll = _coll
		}

		c := Command{}
		if err := json.Unmarshal(e.Data(), &c); err != nil {
			return cdkgo.NewResult(http.StatusBadRequest, err.Error())
		}

		inserts = append(inserts, c.Inserts...)

		wm := getUpdateOrDeleteModels(&c)
		if len(wm) > 0 {
			wms = append(wms, wm...)
		}
	}

	if len(inserts) > 0 {
		opt := options.InsertMany().SetOrdered(false)
		if res, err := s.dbClient.Database(db).Collection(coll).InsertMany(ctx, inserts, opt); err != nil {

			log.Warning("failed to insert many to mongodb", map[string]interface{}{
				log.KeyError: err,
				"inserted":   res.InsertedIDs,
			})

			if err != nil {
				if !mongo.IsDuplicateKeyError(err) ||
					(mongo.IsDuplicateKeyError(err) && !s.cfg.IgnoreDuplicatedError) {
					return cdkgo.NewResult(http.StatusInternalServerError,
						fmt.Sprintf("failed to insert many to mongodb: %s", err))
				}
			}
		}
	}

	if len(wms) > 0 {
		if _, err := s.dbClient.Database(db).Collection(coll).BulkWrite(ctx, wms); err != nil {
			log.Warning("failed to write mongodb", map[string]interface{}{
				log.KeyError: err,
			})
			return cdkgo.NewResult(http.StatusInternalServerError,
				fmt.Sprintf("failed to write mongodb: %s", err))
		}
	}

	return cdkgo.SuccessResult
}

func getAttr(e *ce.Event, key string) (string, error) {
	val, exist := e.Extensions()[key]
	if val == nil && !exist {
		return "", fmt.Errorf("mongodb: attribute %s not found or is empty", key)
	}

	str, ok := val.(string)
	if !ok {
		return "", fmt.Errorf("mongodb: invalid attribute %s=%v", key, val)
	}
	return str, nil
}

func getUpdateOrDeleteModels(c *Command) []mongo.WriteModel {
	wms := make([]mongo.WriteModel, len(c.Updates)+len(c.Deletes))
	var i = 0

	for idx := range c.Updates {
		u := c.Updates[idx]
		if u.UpdateMany {
			wms[i] = &mongo.UpdateManyModel{
				Filter: u.Filter,
				Update: u.Update,
			}
		} else {
			wms[i] = &mongo.UpdateOneModel{
				Filter: u.Filter,
				Update: u.Update,
			}
		}
		i++
	}

	for idx := range c.Deletes {
		d := c.Deletes[idx]
		if d.DeleteMany {
			wms[i] = &mongo.DeleteManyModel{
				Filter: d.Filter,
			}
		} else {
			wms[i] = &mongo.DeleteOneModel{
				Filter: d.Filter,
			}
		}
		i++
	}
	return wms
}

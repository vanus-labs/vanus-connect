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
	"strconv"
	"strings"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo/writeconcern"

	ce "github.com/cloudevents/sdk-go/v2"
	"github.com/pkg/errors"
	cdkgo "github.com/vanus-labs/cdk-go"
	"github.com/vanus-labs/cdk-go/config"
	"github.com/vanus-labs/cdk-go/connector"
	"github.com/vanus-labs/cdk-go/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.uber.org/ratelimit"
)

const (
	mongoSinkDatabase   = "xvdatabasedb"
	mongoSinkCollection = "xvdatabasecoll"

	name = "Sink MongoDB"

	debeziumConnector = "iodebeziumconnector"
)

var _ cdkgo.SinkConfigAccessor = &Config{}

type Config struct {
	cdkgo.SinkConfig      `json:",inline" yaml:",inline"`
	ConnectionURI         string         `json:"connection_uri" yaml:"connection_uri"`
	Credential            Credential     `json:"credential" yaml:"credential"`
	ConvertConfig         *ConvertConfig `json:"convert" yaml:"convert"`
	IgnoreDuplicatedError bool           `json:"ignore_duplicated_error" yaml:"ignore_duplicated_error"`
	Parallelism           int            `json:"parallelism" yaml:"parallelism"`
	InsertWriteConcern    string         `json:"insert_write_concern" yaml:"insert_write_concern"`
	RateLimit             int            `json:"rate_limit" yaml:"rate_limit"`
	DebugSkip             bool           `json:"debug_skip" yaml:"debug_skip"`
	BulkSize              int            `json:"bulk_size" yaml:"bulk_size"`
	Upsert                bool           `json:"upsert" yaml:"upsert"`
	ackDisable            bool
	writeCon              writeconcern.Option
}

func (c *Config) GetWriteConcern() writeconcern.Option {
	return c.writeCon
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
	if c.ConvertConfig != nil {
		err := c.ConvertConfig.Validate()
		if err != nil {
			return err
		}
	}
	if c.Parallelism == 0 {
		c.Parallelism = 1
	}

	if c.Parallelism > 32 {
		c.Parallelism = 32
		log.Info("the parallelism is exceeded than the maximum value of 32, set it to 32", nil)
	}

	var opt = writeconcern.WMajority()
	if c.InsertWriteConcern != "" {
		w := strings.Split(c.InsertWriteConcern, ":")
		if len(w) == 2 && w[0] == "w" {
			i, err := strconv.ParseInt(w[1], 10, 64)
			if err != nil || i > 3 {
				log.Info("invalid insert_write_concern, use majority", map[string]interface{}{
					"config":     c.InsertWriteConcern,
					"ref":        "https://www.mongodb.com/docs/manual/reference/write-concern/",
					log.KeyError: err,
				})
			}
			if i == 0 {
				log.Warning("insert unacknowledged is enable, watch carefully your mongodb instance ", nil)
				c.ackDisable = true
			}
			opt = writeconcern.W(int(i))
		}
	} else {
		log.Info("use default writeConcern: majority", nil)
	}
	c.writeCon = opt

	if c.RateLimit == 0 {
		if c.ackDisable {
			c.RateLimit = 5000
		} else {
			c.RateLimit = 1 << 20
		}
	}

	log.Info("config", map[string]interface{}{
		"rate_limit": c.RateLimit,
	})

	if c.BulkSize == 0 {
		c.BulkSize = 4
	}

	log.Info("config", map[string]interface{}{
		"bulk_size": c.BulkSize,
	})

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
	cfg           *Config
	dbClient      *mongo.Client
	logger        log.Logger
	convertStruct *convertStruct
	db            map[string]map[string]*mongo.Collection
	mutex         sync.RWMutex
	rateLimit     ratelimit.Limiter
}

func NewMongoSink() cdkgo.Sink {
	return &mongoSink{
		db: map[string]map[string]*mongo.Collection{},
	}
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
		log.Info("mongodb is connected", nil)
	}
	s.dbClient = mongoClient
	if s.cfg.ConvertConfig != nil {
		s.convertStruct = newConvert(s.cfg.ConvertConfig)
	}
	s.rateLimit = ratelimit.New(s.cfg.RateLimit)

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
	Filter     map[string]interface{} `json:"filter"`
	Update     map[string]interface{} `json:"update"`
	UpdateMany bool                   `json:"update_many"`
}

type Delete struct {
	Filter     bson.M `json:"filter"`
	DeleteMany bool   `json:"delete_many"`
}

func (s *mongoSink) Arrived(ctx context.Context, events ...*ce.Event) connector.Result {
	if s.cfg.DebugSkip {
		return cdkgo.SuccessResult
	}
	var dbName string
	var collName string
	wms := make([]mongo.WriteModel, 0)
	start := time.Now()
	var inserts []interface{}
	for idx := range events {
		s.rateLimit.Take()
		e := events[idx]

		v, exist := e.Extensions()[debeziumConnector]
		if exist {
			switch v {
			case "mysql":
				_e, err := s.convertEvents(e)
				if err != nil {
					return cdkgo.NewResult(http.StatusBadRequest, err.Error())
				}
				e = _e[0]
			}
		}

		if dbName == "" {
			_db, err := getAttr(e, mongoSinkDatabase)
			if err != nil {
				return cdkgo.NewResult(http.StatusBadRequest, err.Error())
			}
			dbName = _db
		}

		if collName == "" {
			_coll, err := getAttr(e, mongoSinkCollection)
			if err != nil {
				return cdkgo.NewResult(http.StatusBadRequest, err.Error())
			}
			collName = _coll
		}

		c := Command{}
		d := json.NewDecoder(bytes.NewReader(e.Data()))
		d.UseNumber()

		if err := d.Decode(&c); err != nil {
			return cdkgo.NewResult(http.StatusBadRequest, err.Error())
		}

		c.Inserts = recursive(c.Inserts).([]interface{})
		inserts = append(inserts, c.Inserts...)

		c.Updates = recursive(c.Updates).([]Update)
		c.Deletes = recursive(c.Deletes).([]Delete)
		wm := getUpdateOrDeleteModels(&c, s.cfg.Upsert)
		if len(wm) > 0 {
			wms = append(wms, wm...)
		}
	}

	cur := time.Now()
	if cur.Sub(start) > 10*time.Millisecond {
		log.Info("preparing data takes too long", map[string]interface{}{
			"numbers": len(inserts),
			"used":    cur.Sub(start),
		})
	}

	start = time.Now()
	getColl := func() *mongo.Collection {
		d, exist := s.db[dbName]
		if !exist {
			return nil
		}
		return d[collName]
	}

	s.mutex.RLock()
	collInstance := getColl()
	s.mutex.RUnlock()

	if collInstance == nil {
		s.mutex.Lock()
		collInstance = getColl()
		if collInstance == nil {
			collOpt := options.Collection().SetWriteConcern(writeconcern.New(s.cfg.GetWriteConcern()))
			d, exist := s.db[dbName]
			if !exist {
				d = make(map[string]*mongo.Collection)
				s.db[dbName] = d
			}
			collInstance = s.dbClient.Database(dbName).Collection(collName, collOpt)
			d[collName] = collInstance
		}
		s.mutex.Unlock()
	}

	cur = time.Now()
	if cur.Sub(start) > 10*time.Millisecond {
		log.Info("take collection too long", map[string]interface{}{
			"numbers": len(inserts),
			"used":    cur.Sub(start),
		})
	}

	start = time.Now()
	var err error
	if len(inserts) > 0 {
		para := s.cfg.Parallelism
		if para > len(inserts) {
			para = len(inserts)
		}
		avg := len(inserts) / para
		from := 0
		end := from + avg
		wg := sync.WaitGroup{}
		mutex := sync.Mutex{}
		insert := func(data []interface{}) {
			defer wg.Done()
			opt := options.InsertMany().SetOrdered(false)

			if res, _err := collInstance.InsertMany(ctx, data, opt); _err != nil {
				if _err == mongo.ErrUnacknowledgedWrite && s.cfg.ackDisable {
					return
				}
				log.Warning("failed to insert many to mongodb", map[string]interface{}{
					log.KeyError: _err,
					"inserted":   res,
				})
				mutex.Lock()
				if err == nil {
					err = _err
				} else if mongo.IsDuplicateKeyError(err) && !mongo.IsDuplicateKeyError(_err) {
					err = _err
				}
				mutex.Unlock()
			}
		}

		for idx := 0; idx < para; idx++ {
			wg.Add(1)

			if idx == para-1 {
				go insert(inserts[from:])
			} else {
				go insert(inserts[from:end])
			}
			from = end
			end = from + avg
		}
		wg.Wait()
	}

	if err != nil {
		if !mongo.IsDuplicateKeyError(err) ||
			(mongo.IsDuplicateKeyError(err) && !s.cfg.IgnoreDuplicatedError) {
			return cdkgo.NewResult(http.StatusInternalServerError,
				fmt.Sprintf("failed to insert many to mongodb: %s", err))
		}
	}

	cur = time.Now()
	if cur.Sub(start) > 100*time.Millisecond {
		log.Info("insert mongodb takes too long", map[string]interface{}{
			"numbers": len(inserts),
			"used":    cur.Sub(start),
		})
	}

	start = time.Now()
	for from := 0; from < len(wms); {
		var models []mongo.WriteModel
		end := from + s.cfg.BulkSize
		if end >= len(wms) {
			models = wms[from:]
		} else {
			models = wms[from:end]
		}
		if _, err := s.dbClient.Database(dbName).Collection(collName).BulkWrite(ctx, models); err != nil {
			log.Warning("failed to write mongodb", map[string]interface{}{
				log.KeyError: err,
			})
			return cdkgo.NewResult(http.StatusInternalServerError,
				fmt.Sprintf("failed to write mongodb: %s", err))
		}
		from = end
	}

	cur = time.Now()
	if cur.Sub(start) > 100*time.Millisecond {
		log.Info("update to mongodb takes too long", map[string]interface{}{
			"numbers": len(wms),
			"used":    cur.Sub(start),
		})
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

func getUpdateOrDeleteModels(c *Command, upsert bool) []mongo.WriteModel {
	wms := make([]mongo.WriteModel, len(c.Updates)+len(c.Deletes))
	var i = 0

	for idx := range c.Updates {
		u := c.Updates[idx]
		if u.UpdateMany {
			wms[i] = &mongo.UpdateManyModel{
				Filter: u.Filter,
				Update: u.Update,
				Upsert: &upsert,
			}
		} else {
			wms[i] = &mongo.UpdateOneModel{
				Filter: u.Filter,
				Update: u.Update,
				Upsert: &upsert,
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

func recursive(val interface{}) interface{} {
	switch val.(type) {
	case json.Number:
		v := val.(json.Number)
		if i, err := v.Int64(); err == nil {
			val = i
		} else {
			val, _ = v.Float64()
		}
	case map[string]interface{}:
		v := val.(map[string]interface{})
		for k := range v {
			v[k] = recursive(v[k])
		}
	case []interface{}:
		v := val.([]interface{})
		for idx := range v {
			v[idx] = recursive(v[idx])
		}
	case []Update:
		v := val.([]Update)
		for idx := range v {
			v[idx] = recursive(v[idx]).(Update)
		}
	case Update:
		v := val.(Update)
		v.Filter = recursive(v.Filter).(map[string]interface{})
		v.Update = recursive(v.Update).(map[string]interface{})
	case []Delete:
		v := val.([]Delete)
		for idx := range v {
			v[idx] = recursive(v[idx]).(Delete)
		}
	case Delete:
		v := val.(Delete)
		v.Filter = recursive(v.Filter).(map[string]interface{})
	default:
		return val
	}
	return val
}

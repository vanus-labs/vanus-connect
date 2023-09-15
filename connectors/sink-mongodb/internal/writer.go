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
	"sync"

	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/mongo"
)

type InsertWriter struct {
	lock      sync.Mutex
	data      []interface{}
	size      int
	flushSize int
	coll      *mongo.Collection
	logger    zerolog.Logger
}

func NewInsertWriter(dbClient *mongo.Client, logger zerolog.Logger, dbName, collName string, flushSize int) *InsertWriter {
	return &InsertWriter{
		data:      make([]interface{}, 0),
		coll:      dbClient.Database(dbName).Collection(collName),
		logger:    logger,
		flushSize: flushSize,
	}
}

func (w *InsertWriter) Size() int {
	w.lock.Lock()
	defer w.lock.Unlock()
	return w.size
}

func (w *InsertWriter) Write(data interface{}) {
	w.lock.Lock()
	defer w.lock.Unlock()
	w.data = append(w.data, data)
	w.size++
	if w.size >= w.flushSize {
		err := w.flush()
		w.logger.Warn().Err(err).Str("collection", w.coll.Name()).Msg("insert failed")
	}
}

func (w *InsertWriter) Flush() error {
	w.lock.Lock()
	defer w.lock.Unlock()
	if w.size == 0 {
		return nil
	}
	return w.flush()
}

func (w *InsertWriter) flush() error {
	_, err := w.coll.InsertMany(context.TODO(), w.data)
	if err != nil {
		return err
	}
	w.logger.Info().Int("size", w.size).Str("collection", w.coll.Name()).Msg("insert success")
	w.data = nil
	w.size = 0
	return nil
}

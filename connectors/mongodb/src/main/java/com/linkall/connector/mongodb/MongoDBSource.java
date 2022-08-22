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

package com.linkall.connector.mongodb;

import com.linkall.vance.core.Adapter;

import java.util.HashMap;
import java.util.Map;
import java.util.Properties;

public class MongoDBSource extends com.linkall.connector.mongodb.debezium.Source implements com.linkall.vance.core.Source {

    private static final String DEBEZIUM_CONNECTOR = "io.debezium.connector.mongodb.MongoDbConnector";

    @Override
    public Adapter getAdapter() {
        return new MongoDBAdapter();
    }

    @Override
    public String getConnectorClass() {
        return DEBEZIUM_CONNECTOR;
    }

    @Override
    public Map<String, Object> getConfigOffset() {
        return new HashMap<>();
    }

    @Override
    // https://debezium.io/documentation/reference/stable/connectors/mongodb.html#mongodb-connector-properties
    public Properties getDebeziumProperties() {
        final Properties props = new Properties();

        props.setProperty("name", "test");
        props.setProperty("mongodb.hosts", "44.242.140.28:27017");
        props.setProperty("mongodb.name", "test");
//        props.setProperty("mongodb.user", "admin");
//        props.setProperty("mongodb.password", "admin");
//        props.setProperty("mongodb.authsource", "admin");

        return props;
    }
}

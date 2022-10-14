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

package com.linkall.connector.source.mongodb;

import com.alibaba.fastjson.JSON;
import com.alibaba.fastjson.JSONObject;
import com.linkall.vance.core.Adapter;
import org.apache.logging.log4j.util.Strings;

import java.io.FileNotFoundException;
import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.HashMap;
import java.util.Map;
import java.util.Properties;
import java.util.stream.Collectors;
import java.util.stream.Stream;

public class MongoDBSource extends com.linkall.cdk.database.debezium.DebeziumSource
        implements com.linkall.vance.core.Source {
    private static final String DEBEZIUM_CONNECTOR = "io.debezium.connector.mongodb.MongoDbConnector";

    private final String connectorName;
    private final Map<String, Object> config;
    private Map<String, String> secret;

    public MongoDBSource() throws IOException {
        super();
        String configPath = System.getenv("CONNECTOR_CONFIG");
        if (Strings.isBlank(configPath)) {
            configPath = "/vance/config/config.json";
        }
        Path cp = Paths.get(configPath);
        if (!cp.toFile().exists()) {
            throw new FileNotFoundException("the config.json not found in [ " + cp.toAbsolutePath() + " ]");
        }
        byte[] data = Files.readAllBytes(cp);
        config = JSON.parseObject(data, Map.class);

        String secretPath = System.getenv("CONNECTOR_SECRET");
        if (Strings.isBlank(secretPath)) {
            secretPath = "/vance/config/secret.json";
        }

        Path sp = Paths.get(secretPath);
        if (sp.toFile().exists()) {
            data = Files.readAllBytes(sp);
            secret = JSON.parseObject(data, Map.class);
        }
        connectorName = config.get("name").toString();
    }

    @Override
    public Adapter getAdapter() {
        return new MongoDBAdapter();
    }

    @Override
    public String getConnectorClass() {
        return DEBEZIUM_CONNECTOR;
    }

    @Override
    public String getDatabase() {
        return connectorName;
    }

    @Override
    public String getStoreOffsetKey() {
        return null;
    }

    @Override
    public Map<String, Object> getConfigOffset() {
        return new HashMap<>();
    }

    @Override
    // https://debezium.io/documentation/reference/stable/connectors/mongodb.html#mongodb-connector-properties
    public Properties getDebeziumProperties() {
        final Properties props = new Properties();
        props.setProperty("name", connectorName);
        props.setProperty("mongodb.hosts", config.get("db_hosts").toString());
        props.setProperty("mongodb.name", config.get("db_name").toString());
        props.setProperty("capture.mode", "change_streams_update_full");
        if (secret != null && secret.size() > 0) {
            props.setProperty("mongodb.user", secret.get("username"));
            props.setProperty("mongodb.password", secret.get("password"));
            props.setProperty("mongodb.authsource", secret.getOrDefault("authSource", "admin"));
        }

        if (config.get("database") != null) {
            JSONObject db = (JSONObject) config.get("database");
            if (db.getJSONArray("include").size() > 0 && db.getJSONArray("exclude").size() > 0) {
                throw new IllegalArgumentException("the database.include and database.exclude can't be set together");
            }
            if (db.getJSONArray("include").size() > 0) {
                props.setProperty("database.include.list", tableFormat(db.getJSONArray("include").stream()));
            }

            if (db.getJSONArray("exclude").size() > 0) {
                props.setProperty("database.exclude.list", tableFormat(db.getJSONArray("exclude").stream()));
            }
        }
        if (config.get("collection") != null) {
            JSONObject collection = (JSONObject) config.get("collection");
            if (collection.getJSONArray("include").size() > 0 && collection.getJSONArray("exclude").size() > 0) {
                throw new IllegalArgumentException("the collection.include and collection.exclude can't be set together");
            }
            if (collection.getJSONArray("include").size() > 0) {
                props.setProperty("collection.include.list", tableFormat(collection.getJSONArray("include").stream()));
            }

            if (collection.getJSONArray("exclude").size() > 0) {
                props.setProperty("collection.exclude.list", tableFormat(collection.getJSONArray("exclude").stream()));
            }
        }

        return props;
    }

    public String tableFormat(Stream<Object> table) {
        return table
                .map(stream -> this.getDatabase() + "." + stream)
                .collect(Collectors.joining(","));
    }
}

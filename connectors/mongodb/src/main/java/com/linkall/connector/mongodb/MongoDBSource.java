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

import com.alibaba.fastjson.JSON;
import com.linkall.connector.mongodb.debezium.DebeziumSource;
import com.linkall.vance.core.Adapter;
import org.apache.logging.log4j.util.Strings;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.FileNotFoundException;
import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.Arrays;
import java.util.HashMap;
import java.util.Map;
import java.util.Properties;
import java.util.stream.Collectors;
import java.util.stream.Stream;

public class MongoDBSource extends DebeziumSource implements com.linkall.vance.core.Source {

    private static final String DEBEZIUM_CONNECTOR = "io.debezium.connector.mongodb.MongoDbConnector";

    private String connectorName;
    private Map<String, Object> config;
    private Map<String, String> secret;

    public MongoDBSource() throws IOException {
        super();
        String home = System.getenv("MONGODB_CONNECTOR_HOME");
        if (Strings.isBlank(home)) {
            home = "/etc/vance/mongodb";
        }
        Path cp = Paths.get(home, "config.json");
        if (!cp.toFile().exists()) {
            throw new FileNotFoundException("the config.json not found in [ " + cp.toAbsolutePath().toString() + " ]");
        }
        byte[] data = Files.readAllBytes(cp);
        config = JSON.parseObject(data, Map.class);

        Path sp = Paths.get(home, "secret.json");
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
        if (secret.size() > 0) {
            props.setProperty("mongodb.user", secret.get("user"));
            props.setProperty("mongodb.password", secret.get("password"));
            props.setProperty("mongodb.authsource", secret.getOrDefault("authsource", "admin"));
        }

        if (config.get("database") != null) {
            Map<String, String[]> db = (Map<String, String[]>) config.get("database");
            if (db.get("include").length > 0 && db.get("exclude").length > 0) {
                throw new IllegalArgumentException("the database.include and database.exclude can't be set together");
            }
            if (db.get("include").length > 0) {
                props.setProperty("database.include.list", tableFormat(Arrays.stream(db.get("include"))));
            }

            if (db.get("exclude").length > 0) {
                props.setProperty("database.exclude.list", tableFormat(Arrays.stream(db.get("exclude"))));
            }
        }
        if (config.get("collection") != null) {
            Map<String, String[]> collection = (Map<String, String[]>) config.get("collection");
            if (collection.get("include").length > 0 && collection.get("exclude").length > 0) {
                throw new IllegalArgumentException("the collection.include and collection.exclude can't be set together");
            }
            if (collection.get("include").length > 0) {
                props.setProperty("collection.include.list", tableFormat(Arrays.stream(collection.get("include"))));
            }

            if (collection.get("exclude").length > 0) {
                props.setProperty("collection.exclude.list", tableFormat(Arrays.stream(collection.get("exclude"))));
            }
        }

        return props;
    }

    public String tableFormat(Stream<String> table) {
        return table
                .map(stream -> this.getDatabase() + "." + stream)
                .collect(Collectors.joining(","));
    }
}

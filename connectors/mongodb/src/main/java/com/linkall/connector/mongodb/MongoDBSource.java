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
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.*;
import java.util.stream.Collectors;
import java.util.stream.Stream;

public class MongoDBSource extends DebeziumSource implements com.linkall.vance.core.Source {
    private static final Logger LOGGER = LoggerFactory.getLogger(MongoDBSource.class);

    private static final String DEBEZIUM_CONNECTOR = "io.debezium.connector.mongodb.MongoDbConnector";
    private static final String CONNECTOR_NAME = System.getProperty("CONNECTOR_NAME", "test");
    private static final String MONGODB_NAME = System.getProperty("MONGODB_NAME", "test");
    private static final String MONGODB_HOSTS = System.getProperty("MONGODB_HOSTS", "");
    private static final String DB_INCLUDE_LIST = System.getProperty("DB_INCLUDE_LIST", "");
    private static String secretFile = System.getProperty("SECRET_FILE", "secret.json");

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
        return CONNECTOR_NAME;
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
    public Properties getDebeziumProperties() throws IOException {
        final Properties props = new Properties();
        props.setProperty("name", CONNECTOR_NAME);
        props.setProperty("database.server.name", CONNECTOR_NAME);
//        props.setProperty("mongodb.hosts", MONGODB_HOSTS);
        props.setProperty("mongodb.hosts", "44.242.140.28:27017");
        props.setProperty("mongodb.name", MONGODB_NAME);

        // table selection
        props.setProperty("database.include.list", DB_INCLUDE_LIST);
//        if (!config.getIncludeTables().isEmpty()) {
//            props.setProperty("table.include.list", tableFormat(config.getIncludeTables().stream()));
//        } else {
//            props.setProperty("table.exclude.list", tableFormat(getExcludedTables(config.getExcludeTables()).stream()));
//        }
        try {
            Path p = Paths.get(secretFile);
            if (p.toFile().exists()) {
                byte[] data = Files.readAllBytes(Paths.get(secretFile));
                Map<String, String> secret = JSON.parseObject(data, Map.class);
//                props.setProperty("mongodb.user", secret.get("user"));
//                props.setProperty("mongodb.password",secret.get("password"));
//                props.setProperty("mongodb.authsource", "authsource");
            }
        } catch (IOException e) {
            LOGGER.error("read secret failed, error: {}", e.getMessage());
            throw e;
        }
        return props;
    }

    public Set<String> getSystemExcludedTables() {
        return new HashSet<>(Arrays.asList("information_schema", "mongodb", "performance_schema", "sys"));
    }

    public String tableFormat(Stream<String> table) {
        return table
                .map(stream -> this.getDatabase() + "." + stream)
                .collect(Collectors.joining(","));
    }

    public Set<String> getExcludedTables(Set<String> excludeTables) {
        Set<String> exclude = new HashSet<>(getSystemExcludedTables());
        exclude.addAll(excludeTables);
        return exclude;
    }
}

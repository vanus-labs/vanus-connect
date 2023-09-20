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

import com.fasterxml.jackson.annotation.JsonProperty;
import com.linkall.cdk.database.debezium.DebeziumConfig;
import io.debezium.connector.mongodb.MongoDbConnector;

import java.util.Arrays;
import java.util.Properties;
import java.util.stream.Collectors;

public class MongoDBConfig extends DebeziumConfig {
    private static final String DEBEZIUM_CLASS = "io.debezium.connector.mongodb.MongoDbConnector";

    @JsonProperty("name")
    private String name;

    @JsonProperty("connection_url")
    private String connectionUrl;

    @JsonProperty("hosts")
    private String[] hosts;

    @JsonProperty("credential")
    private MongoDBCredentials credentials;

    @JsonProperty("database_include")
    private String[] includeDatabases;
    @JsonProperty("database_exclude")
    private String[] excludeDatabases;

    @JsonProperty("collection_include")
    private String[] includeCollections;
    @JsonProperty("collection_exclude")
    private String[] excludeCollections;

    public MongoDBConfig() {
    }

    @Override
    public Class<?> secretClass() {
        return MongoDBCredentials.class;
    }

    @Override
    // https://debezium.io/documentation/reference/stable/connectors/mongodb.html#mongodb-connector-properties
    public Properties getDebeziumProperties() {
        final Properties props = new Properties();
        props.setProperty("connector.class", MongoDbConnector.class.getCanonicalName());
        props.setProperty("name", name);
        if (connectionUrl != null) {
            props.setProperty("mongodb.connection.string", connectionUrl);
        } else {
            props.setProperty("mongodb.hosts", String.join(",", this.hosts));
        }
        props.setProperty("capture.mode", "change_streams_update_full");
        if (credentials != null) {
            props.putAll(credentials.getProperties());
        }

        props.setProperty("topic.prefix", name);

        if (includeDatabases != null && includeDatabases.length > 0 &&
                excludeDatabases != null && excludeDatabases.length > 0) {
            throw new IllegalArgumentException("the database.include and database.exclude can't be set together");
        }
        if (includeDatabases != null && includeDatabases.length > 0) {
            props.setProperty("database.include.list", Arrays.stream(includeDatabases).collect(Collectors.joining(",")));
        }

        if (excludeDatabases != null && excludeDatabases.length > 0) {
            props.setProperty("database.exclude.list", Arrays.stream(excludeDatabases).collect(Collectors.joining(",")));
        }

        if (includeCollections != null && includeCollections.length > 0
                && excludeCollections != null && excludeCollections.length > 0) {
            throw new IllegalArgumentException("the collection.include and collection.exclude can't be set together");
        }
        if (includeCollections != null && includeCollections.length > 0) {
            props.setProperty("collection.include.list", Arrays.stream(includeCollections).collect(Collectors.joining(",")));
        }

        if (excludeCollections != null && excludeCollections.length > 0) {
            props.setProperty("collection.exclude.list", Arrays.stream(excludeCollections).collect(Collectors.joining(",")));
        }

        return props;
    }

    @Override
    protected Object getOffset() {
        return null;
    }
}

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

package com.linkall.connector.mongodb.debezium;

import java.util.HashSet;
import java.util.Set;

public class Config {
    private String host;
    private String port;
    private String username;
    private String password;
    private String database;
    private Set<String> includeTables;
    private Set<String> excludeTables;
    private String storeOffsetKey;

    public Config(
            String host,
            String port,
            String username,
            String password,
            String database,
            String includeTable,
            String excludeTable,
            String storeOffsetKey) {
        this.host = host;
        this.port = port;
        this.username = username;
        this.password = password;
        this.database = database;
        this.storeOffsetKey = storeOffsetKey;
        includeTables = new HashSet<>();
        excludeTables = new HashSet<>();
        if (includeTable != null) {
            for (String tableName : includeTable.split(",")) {
                includeTables.add(tableName);
            }
        }
        if (excludeTable != null) {
            for (String tableName : excludeTable.split(",")) {
                excludeTables.add(tableName);
            }
        }
    }

    public String getHost() {
        return host;
    }

    public void setHost(String host) {
        this.host = host;
    }

    public String getPort() {
        return port;
    }

    public void setPort(String port) {
        this.port = port;
    }

    public String getUsername() {
        return username;
    }

    public void setUsername(String username) {
        this.username = username;
    }

    public String getPassword() {
        return password;
    }

    public void setPassword(String password) {
        this.password = password;
    }

    public String getDatabase() {
        return database;
    }

    public void setDatabase(String database) {
        this.database = database;
    }

    public Set<String> getIncludeTables() {
        return includeTables;
    }

    public void setIncludeTables(Set<String> includeTables) {
        this.includeTables = includeTables;
    }

    public Set<String> getExcludeTables() {
        return excludeTables;
    }

    public void setExcludeTables(Set<String> excludeTables) {
        this.excludeTables = excludeTables;
    }

    public String getStoreOffsetKey() {
        return storeOffsetKey;
    }

    public void setStoreOffsetKey(String storeOffsetKey) {
        this.storeOffsetKey = storeOffsetKey;
    }
}

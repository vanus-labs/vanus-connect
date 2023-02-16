package com.linkall.sink.snowflake;

import com.fasterxml.jackson.annotation.JsonProperty;

import java.util.Properties;

public class DbConfig {
    @JsonProperty("host")
    private String host;
    @JsonProperty("username")
    private String username;
    @JsonProperty("password")
    private String password;
    @JsonProperty("role")
    private String role;
    @JsonProperty("warehouse")
    private String warehouse;
    @JsonProperty("database")
    private String database;
    @JsonProperty("schema")
    private String schema;
    @JsonProperty("table")
    private String table;
    @JsonProperty("properties")
    private Properties properties;

    public String getHost() {
        return host;
    }

    public String getUsername() {
        return username;
    }

    public String getPassword() {
        return password;
    }

    public String getRole() {
        return role;
    }

    public String getWarehouse() {
        return warehouse;
    }

    public String getDatabase() {
        return database;
    }

    public String getSchema() {
        return schema;
    }

    public String getTable() {
        return table;
    }

    public Properties getProperties() {
        return properties;
    }
}

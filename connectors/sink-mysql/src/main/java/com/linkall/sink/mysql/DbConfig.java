package com.linkall.sink.mysql;

import com.fasterxml.jackson.annotation.JsonProperty;

import java.util.Properties;

public class DbConfig {
    @JsonProperty("host")
    private String host;
    @JsonProperty("port")
    private int port;
    @JsonProperty("username")
    private String username;
    @JsonProperty("password")
    private String password;
    @JsonProperty("database")
    private String database;

    @JsonProperty("table_name")
    private String tableName;

    @JsonProperty("properties")
    private Properties properties;

    public String getHost() {
        return host;
    }

    public int getPort() {
        return port;
    }

    public String getUsername() {
        return username;
    }

    public String getPassword() {
        return password;
    }

    public String getDatabase() {
        return database;
    }

    public String getTableName() {
        return tableName;
    }

    public Properties getProperties() {
        return properties;
    }
}

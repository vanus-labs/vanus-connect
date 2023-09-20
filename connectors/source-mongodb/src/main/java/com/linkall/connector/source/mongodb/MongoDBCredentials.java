package com.linkall.connector.source.mongodb;

import com.fasterxml.jackson.annotation.JsonProperty;

import java.util.Properties;

public class MongoDBCredentials {
    @JsonProperty("username")
    private String username;
    @JsonProperty("password")
    private String password;
    @JsonProperty("auth_source")
    private String authSource;

    public MongoDBCredentials() {
    }

    public Properties getProperties() {
        final Properties props = new Properties();
        props.setProperty("mongodb.user", this.username);
        props.setProperty("mongodb.password", this.password);
        props.setProperty("mongodb.authsource", this.authSource);
        return props;
    }
}

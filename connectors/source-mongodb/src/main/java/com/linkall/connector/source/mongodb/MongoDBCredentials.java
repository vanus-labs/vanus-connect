package com.linkall.connector.source.mongodb;

import java.util.Properties;

public class MongoDBCredentials {
    private String username ;
    private String password ;
    private String authSource ;

    public MongoDBCredentials() {}

    public Properties getProperties() {
        final Properties props = new Properties();
        props.setProperty("mongodb.user", this.username);
        props.setProperty("mongodb.password", this.password);
        props.setProperty("mongodb.authsource", this.authSource);
        return props;
    }

    public void setUsername(String username) {
        this.username = username;
    }

    public void setPassword(String password) {
        this.password = password;
    }

    public void setAuthSource(String authSource) {
        this.authSource = authSource;
    }
}

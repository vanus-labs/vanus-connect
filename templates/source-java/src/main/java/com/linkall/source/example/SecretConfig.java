package com.linkall.source.example;

import com.fasterxml.jackson.annotation.JsonProperty;

public class SecretConfig {
    @JsonProperty("username")
    private String username;
    @JsonProperty("password")
    private String password;

    public String getUsername() {
        return username;
    }

    public String getPassword() {
        return password;
    }
}

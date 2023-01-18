package com.linkall.source.github;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.linkall.cdk.config.SourceConfig;

public class GitHubConfig extends SourceConfig {
    @JsonProperty("port")
    private int port;

    @JsonProperty("secret")
    private SecretConfig secretConfig;

    @Override
    public Class<?> secretClass() {
        return SecretConfig.class;
    }

    public int getPort() {
        return port;
    }

    public SecretConfig getSecretConfig() {
        return secretConfig;
    }
}

package com.linkall.source.example;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.linkall.cdk.config.SourceConfig;

public class ExampleConfig extends SourceConfig {
    @JsonProperty("secret")
    private SecretConfig secret;

    @Override
    public Class secretClass() {
        // TODO
        return SecretConfig.class;
    }

    public SecretConfig getSecret() {
        return secret;
    }
}

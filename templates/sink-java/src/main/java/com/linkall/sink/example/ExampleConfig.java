package com.linkall.sink.example;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.linkall.cdk.config.SinkConfig;

public class ExampleConfig extends SinkConfig {
    @JsonProperty
    private int num;

    public int getNum() {
        return num;
    }

    @Override
    public Class<?> secretClass() {
        // TODO
        return super.secretClass();
    }
}

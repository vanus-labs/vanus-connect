package com.linkall.sink.example;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.linkall.cdk.config.SinkConfig;

public class CloudEventConfig extends SinkConfig {
    @JsonProperty
    private int num;

    @JsonProperty("target")
    private String target;


    public int getNum() {
        return num;
    }

    public String getTarget() {
        return target;
    }

    @Override
    public Class<?> secretClass() {
        // TODO
        return super.secretClass();
    }
}

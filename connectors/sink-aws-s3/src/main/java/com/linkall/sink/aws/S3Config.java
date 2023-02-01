package com.linkall.sink.aws;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.linkall.cdk.config.SinkConfig;

public class S3Config extends SinkConfig {
    @JsonProperty("aws")
    private SecretConfig secretConfig;

    @JsonProperty("region")
    private String region;

    @JsonProperty("bucket")
    private String bucket;

    @JsonProperty("flush_size")
    private Integer flushSize;

    @JsonProperty("scheduled_interval")
    private Integer scheduledInterval;

    @JsonProperty("time_interval")
    private TimeInterval timeInterval;

    @Override
    public Class<?> secretClass() {
        return SecretConfig.class;
    }

    public SecretConfig getSecretConfig() {
        return secretConfig;
    }

    public String getRegion() {
        return region;
    }

    public String getBucket() {
        return bucket;
    }

    public Integer getFlushSize() {
        return flushSize;
    }

    public Integer getScheduledInterval() {
        return scheduledInterval;
    }

    public TimeInterval getTimeInterval() {
        return timeInterval;
    }

    public enum TimeInterval {
        HOURLY, DAILY
    }
}

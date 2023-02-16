package com.linkall.sink.snowflake;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.linkall.cdk.config.SinkConfig;

public class SnowflakeConfig extends SinkConfig {
    @JsonProperty("snowflake")
    private DbConfig snowflake;
    @JsonProperty("flush_time")
    private long flushTime;
    @JsonProperty("flush_size_bytes")
    private long sizeBytes;

    @Override
    public Class<?> secretClass() {
        return DbConfig.class;
    }

    public DbConfig getSnowflake() {
        return snowflake;
    }


    public long getFlushTime() {
        return flushTime;
    }

    public long getSizeBytes() {
        return sizeBytes;
    }
}

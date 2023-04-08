package com.linkall.sink.mysql;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.linkall.cdk.config.SinkConfig;

public class JdbcConfig extends SinkConfig {

    @JsonProperty("db")
    private DbConfig dbConfig;
    @JsonProperty("insert_mode")
    private InsertMode insertMode;
    @JsonProperty("commit_interval")
    private Long commitInterval;
    @JsonProperty("commit_size")
    private Integer commitSize;

    @Override
    public Class<?> secretClass() {
        return DbConfig.class;
    }

    public DbConfig getDbConfig() {
        return dbConfig;
    }

    public InsertMode getInsertMode() {
        return insertMode;
    }

    public Long getCommitInterval() {
        return commitInterval;
    }

    public Integer getCommitSize() {
        return commitSize;
    }

    public enum InsertMode {
        INSERT, UPSERT
    }
}

package com.linkall.source.kafka;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.linkall.cdk.config.SourceConfig;

import java.util.List;

public class KafkaConfig extends SourceConfig {
    @JsonProperty("bootstrap_servers")
    private String bootstrapServers;

    @JsonProperty("group_id")
    private String groupId;

    @JsonProperty("topics")
    private List<String> topics;

    public String getBootstrapServers() {
        return bootstrapServers;
    }

    public List<String> getTopics() {
        return topics;
    }

    public String getGroupId() {
        return groupId;
    }
}

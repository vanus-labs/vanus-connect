package com.linkall.source.aws.sns;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.linkall.cdk.config.SourceConfig;

public class SnsConfig extends SourceConfig {
    @JsonProperty("aws")
    private SecretConfig secretConfig;

    @JsonProperty("sns_arn")
    private String snsArn;

    @JsonProperty("protocol")
    private String protocol;

    @JsonProperty("endpoint")
    private String endpoint;

    @JsonProperty("port")
    private Integer port;

    @Override
    public Class<?> secretClass() {
        return SecretConfig.class;
    }

    public SecretConfig getSecretConfig() {
        return secretConfig;
    }

    public String getSnsArn() {
        return snsArn;
    }

    public String getProtocol() {
        return protocol;
    }

    public String getEndpoint(){
        return endpoint;
    }
    public Integer getPort() {
        return port;
    }
}

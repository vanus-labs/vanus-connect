package com.linkall.source.aws.sqs;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.linkall.cdk.config.SourceConfig;

public class SqsConfig extends SourceConfig {
    @JsonProperty("aws")
    private SecretConfig secretConfig;

    @JsonProperty("sqs_arn")
    private String sqsArn;

    @Override
    public Class<?> secretClass() {
        return SecretConfig.class;
    }

    public SecretConfig getSecretConfig() {
        return secretConfig;
    }

    public String getSqsArn() {
        return sqsArn;
    }

}

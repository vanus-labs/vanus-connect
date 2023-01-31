package com.linkall.source.aws.sqs;

import com.fasterxml.jackson.annotation.JsonProperty;

public class SecretConfig {
    @JsonProperty("access_key_id")
    private String accessKeyID;

    @JsonProperty("secret_access_key")
    private String secretAccessKey;

    public String getAccessKeyID() {
        return accessKeyID;
    }

    public String getSecretAccessKey() {
        return secretAccessKey;
    }
}

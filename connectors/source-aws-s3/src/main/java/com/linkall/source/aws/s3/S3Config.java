package com.linkall.source.aws.s3;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.linkall.cdk.config.SourceConfig;

import java.util.List;

public class S3Config extends SourceConfig {
    @JsonProperty("aws")
    private SecretConfig secretConfig;

    @JsonProperty("s3_bucket_arn")
    private String bucketArn;

    @JsonProperty("s3_events")
    private List<String> s3Events;

    @JsonProperty("region")
    private String region;


    @JsonProperty("sqs_arn")
    private String sqsArn;

    @Override
    public Class<?> secretClass() {
        return SecretConfig.class;
    }

    public SecretConfig getSecretConfig() {
        return secretConfig;
    }

    public String getBucketArn() {
        return bucketArn;
    }

    public List<String> getS3Events() {
        return s3Events;
    }

    public String getRegion() {
        return region;
    }

    public String getSqsArn() {
        return sqsArn;
    }
}

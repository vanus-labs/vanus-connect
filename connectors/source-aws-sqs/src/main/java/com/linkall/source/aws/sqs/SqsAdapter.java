package com.linkall.source.aws.sqs;


import io.cloudevents.CloudEvent;
import io.cloudevents.core.builder.CloudEventBuilder;

import java.net.URI;
import java.nio.charset.StandardCharsets;
import java.time.OffsetDateTime;

public class SqsAdapter {

    public CloudEvent adapt(SqsContent sqsContent) {
        CloudEventBuilder template = CloudEventBuilder.v1();
        template.withId(sqsContent.getMsgId());
        URI uri = URI.create("cloud.aws.sqs." + sqsContent.getRegion() + "." + sqsContent.getQueueName());
        template.withSource(uri);
        template.withType("com.amazonaws.sqs.message");
        template.withTime(OffsetDateTime.now());
        template.withData("text/plain",sqsContent.getBody().getBytes(StandardCharsets.UTF_8));

        return template.build();
    }
}

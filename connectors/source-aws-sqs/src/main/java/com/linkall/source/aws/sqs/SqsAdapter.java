package com.linkall.source.aws.sqs;

import com.linkall.vance.core.Adapter1;

import io.cloudevents.CloudEvent;
import io.cloudevents.core.builder.CloudEventBuilder;

import java.net.URI;
import java.nio.charset.StandardCharsets;
import java.time.OffsetDateTime;

public class SqsAdapter implements Adapter1<SqsContent> {
    private static final CloudEventBuilder template = CloudEventBuilder.v1();
    @Override
    public CloudEvent adapt(SqsContent sqsContent) {

        template.withId(sqsContent.getMsgId());
        URI uri = URI.create("cloud.aws.sqs."+sqsContent.getRegion()+"."+sqsContent.getQueueName());
        template.withSource(uri);
        template.withType("com.amazonaws.sqs.message");
        template.withDataContentType("text/plain");
        template.withTime(OffsetDateTime.now());
        template.withData(sqsContent.getBody().getBytes(StandardCharsets.UTF_8));

        return template.build();
    }
}

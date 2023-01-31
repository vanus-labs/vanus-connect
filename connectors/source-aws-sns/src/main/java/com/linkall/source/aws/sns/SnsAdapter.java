package com.linkall.source.aws.sns;

import io.cloudevents.CloudEvent;
import io.cloudevents.core.builder.CloudEventBuilder;
import io.vertx.core.buffer.Buffer;
import io.vertx.core.http.HttpServerRequest;
import io.vertx.core.json.JsonObject;
import org.apache.commons.lang.StringUtils;

import java.net.URI;
import java.time.OffsetDateTime;

public class SnsAdapter {

    public CloudEvent adapt(HttpServerRequest httpServerRequest, Buffer buffer) {
        CloudEventBuilder template = CloudEventBuilder.v1();
        JsonObject jsonObject = buffer.toJsonObject();
        template.withId(httpServerRequest.getHeader("X-Amz-Sns-Message-Id"));
        String subscriptionArn = httpServerRequest.getHeader("X-Amz-Sns-Subscription-Arn");
        URI uri = null;
        if (StringUtils.isBlank(subscriptionArn)) {
            uri = URI.create(httpServerRequest.getHeader("X-Amz-Sns-Topic-Arn"));
        } else {
            uri = URI.create(subscriptionArn);
        }
        template.withSource(uri);
        String type = "com.amazonaws.sns." + jsonObject.getString("Type");
        template.withType(type)
                .withDataContentType("application/json");
        String subject = jsonObject.getString("Subject");
        if (!StringUtils.isBlank(subject)) {
            template.withSubject(subject);
        }
        String timeStamp = jsonObject.getString("Timestamp");
        OffsetDateTime time = OffsetDateTime.parse(timeStamp);
        template.withTime(time);
        template.withData(buffer.getBytes());
        return template.build();
    }
}

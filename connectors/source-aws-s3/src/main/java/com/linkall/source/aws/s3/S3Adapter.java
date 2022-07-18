package com.linkall.source.aws.s3;

import com.linkall.vance.core.Adapter1;
import io.cloudevents.CloudEvent;
import io.cloudevents.core.builder.CloudEventBuilder;
import io.cloudevents.types.Time;
import io.vertx.core.json.JsonObject;

import java.net.URI;

public class S3Adapter implements Adapter1<JsonObject> {
    private static final CloudEventBuilder template = CloudEventBuilder.v1();
    @Override
    public CloudEvent adapt(JsonObject record) {
        JsonObject responseElements = record.getJsonObject("responseElements");
        template.withId(responseElements.getString("x-amz-request-id")+"."+responseElements.getString("x-amz-id-2"));
        URI uri = URI.create(record.getString("eventSource")+"."+record.getString("awsRegion")+"."+
                record.getJsonObject("s3").getJsonObject("bucket").getString("name"));
        template.withSource(uri);
        template.withType("com.amazonaws.s3."+record.getString("eventName"));
        template.withDataContentType("application/json");
        String timeStr = record.getString("eventTime");
        template.withTime(Time.parseTime("time", timeStr));
        template.withData(record.getJsonObject("s3").toBuffer().getBytes());
        template.withExtension("keyname",record.getJsonObject("s3").getJsonObject("object").getString("key"));
        return template.build();
    }
}

package com.linkall.source.kafka;

import com.linkall.vance.core.Adapter1;
import io.cloudevents.CloudEvent;
import io.cloudevents.core.builder.CloudEventBuilder;
import io.vertx.core.json.JsonObject;


import java.net.URI;
import java.time.OffsetDateTime;
import java.util.UUID;

public class KafkaAdapter implements Adapter1<KafkaData> {
    private static final CloudEventBuilder template = CloudEventBuilder.v1();

    public CloudEvent adapt(KafkaData kafkaData) {
        template.withId(UUID.randomUUID().toString());
        URI uri = URI.create("vance-kafka-source");
        template.withSource(uri);
        template.withType("kafka.message");
        template.withSubject(kafkaData.key());
        template.withTime(kafkaData.timestamp());
        try{
            new JsonObject(new String(kafkaData.value()));
            template.withDataContentType("application/json");
        }catch(Exception e){
            template.withDataContentType("plain/text");
        }
        template.withData(kafkaData.value());
        return template.build();

    }

}


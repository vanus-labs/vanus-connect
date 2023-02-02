package com.linkall.source.kafka;

import io.cloudevents.CloudEvent;
import io.cloudevents.core.builder.CloudEventBuilder;
import io.vertx.core.json.JsonObject;


import java.net.URI;
import java.util.UUID;

public class KafkaAdapter {

    public CloudEvent adapt(KafkaData kafkaData) {
        CloudEventBuilder template = CloudEventBuilder.v1();
        template.withId(UUID.randomUUID().toString());
        URI uri = URI.create("kafka." + kafkaData.KAFKA_SERVER_URL() +"."+ kafkaData.topic());
        template.withSource(uri);
        template.withType("kafka.message");
        template.withTime(kafkaData.timeStamp());
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


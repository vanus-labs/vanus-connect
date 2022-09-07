package com.linkall.source.kafka;

import java.time.OffsetDateTime;

public record KafkaData(String topic, String key, byte[] value, String KAFKA_SERVER_URL, OffsetDateTime timeStamp) {


}

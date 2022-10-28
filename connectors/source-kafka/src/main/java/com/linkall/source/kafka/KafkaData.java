package com.linkall.source.kafka;

import java.time.OffsetDateTime;

public record KafkaData(String topic, byte[] key, byte[] value, String KAFKA_SERVER_URL, OffsetDateTime timeStamp) {


}

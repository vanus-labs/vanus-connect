package com.linkall.source.kafka;

public record KafkaData(String topic, String key, byte[] value) {


}

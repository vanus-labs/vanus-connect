package com.linkall.source.kafka;

import com.linkall.vance.core.VanceApplication;

public class Entrance {
    public static void main(String[] args) {
        VanceApplication.run(KafkaSource.class);
    }

}
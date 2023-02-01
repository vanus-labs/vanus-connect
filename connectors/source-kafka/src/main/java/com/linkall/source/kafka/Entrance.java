package com.linkall.source.kafka;


import com.linkall.cdk.Application;

public class Entrance {
    public static void main(String[] args) {
        Application.run(KafkaSource.class);
    }
}
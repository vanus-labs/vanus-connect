package com.linkall.sink.aws;

import com.linkall.cdk.runtime.sender.HTTPSender;
import com.linkall.cdk.runtime.sender.Sender;
import io.cloudevents.CloudEvent;
import io.cloudevents.core.builder.CloudEventBuilder;

import java.net.URI;
import java.nio.charset.StandardCharsets;
import java.util.Arrays;
import java.util.UUID;
import java.util.concurrent.Executors;
import java.util.concurrent.ScheduledExecutorService;
import java.util.concurrent.TimeUnit;

public class SendFile {
    public static void main(String[] args) {

        CloudEventBuilder eventTemplate = CloudEventBuilder.v1()
                .withSource(URI.create("https://github.com/vanus-labs/vanus/connectors/sendfile"))
                .withDataContentType("application/json")
                .withType("simulation client");
        CloudEvent event = eventTemplate
                .withId(UUID.randomUUID().toString())
                .withExtension("objectkey","abcd.txt")
                .withData("application/json", "{\"key\": \"value\"}".getBytes(StandardCharsets.UTF_8))
                .build();

        Sender sender = new HTTPSender("http://localhost:8080");
        ScheduledExecutorService threadPool = Executors.newScheduledThreadPool(100);
        threadPool.scheduleAtFixedRate(new Runnable() {
            @Override
            public void run() {
                try {
                    sender.sendEvents(Arrays.asList(event));
                } catch (Throwable e) {
                    e.printStackTrace();
                }
            }
        }, 5L, 1L, TimeUnit.SECONDS);

    }
}

package com.linkall.sink.aws;

import com.linkall.vance.common.file.GenericFileUtil;
import com.linkall.vance.core.http.HttpClient;
import io.cloudevents.CloudEvent;
import io.cloudevents.core.builder.CloudEventBuilder;

import java.io.File;
import java.net.URI;
import java.nio.charset.StandardCharsets;
import java.time.Duration;
import java.time.LocalDateTime;
import java.util.UUID;
import java.util.concurrent.Executors;
import java.util.concurrent.ScheduledExecutorService;
import java.util.concurrent.TimeUnit;

public class SendFile {
    public static void main(String[] args) {

        CloudEventBuilder eventTemplate = CloudEventBuilder.v1()
                .withSource(URI.create("https://github.com/linkall-labs/vance/connectors/sendfile"))
                .withDataContentType("application/json")
                .withType("simulation client");
        CloudEvent event = eventTemplate
                .withId(UUID.randomUUID().toString())
                .withExtension("objectkey","abcd.txt")
                .withData("plain/text", GenericFileUtil.readResource("abcd.txt").getBytes(StandardCharsets.UTF_8))
                .build();

        ScheduledExecutorService threadPool = Executors.newScheduledThreadPool(100);
        threadPool.scheduleAtFixedRate(new Runnable() {
            @Override
            public void run() {
                HttpClient.deliver(event);
            }
        }, 15000L, 100L, TimeUnit.MILLISECONDS);

    }
}

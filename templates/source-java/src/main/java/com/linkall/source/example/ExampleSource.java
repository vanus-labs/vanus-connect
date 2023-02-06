package com.linkall.source.example;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.linkall.cdk.config.Config;
import com.linkall.cdk.connector.Element;
import com.linkall.cdk.connector.Source;
import com.linkall.cdk.connector.Tuple;
import com.linkall.cdk.util.EventUtil;
import io.cloudevents.CloudEvent;
import io.cloudevents.core.builder.CloudEventBuilder;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.net.URI;
import java.time.OffsetDateTime;
import java.util.HashMap;
import java.util.Map;
import java.util.UUID;
import java.util.concurrent.*;

public class ExampleSource implements Source {
    private static final Logger LOGGER = LoggerFactory.getLogger(ExampleSource.class);
    private final ObjectMapper objectMapper;

    private final BlockingQueue<Tuple> queue;
    private final ScheduledExecutorService executor;
    private ExampleConfig config;
    private int num;

    public ExampleSource() {
        objectMapper = new ObjectMapper();
        queue = new ArrayBlockingQueue<>(100);
        executor = Executors.newSingleThreadScheduledExecutor();
    }

    public CloudEvent makeEvent(int i) throws JsonProcessingException {
        CloudEventBuilder template = CloudEventBuilder.v1();
        template.withId(UUID.randomUUID().toString());
        URI uri = URI.create("testSource");
        template.withSource(uri);
        template.withType("testType");
        template.withDataContentType("application/json");
        template.withTime(OffsetDateTime.now());
        Map<String, Object> data = new HashMap<>();
        data.put("number", i);
        data.put("string", "Event Num " + i);
        template.withData("application/json", objectMapper.writeValueAsBytes(data));
        return template.build();
    }

    public void start() {
        executor.scheduleAtFixedRate(() -> {
            try {
                CloudEvent event = makeEvent(num++);
                queue.put(new Tuple(new Element(event, null), () ->
                        LOGGER.info("send event success {}", EventUtil.eventToJson(event))
                        , (success, failed, msg) -> LOGGER.info("send event failed:{}, {}", msg, EventUtil.eventToJson(event))));
            } catch (Exception e) {
                LOGGER.error("error", e);
            }
        }, 3, 10, TimeUnit.SECONDS);
    }

    @Override
    public Class<? extends Config> configClass() {
        // TODO
        return ExampleConfig.class;
    }

    @Override
    public void initialize(Config config) {
        // TODO
        this.config = (ExampleConfig) config;
        start();
    }

    @Override
    public String name() {
        // TODO
        return "ExampleSource";
    }

    @Override
    public void destroy() {
        // TODO
        executor.shutdown();
    }

    @Override
    public BlockingQueue<Tuple> queue() {
        // TODO
        return queue;
    }
}


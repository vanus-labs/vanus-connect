package com.linkall.sink.example;

import com.linkall.cdk.config.Config;
import com.linkall.cdk.connector.Element;
import com.linkall.cdk.connector.Result;
import com.linkall.cdk.connector.Sink;
import com.linkall.cdk.connector.Tuple;
import com.linkall.cdk.util.EventUtil;
import io.cloudevents.CloudEvent;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;


import java.util.concurrent.*;

public class CloudEventSink implements Sink {

    private static final Logger LOGGER = LoggerFactory.getLogger(CloudEventSink.class);
    private CloudEventConfig config;
    private final BlockingQueue<Tuple> queue;
    private final ScheduledExecutorService executor;

    public CloudEventSink() {
        queue = new ArrayBlockingQueue<>(100);
        executor = Executors.newSingleThreadScheduledExecutor();
    }


    @Override
    public Class<? extends Config> configClass() {
        // TODO
        return CloudEventConfig.class;
    }

    @Override
    public void initialize(Config config)  {
        this.config = (CloudEventConfig) config;
    }

    @Override
    public String name() {
        // TODO
        return "CloudEventSink";
    }

    @Override
    public void destroy() {
        // TODO
    }

    @Override
    public Result Arrived(CloudEvent... events) {
        for (CloudEvent event : events) {

            System.out.println(EventUtil.eventToJson(event));

        }
        return null;
    }



}
    


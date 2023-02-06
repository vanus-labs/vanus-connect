package com.linkall.sink.example;

import com.linkall.cdk.config.Config;
import com.linkall.cdk.connector.Result;
import com.linkall.cdk.connector.Sink;
import com.linkall.cdk.util.EventUtil;
import io.cloudevents.CloudEvent;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.concurrent.atomic.AtomicInteger;

public class ExampleSink implements Sink {

    private static final Logger LOGGER = LoggerFactory.getLogger(ExampleSink.class);
    private AtomicInteger eventNum;

    @Override
    public Class<? extends Config> configClass() {
        // TODO
        return ExampleConfig.class;
    }

    @Override
    public void initialize(Config cfg) throws Exception {
        // TODO
        ExampleConfig config = (ExampleConfig) cfg;
        eventNum = new AtomicInteger(config.getNum());
    }

    @Override
    public String name() {
        // TODO
        return "ExampleSink";
    }

    @Override
    public void destroy() {
        // TODO
    }

    @Override
    public Result Arrived(CloudEvent... events) {
        // TODO
        for (CloudEvent event : events) {
            int num = eventNum.addAndGet(1);
            // print number of received events
            LOGGER.info("receive a new event, in total: " + num);
            LOGGER.info(EventUtil.eventToJson(event));
        }
        return null;
    }
}

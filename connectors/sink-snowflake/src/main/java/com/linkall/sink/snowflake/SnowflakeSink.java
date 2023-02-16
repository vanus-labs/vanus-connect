package com.linkall.sink.snowflake;

import com.linkall.cdk.config.Config;
import com.linkall.cdk.connector.Result;
import com.linkall.cdk.connector.Sink;
import io.cloudevents.CloudEvent;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class SnowflakeSink implements Sink {

    private static final Logger LOGGER = LoggerFactory.getLogger(SnowflakeSink.class);

    private SnowflakeService service;

    @Override
    public Result Arrived(CloudEvent... events) {
        for (CloudEvent event : events) {
            try {
                service.addData(event);
            } catch (Exception e) {
                LOGGER.error("writer event failed", e);
                return new Result(500, e.getMessage());
            }
        }
        return Result.SUCCESS;
    }

    @Override
    public Class<? extends Config> configClass() {
        return SnowflakeConfig.class;
    }

    @Override
    public void initialize(Config cfg) throws Exception {
        SnowflakeConfig config = (SnowflakeConfig) cfg;
        this.service = new SnowflakeService(config);
        service.start();
    }

    @Override
    public String name() {
        return "Snowflake Sink";
    }

    @Override
    public void destroy() throws Exception {
        service.stop();
    }
}

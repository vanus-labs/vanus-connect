package com.linkall.sink.display;

import com.linkall.vance.common.json.JsonMapper;
import com.linkall.vance.core.Sink;
import com.linkall.vance.core.http.HttpServer;
import io.vertx.core.json.JsonObject;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.concurrent.atomic.AtomicInteger;

public class DisplaySink implements Sink {
    private static final Logger LOGGER = LoggerFactory.getLogger(DisplaySink.class);
    private static final AtomicInteger eventNum = new AtomicInteger(0);

    @Override
    public void start(){
        HttpServer server = HttpServer.createHttpServer();
        server.ceHandler(event -> {
            int num = eventNum.addAndGet(1);
            LOGGER.info("receive a new event, in total: "+num);
            JsonObject js = JsonMapper.wrapCloudEvent(event);
            LOGGER.info(js.encodePrettily());
        });
        server.listen();
    }

}
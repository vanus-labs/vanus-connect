package com.linkall.source.http;

import com.linkall.vance.core.Adapter;
import com.linkall.vance.core.Adapter2;
import com.linkall.vance.core.Source;
import com.linkall.vance.core.http.HttpClient;
import com.linkall.vance.core.http.HttpServer;
import io.vertx.core.buffer.Buffer;
import io.vertx.core.http.HttpServerRequest;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.concurrent.atomic.AtomicInteger;


public class HttpSource implements Source {
    private static final Logger LOGGER = LoggerFactory.getLogger(HttpSource.class);
    private static final AtomicInteger eventNum = new AtomicInteger(0);

    @SuppressWarnings("unchecked")
    @Override
    public void start(){
        HttpServer server = HttpServer.createHttpServer();
        HttpClient.setDeliverSuccessHandler(resp->{
            int num = eventNum.addAndGet(1);
            LOGGER.info("send event in total: "+num);
        });
        server.simpleHandler(((Adapter2<HttpServerRequest,Buffer>) getAdapter()));
        server.listen();
    }

    @Override
    public Adapter getAdapter() {
        return new HttpAdapter();
    }
}
package com.linkall.source.http;

import com.linkall.core.Adapter;
import com.linkall.core.Adapter2;
import com.linkall.core.Source;
import com.linkall.core.http.HttpClient;
import com.linkall.core.http.HttpServer;
import io.cloudevents.CloudEvent;
import io.cloudevents.core.builder.CloudEventBuilder;
import io.vertx.core.buffer.Buffer;
import io.vertx.core.http.HttpServerRequest;
import io.vertx.core.json.JsonObject;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.net.URI;
import java.time.OffsetDateTime;
import java.util.UUID;
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
class HttpAdapter implements Adapter2<HttpServerRequest,Buffer> {
    private static final CloudEventBuilder template = CloudEventBuilder.v1();
    @Override
    public CloudEvent adapt(HttpServerRequest req, Buffer buffer) {
        template.withId(UUID.randomUUID().toString());
        URI uri = URI.create("vance-http-source");
        template.withSource(uri);
        template.withType("http");
        template.withDataContentType("application/json");
        template.withTime(OffsetDateTime.now());

        JsonObject data = new JsonObject();
        JsonObject headers = new JsonObject();
        req.headers().forEach((m)-> headers.put(m.getKey(),m.getValue()));
        data.put("headers",headers);
        String contentType = req.getHeader("content-type");
        if(null != contentType && contentType.equals("application/json")){
            JsonObject body = buffer.toJsonObject();
            data.put("body",body);
        }else{
            String myData = new String(buffer.getBytes());
            JsonObject body = new JsonObject();
            body.put("data",myData);
            data.put("body",body);
        }
        template.withData(data.toBuffer().getBytes());

        return template.build();
    }
}

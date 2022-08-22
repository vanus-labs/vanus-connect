package com.linkall.connector.mongodb.debezium;

import com.linkall.vance.core.Adapter1;
import io.cloudevents.CloudEvent;
import io.debezium.engine.ChangeEvent;
import io.debezium.engine.DebeziumEngine;
import io.vertx.core.Vertx;
import io.vertx.ext.web.client.WebClient;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.List;
import java.util.concurrent.CountDownLatch;

public class RecordConsumer
        implements DebeziumEngine.ChangeConsumer<ChangeEvent<String, String>> {
    private static final Logger LOGGER = LoggerFactory.getLogger(RecordConsumer.class);
    private final Adapter1<String> adapter;
    private final WebClient webClient;

    public RecordConsumer(Adapter1<String> adapter) {
        this.adapter = adapter;
        webClient = WebClient.create(Vertx.vertx());
    }

    @Override
    public void handleBatch(List<ChangeEvent<String, String>> records, DebeziumEngine.RecordCommitter<ChangeEvent<String, String>> committer) throws InterruptedException {
        CountDownLatch latch = new CountDownLatch(records.size());
        for (ChangeEvent<String, String> record : records) {
            LOGGER.debug("Received event '{}'", record);
            if (record.value() == null) {
                latch.countDown();
                continue;
            }
            // TODO add to dead letter
            CloudEvent ceEvent = this.adapter.adapt(record.value());
//            Future<HttpResponse<Buffer>> responseFuture =
//                    VertxMessageFactory.createWriter(webClient.postAbs(ConfigUtil.getVanceSink()))
//                            .writeStructured(ceEvent, JsonFormat.CONTENT_TYPE);
//            responseFuture.onComplete(
//                    ar -> {
//                        if (ar.failed()) {
//                            LOGGER.warn("Error to send record: {},error: {}", record, ar.cause());
//                        } else if (ar.result().statusCode() == HTTP_OK
//                                || ar.result().statusCode() == HTTP_NO_CONTENT
//                                || ar.result().statusCode() == HTTP_ACCEPTED) {
//                            LOGGER.debug("Success to send cloudEvent：{}", ceEvent.getId());
//                        } else {
//                            LOGGER.warn(
//                                    "Failed to send record: {},statusCode: {}, body: {}",
//                                    record,
//                                    ar.result().statusCode(),
//                                    ar.result().bodyAsString());
//                        }
//                        try {
//                            committer.markProcessed(record);
//                        } catch (InterruptedException e) {
//                            LOGGER.warn(
//                                    "Failed to mark processed record: {},error: {}",
//                                    record.value().sourceOffset(),
//                                    e);
//                        }
//                        latch.countDown();
//                    });
        }
        latch.await();
        committer.markBatchFinished();
    }
}

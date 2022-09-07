package com.linkall.source.kafka;

import com.linkall.vance.core.Adapter;
import com.linkall.vance.core.Source;
import io.vertx.core.Vertx;
import io.vertx.ext.web.client.WebClient;

import java.util.concurrent.ConcurrentLinkedQueue;
import java.util.concurrent.Executors;
import java.util.concurrent.ScheduledExecutorService;
import java.util.concurrent.TimeUnit;


public class KafkaSource implements Source  {

    public final static Vertx vertx = Vertx.vertx();
    public final static WebClient webClient = WebClient.create(vertx);

    public void start(){
        KafkaWorker worker = new KafkaWorker("kafkawork",false);
        worker.start();
    }

    public Adapter getAdapter() {
        return new KafkaAdapter();
    }

}

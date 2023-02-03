package com.linkall.source.kafka;

import com.linkall.cdk.config.Config;
import com.linkall.cdk.connector.Source;
import com.linkall.cdk.connector.Tuple;
import io.vertx.core.Vertx;
import io.vertx.ext.web.client.WebClient;

import java.util.concurrent.BlockingQueue;
import java.util.concurrent.LinkedBlockingQueue;


public class KafkaSource implements Source {

    public final static Vertx vertx = Vertx.vertx();
    public final static WebClient webClient = WebClient.create(vertx);

    private BlockingQueue<Tuple> queue;
    private KafkaConfig config;

    public KafkaSource() {
        queue = new LinkedBlockingQueue<>(100);
    }

    public void start() {
        KafkaWorker worker = new KafkaWorker("kafkawork", false, config, queue);
        worker.start();
    }

    @Override
    public BlockingQueue<Tuple> queue() {
        return queue;
    }

    @Override
    public Class<? extends Config> configClass() {
        return KafkaConfig.class;
    }

    @Override
    public void initialize(Config config) throws Exception {
        this.config = (KafkaConfig) config;
        start();
    }

    @Override
    public String name() {
        return "KafkaSource";
    }

    @Override
    public void destroy() throws Exception {

    }
}

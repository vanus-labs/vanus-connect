package com.linkall.source.kafka;

import com.linkall.vance.common.config.ConfigUtil;
import io.cloudevents.CloudEvent;
import io.cloudevents.http.vertx.VertxMessageFactory;
import io.cloudevents.jackson.JsonFormat;
import io.vertx.circuitbreaker.CircuitBreaker;
import io.vertx.circuitbreaker.CircuitBreakerOptions;
import io.vertx.core.Future;
import io.vertx.core.buffer.Buffer;
import io.vertx.core.impl.logging.Logger;
import io.vertx.core.impl.logging.LoggerFactory;
import io.vertx.ext.web.client.HttpResponse;
import kafka.utils.ShutdownableThread;
import org.apache.kafka.clients.consumer.*;
import org.apache.kafka.common.TopicPartition;


import java.time.*;
import java.util.*;
import java.util.concurrent.*;
import java.util.stream.Collectors;


public class KafkaWorker extends ShutdownableThread {
    private static final Logger LOGGER = LoggerFactory.getLogger(KafkaWorker.class);
    private final KafkaConsumer<byte[], byte[]> consumer;
    private final String topicList;
    private final KafkaAdapter adapter;
    private final String KAFKA_SERVER_URL;

    public final ConcurrentHashMap<TopicPartition,Long> offsets = new ConcurrentHashMap<>();


    public KafkaWorker(String name, boolean isInterruptible)  {
        super(name, isInterruptible);

        KAFKA_SERVER_URL = ConfigUtil.getString("KAFKA_SERVER_URL");
        String KAFKA_SERVER_PORT = ConfigUtil.getString("KAFKA_SERVER_PORT");
        String CLIENT_ID = ConfigUtil.getString("CLIENT_ID");
        topicList = ConfigUtil.getString("TOPIC_LIST");

        Properties properties = new Properties();
        properties.put(ConsumerConfig.BOOTSTRAP_SERVERS_CONFIG, KAFKA_SERVER_URL + ":" + KAFKA_SERVER_PORT);
        properties.put(ConsumerConfig.GROUP_ID_CONFIG, CLIENT_ID);
        properties.put(ConsumerConfig.ENABLE_AUTO_COMMIT_CONFIG, "false");
        properties.put(ConsumerConfig.AUTO_OFFSET_RESET_CONFIG, "earliest");
        properties.put(ConsumerConfig.AUTO_COMMIT_INTERVAL_MS_CONFIG, "1000");
        properties.put(ConsumerConfig.SESSION_TIMEOUT_MS_CONFIG, "30000");
        properties.put(ConsumerConfig.KEY_DESERIALIZER_CLASS_CONFIG, "org.apache.kafka.common.serialization.ByteArrayDeserializer");
        properties.put(ConsumerConfig.VALUE_DESERIALIZER_CLASS_CONFIG, "org.apache.kafka.common.serialization.ByteArrayDeserializer");

        consumer = new KafkaConsumer<>(properties);
        String[] topicsListArray = null;
        boolean haveSpace = topicList.contains(", ");
        if (haveSpace) {
            topicsListArray = topicList.split(", ");
        } else
            topicsListArray = topicList.split(",");

       
        consumer.subscribe(Arrays.asList(topicsListArray));

        adapter = new KafkaAdapter();
    }


    @Override
    public void doWork() {
        ConsumerRecords<byte[], byte[]> records = consumer.poll(Duration.ofSeconds(25));
        System.out.println("records.partitions() size: "+records.partitions().size());
        for (TopicPartition partition : records.partitions()) {
            List<ConsumerRecord<byte[], byte[]>> partitionRecords = records.records(partition);
            System.out.println("partitionRecords size: "+partitionRecords.size());
            ConcurrentHashMap<Long,Boolean> cm = new ConcurrentHashMap<>();

            for (ConsumerRecord<byte[], byte[]> record : partitionRecords) {
                OffsetDateTime timeStamp = new Date(record.timestamp()).toInstant().atOffset( ZoneOffset.UTC );
                KafkaData kafkaData = new KafkaData(record.topic(), record.key(), record.value(), KAFKA_SERVER_URL, timeStamp);
                CloudEvent event = adapter.adapt(kafkaData);
                String sink = ConfigUtil.getVanceSink();
                System.out.println("message: " + Arrays.toString(record.value()));
                System.out.println("offset: " + record.offset());


                CircuitBreaker breaker = CircuitBreaker.create("my-circuit-breaker", KafkaSource.vertx,
                        new CircuitBreakerOptions()
                                .setMaxRetries(5)
                                .setTimeout(30000)
                );
                breaker.<String>execute(promise -> {
                    LOGGER.info("try to send request");

            Future<HttpResponse<Buffer>> responseFuture = VertxMessageFactory.createWriter(KafkaSource.webClient.postAbs(sink))
                    .writeStructured(event, JsonFormat.CONTENT_TYPE);

                    responseFuture.onSuccess(resp-> {
                        promise.complete();
                        LOGGER.info("send task success");
                        cm.put(record.offset(),true);
                    });
                    responseFuture.onFailure(System.err::println);
                }).onFailure(t->{
                    LOGGER.info("send task failed");
                    cm.put(record.offset(),false);
                    System.out.println("The targeted server is down!");
                    System.exit(1);
                });

                breaker.close();
            }
            while(cm.size()!=partitionRecords.size()){
                try {
                    Thread.sleep(1000);
                } catch (InterruptedException e) {
                    throw new RuntimeException(e);
                }
            }

            TreeMap<Long,Boolean> tm = new TreeMap<>(cm);
            long lastOffset = partitionRecords.get(partitionRecords.size() - 1).offset()+1;
            Iterator<Map.Entry<Long,Boolean>> it = tm.entrySet().iterator();
            while (it.hasNext()){
                Map.Entry<Long,Boolean> entry = it.next();
                if(!entry.getValue()){
                    lastOffset = entry.getKey();
                    break;
                }
            }
            System.out.println("lastOffset: "+lastOffset);
            consumer.commitSync(Collections.singletonMap(partition, new OffsetAndMetadata(lastOffset)));

        }

    }
}

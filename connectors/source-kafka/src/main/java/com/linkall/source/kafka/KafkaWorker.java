package com.linkall.source.kafka;

import com.linkall.cdk.connector.Element;
import com.linkall.cdk.connector.Tuple;
import io.cloudevents.CloudEvent;
import kafka.utils.ShutdownableThread;
import org.apache.kafka.clients.consumer.*;
import org.apache.kafka.common.TopicPartition;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.time.Duration;
import java.time.OffsetDateTime;
import java.time.ZoneOffset;
import java.util.*;
import java.util.concurrent.BlockingQueue;
import java.util.concurrent.ConcurrentHashMap;


public class KafkaWorker extends ShutdownableThread {
    private static final Logger LOGGER = LoggerFactory.getLogger(KafkaWorker.class);
    private final KafkaConsumer<byte[], byte[]> consumer;
    private final KafkaAdapter adapter;

    public final ConcurrentHashMap<TopicPartition, Long> offsets = new ConcurrentHashMap<>();
    private BlockingQueue<Tuple> queue;
    private KafkaConfig config;

    public KafkaWorker(String name, boolean isInterruptible, KafkaConfig config, BlockingQueue<Tuple> queue) {
        super(name, isInterruptible);
        this.config = config;
        this.queue = queue;

        Properties properties = new Properties();
        properties.put(ConsumerConfig.BOOTSTRAP_SERVERS_CONFIG, config.getBootstrapServers());
        properties.put(ConsumerConfig.GROUP_ID_CONFIG, config.getGroupId());
        properties.put(ConsumerConfig.ENABLE_AUTO_COMMIT_CONFIG, "false");
        properties.put(ConsumerConfig.AUTO_OFFSET_RESET_CONFIG, "earliest");
        properties.put(ConsumerConfig.AUTO_COMMIT_INTERVAL_MS_CONFIG, "1000");
        properties.put(ConsumerConfig.SESSION_TIMEOUT_MS_CONFIG, "30000");
        properties.put(ConsumerConfig.KEY_DESERIALIZER_CLASS_CONFIG, "org.apache.kafka.common.serialization.ByteArrayDeserializer");
        properties.put(ConsumerConfig.VALUE_DESERIALIZER_CLASS_CONFIG, "org.apache.kafka.common.serialization.ByteArrayDeserializer");

        consumer = new KafkaConsumer<>(properties);

        consumer.subscribe(config.getTopics());

        adapter = new KafkaAdapter();
    }


    @Override
    public void doWork() {
        ConsumerRecords<byte[], byte[]> records = consumer.poll(Duration.ofSeconds(25));
        System.out.println("records.partitions() size: " + records.partitions().size());
        for (TopicPartition partition : records.partitions()) {
            List<ConsumerRecord<byte[], byte[]>> partitionRecords = records.records(partition);
            System.out.println("partitionRecords size: " + partitionRecords.size());
            ConcurrentHashMap<Long, Boolean> cm = new ConcurrentHashMap<>();

            for (ConsumerRecord<byte[], byte[]> record : partitionRecords) {
                OffsetDateTime timeStamp = new Date(record.timestamp()).toInstant().atOffset(ZoneOffset.UTC);
                KafkaData kafkaData = new KafkaData(record.topic(), record.key(), record.value(), config.getBootstrapServers(), timeStamp);
                CloudEvent event = adapter.adapt(kafkaData);
                Tuple tuple = new Tuple(new Element(event, record), () -> {
                    LOGGER.info("send event success,topic:{},offset:{},id:{}", record.topic(), record.offset(), event.getId());
                    cm.put(record.offset(), true);
                }, (success, failed, msg) -> {
                    LOGGER.info("send event failed,topic:{},offset:{},id:{},msg:{}", record.topic(), record.offset(), event.getId(), msg);
                    cm.put(record.offset(), false);
                    System.out.println("The targeted server is down!");
                    System.exit(1);
                });
                try {
                    queue.put(tuple);
                } catch (InterruptedException e) {
                    LOGGER.warn("put event interrupted");
                }
            }
            while (cm.size()!=partitionRecords.size()) {
                try {
                    Thread.sleep(1000);
                } catch (InterruptedException e) {
                    throw new RuntimeException(e);
                }
            }

            TreeMap<Long, Boolean> tm = new TreeMap<>(cm);
            long lastOffset = partitionRecords.get(partitionRecords.size() - 1).offset() + 1;
            Iterator<Map.Entry<Long, Boolean>> it = tm.entrySet().iterator();
            while (it.hasNext()) {
                Map.Entry<Long, Boolean> entry = it.next();
                if (!entry.getValue()) {
                    lastOffset = entry.getKey();
                    break;
                }
            }
            System.out.println("lastOffset: " + lastOffset);
            consumer.commitSync(Collections.singletonMap(partition, new OffsetAndMetadata(lastOffset)));

        }

    }
}

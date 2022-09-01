package com.linkall.source.kafka;

import com.linkall.vance.common.config.ConfigUtil;
import io.cloudevents.CloudEvent;
import io.cloudevents.http.vertx.VertxMessageFactory;
import io.cloudevents.jackson.JsonFormat;
import io.vertx.core.Future;
import io.vertx.core.buffer.Buffer;
import io.vertx.core.impl.logging.Logger;
import io.vertx.core.impl.logging.LoggerFactory;
import io.vertx.ext.web.client.HttpResponse;
import kafka.utils.ShutdownableThread;
import org.apache.kafka.clients.consumer.*;
import org.apache.kafka.common.TopicPartition;

import java.time.Duration;
import java.util.*;

public class KafkaWorker extends ShutdownableThread {
    private static final Logger LOGGER = LoggerFactory.getLogger(KafkaWorker.class);
    private final KafkaConsumer<byte[], byte[]> consumer;
    private final String topicList;
    private final KafkaAdapter adapter;


    public KafkaWorker(String name, boolean isInterruptible) {
        super(name, isInterruptible);

        String KAFKA_SERVER_URL = ConfigUtil.getString("KAFKA_SERVER_URL");
        String KAFKA_SERVER_PORT = ConfigUtil.getString("KAFKA_SERVER_PORT");
        String CLIENT_ID = ConfigUtil.getString("CLIENT_ID");
        topicList = ConfigUtil.getString("TOPIC_LIST");

        Properties properties = new Properties();
        properties.put(ConsumerConfig.BOOTSTRAP_SERVERS_CONFIG, KAFKA_SERVER_URL + ":" + KAFKA_SERVER_PORT);
        properties.put(ConsumerConfig.GROUP_ID_CONFIG, CLIENT_ID);
        properties.put(ConsumerConfig.ENABLE_AUTO_COMMIT_CONFIG, "true");
        properties.put(ConsumerConfig.AUTO_COMMIT_INTERVAL_MS_CONFIG, "1000");
        properties.put(ConsumerConfig.SESSION_TIMEOUT_MS_CONFIG, "30000");
        properties.put(ConsumerConfig.KEY_DESERIALIZER_CLASS_CONFIG, "org.apache.kafka.common.serialization.ByteArrayDeserializer");
        properties.put(ConsumerConfig.VALUE_DESERIALIZER_CLASS_CONFIG, "org.apache.kafka.common.serialization.ByteArrayDeserializer");

        consumer = new KafkaConsumer<>(properties);
        adapter = new KafkaAdapter();
    }


    @Override
    public void doWork() {
        String[] topicsListArray = topicList.split(", ");
        consumer.subscribe(Arrays.asList(topicsListArray));
        ConsumerRecords<byte[], byte[]> records = consumer.poll(Duration.ofSeconds(25));
        for (TopicPartition partition : records.partitions()) {
            List<ConsumerRecord<byte[], byte[]>> partitionRecords = records.records(partition);
            for (ConsumerRecord<byte[], byte[]> record : partitionRecords) {

                String key64 = Base64.getEncoder().encodeToString(record.key());
                KafkaData kafkaData = new KafkaData(record.topic(), key64, record.value());

                CloudEvent event = adapter.adapt(kafkaData);
                String sink = ConfigUtil.getVanceSink();

                Future<HttpResponse<Buffer>> responseFuture;
                responseFuture = VertxMessageFactory.createWriter(KafkaSource.webClient.postAbs(sink))
                        .writeStructured(event, JsonFormat.CONTENT_TYPE);

                responseFuture.onSuccess(resp-> {
                            LOGGER.info("send CloudEvent success");
                            long lastOffset = partitionRecords.get(partitionRecords.size() - 1).offset();
                            consumer.commitSync(Collections.singletonMap(partition, new OffsetAndMetadata(lastOffset + 1)));
                        })

                        .onFailure(t-> LOGGER.info("send task failed"));

            }

        }
    }
}

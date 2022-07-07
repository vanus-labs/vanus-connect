package com.linkall.source.mysql;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.google.common.base.Preconditions;
import com.linkall.vance.common.store.KVStoreFactory;
import com.linkall.vance.core.KVStore;
import org.apache.kafka.connect.runtime.WorkerConfig;
import org.apache.kafka.connect.storage.MemoryOffsetBackingStore;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.nio.ByteBuffer;
import java.nio.charset.StandardCharsets;
import java.util.Map;
import java.util.Optional;
import java.util.stream.Collectors;

public class KvStoreOffsetBackingStore extends MemoryOffsetBackingStore {
  public static final String OFFSET_STORAGE_KV_STORE_KEY_CONFIG = "offset.storage.kv.key";
  private static final Logger logger = LoggerFactory.getLogger(KvStoreOffsetBackingStore.class);
  private static final String DEFAULT_KEY_NAME = "vance_debezium_mysql_offset";
  private final ObjectMapper objectMapper = new ObjectMapper();
  private final KVStore store;
  private String keyName;

  public KvStoreOffsetBackingStore() {
    store = KVStoreFactory.createKVStore();
  }

  private static String byteBufferToString(final ByteBuffer byteBuffer) {
    Preconditions.checkNotNull(byteBuffer);
    return new String(byteBuffer.array(), StandardCharsets.UTF_8);
  }

  private static ByteBuffer stringToByteBuffer(final String s) {
    Preconditions.checkNotNull(s);
    return ByteBuffer.wrap(s.getBytes(StandardCharsets.UTF_8));
  }

  @Override
  public void configure(WorkerConfig config) {
    super.configure(config);
    Map<String, String> map = config.originalsStrings();
    keyName =
        Optional.ofNullable(map.get(OFFSET_STORAGE_KV_STORE_KEY_CONFIG)).orElse(DEFAULT_KEY_NAME);
  }

  @Override
  public synchronized void start() {
    super.start();
    logger.info("Starting KvStoreOffsetBackingStore with key {}", keyName);
    load();
  }

  @Override
  public synchronized void stop() {
    super.stop();
    logger.info("Stopped KvStoreOffsetBackingStore");
  }

  @SuppressWarnings("unchecked")
  private void load() {
    String value = store.get(keyName);
    if (value == null || value.length() == 0) {
      return;
    }
    logger.info("Load offset: {}", value);
    Map<String, String> mapAsString = null;
    try {
      mapAsString = objectMapper.readValue(value, Map.class);
    } catch (JsonProcessingException e) {
      logger.error("Fail to load offset: {}, error: {}", value, e);
      throw new RuntimeException(e);
    }
    data =
        mapAsString.entrySet().stream()
            .collect(
                Collectors.toMap(
                    e -> stringToByteBuffer(e.getKey()), e -> stringToByteBuffer(e.getValue())));
  }

  @Override
  protected void save() {
    Map<String, String> raw =
        data.entrySet().stream()
            .collect(
                Collectors.toMap(
                    e -> byteBufferToString(e.getKey()), e -> byteBufferToString(e.getValue())));
    try {
      String value = objectMapper.writeValueAsString(raw);
      store.put(keyName, value);
    } catch (JsonProcessingException e) {
      logger.warn("Fail to save offset: {}, error: {}", raw, e);
    }
  }
}

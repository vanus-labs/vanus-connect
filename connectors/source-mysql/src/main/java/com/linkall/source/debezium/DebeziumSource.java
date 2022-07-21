package com.linkall.source.debezium;

import com.linkall.vance.common.config.ConfigUtil;
import com.linkall.vance.common.config.SecretUtil;
import com.linkall.vance.core.Adapter1;
import com.linkall.vance.core.Source;
import io.debezium.embedded.Connect;
import io.debezium.engine.ChangeEvent;
import io.debezium.engine.DebeziumEngine;
import io.debezium.engine.spi.OffsetCommitPolicy;
import org.apache.kafka.connect.json.JsonConverter;
import org.apache.kafka.connect.json.JsonConverterConfig;
import org.apache.kafka.connect.source.SourceRecord;
import org.apache.kafka.connect.storage.Converter;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.nio.charset.StandardCharsets;
import java.util.*;
import java.util.concurrent.Executor;
import java.util.concurrent.Executors;
import java.util.stream.Collectors;
import java.util.stream.Stream;

public abstract class DebeziumSource implements Source {
  private static final Logger LOGGER = LoggerFactory.getLogger(DebeziumSource.class);
  private DebeziumEngine<ChangeEvent<SourceRecord, SourceRecord>> engine;
  private Executor executor;
  private final DebeziumEngine.ChangeConsumer<ChangeEvent<SourceRecord, SourceRecord>> consumer;
  protected final DbConfig config;

  public DebeziumSource() {
    consumer = new DebeziumRecordConsumer((Adapter1<SourceRecord>) getAdapter());
    config =
        new DbConfig(
            SecretUtil.getString("host"),
            SecretUtil.getString("port"),
            SecretUtil.getString("username"),
            SecretUtil.getString("password"),
            SecretUtil.getString("dbName"),
            ConfigUtil.getString("include_table"),
            ConfigUtil.getString("exclude_table"),
            ConfigUtil.getString("store_offset_key"));
  }

  public abstract String getConnectorClass();

  public abstract Map<String, Object> getConfigOffset();

  public abstract Properties getDebeziumProperties();

  public void start() throws Exception {
    engine =
        DebeziumEngine.create(Connect.class)
            .using(getProperties())
            .using(OffsetCommitPolicy.always())
            .notifying(consumer)
            .using(
                (success, message, error) -> {
                  LOGGER.info(
                      "Debezium engine shutdown,success: {},message: {},error:{}",
                      success,
                      message,
                      error);
                })
            .build();
    executor = Executors.newSingleThreadExecutor();
    executor.execute(engine);
  }

  private Properties getProperties() {
    final Properties props = new Properties();

    // debezium engine configuration
    props.setProperty("connector.class", getConnectorClass());
    // snapshot config
    props.setProperty("snapshot.mode", "initial");
    // DO NOT include schema change, e.g. DDL
    props.setProperty("include.schema.changes", "false");
    // disable tombstones
    props.setProperty("tombstones.on.delete", "false");

    // offset
    props.setProperty("offset.storage", KvStoreOffsetBackingStore.class.getCanonicalName());
    if (config.getStoreOffsetKey() != null && !config.getStoreOffsetKey().isEmpty()) {
      props.setProperty(
          KvStoreOffsetBackingStore.OFFSET_STORAGE_KV_STORE_KEY_CONFIG, config.getStoreOffsetKey());
    }
    Map<String, Object> configOffset = getConfigOffset();
    if (configOffset != null && configOffset.size() > 0) {
      Converter valueConverter = new JsonConverter();
      Map<String, Object> valueConfigs = new HashMap<>();
      valueConfigs.put(JsonConverterConfig.SCHEMAS_ENABLE_CONFIG, false);
      valueConverter.configure(valueConfigs, false);
      byte[] offsetValue = valueConverter.fromConnectData(config.getDatabase(), null, configOffset);
      props.setProperty(
          KvStoreOffsetBackingStore.OFFSET_CONFIG_VALUE,
          new String(offsetValue, StandardCharsets.UTF_8));
    }

    props.setProperty("offset.flush.interval.ms", "1000");

    // history
    props.setProperty("database.history", "io.debezium.relational.history.FileDatabaseHistory");
    props.setProperty("database.history.file.filename", "/tmp/mysql/history.data");

    // https://debezium.io/documentation/reference/configuration/avro.html
    props.setProperty("key.converter.schemas.enable", "false");
    props.setProperty("value.converter.schemas.enable", "false");

    // debezium names
    props.setProperty("name", config.getDatabase());
    props.setProperty("database.server.name", config.getDatabase());

    // db connection configuration
    props.setProperty("database.hostname", config.getHost());
    props.setProperty("database.port", config.getPort());
    props.setProperty("database.user", config.getUsername());
    props.setProperty("database.dbname", config.getDatabase());
    props.setProperty("database.password", config.getPassword());

    // https://debezium.io/documentation/reference/1.9/connectors/mysql.html#mysql-property-binary-handling-mode
    props.setProperty("binary.handling.mode", "base64");

    // table selection
    props.setProperty("database.include.list", config.getDatabase());
    if (!config.getIncludeTables().isEmpty()) {
      props.setProperty("table.include.list", tableFormat(config.getIncludeTables().stream()));
    } else {
      props.setProperty(
          "table.exclude.list", tableFormat(getExcludedTables(config.getExcludeTables()).stream()));
    }

    props.putAll(getDebeziumProperties());
    return props;
  }

  public Set<String> getExcludedTables(Set<String> excludeTables) {
    Set<String> exclude = new HashSet<>(getSystemExcludedTables());
    exclude.addAll(excludeTables);
    return exclude;
  }

  public Set<String> getSystemExcludedTables() {
    return new HashSet<>(Arrays.asList("information_schema", "mysql", "performance_schema", "sys"));
  }

  public String tableFormat(Stream<String> table) {
    return table
        .map(stream -> config.getDatabase() + "." + stream)
        .collect(Collectors.joining(","));
  }
}

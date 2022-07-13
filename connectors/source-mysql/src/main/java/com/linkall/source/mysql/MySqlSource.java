package com.linkall.source.mysql;

import com.linkall.vance.common.env.EnvUtil;
import com.linkall.vance.core.Adapter;
import com.linkall.vance.core.Adapter1;
import com.linkall.vance.core.Source;
import io.debezium.embedded.Connect;
import io.debezium.engine.ChangeEvent;
import io.debezium.engine.DebeziumEngine;
import io.debezium.engine.spi.OffsetCommitPolicy;
import org.apache.kafka.connect.source.SourceRecord;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.Arrays;
import java.util.HashSet;
import java.util.Properties;
import java.util.Set;
import java.util.concurrent.Executor;
import java.util.concurrent.Executors;
import java.util.stream.Collectors;
import java.util.stream.Stream;

public class MySqlSource implements Source {
  private static final Logger LOGGER = LoggerFactory.getLogger(MySqlSource.class);
  private DebeziumEngine<ChangeEvent<SourceRecord, SourceRecord>> engine;
  private Executor executor;
  private final DebeziumEngine.ChangeConsumer<ChangeEvent<SourceRecord, SourceRecord>> consumer;
  private final MySqlConfig config;

  public MySqlSource() {
    consumer = new DebeziumRecordConsumer((Adapter1<SourceRecord>) getAdapter());
    config =
        new MySqlConfig(
            EnvUtil.getEnvOrConfig("host"),
            EnvUtil.getEnvOrConfig("port"),
            EnvUtil.getEnvOrConfig("username"),
            EnvUtil.getEnvOrConfig("password"),
            EnvUtil.getEnvOrConfig("database"),
            EnvUtil.getEnvOrConfig("include_table"),
            EnvUtil.getEnvOrConfig("exclude_table"));
  }

  @Override
  public Adapter getAdapter() {
    return new MySqlAdapter();
  }

  @Override
  public void start() throws Exception {
    engine =
        DebeziumEngine.create(Connect.class)
            .using(getDebeziumProperties())
            .using(OffsetCommitPolicy.always())
            .notifying(consumer)
            .using(
                (success, message, error) -> {
                  LOGGER.info("Debezium engine shutdown.");
                })
            .build();
    executor = Executors.newSingleThreadExecutor();
    executor.execute(engine);
  }

  private Properties getDebeziumProperties() {
    final Properties props = new Properties();

    // debezium engine configuration
    props.setProperty("connector.class", "io.debezium.connector.mysql.MySqlConnector");
    // snapshot config
    props.setProperty("snapshot.mode", "initial");
    // DO NOT include schema change, e.g. DDL
    props.setProperty("include.schema.changes", "false");
    // disable tombstones
    props.setProperty("tombstones.on.delete", "false");

    // offset
    props.setProperty("offset.storage", KvStoreOffsetBackingStore.class.getCanonicalName());
    //        props.setProperty("offset.storage.file.filename", "/tmp/offset.dat");
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
    if (config.getIncludeTables() != null && config.getIncludeTables().length > 0) {
      props.setProperty(
          "table.include.list", tableFormat(Arrays.stream(config.getIncludeTables())));
    } else {
      props.setProperty(
          "table.exclude.list", tableFormat(getExcludedTables(config.getExcludeTables()).stream()));
    }
    props.setProperty("database.include.list", config.getDatabase());

    // convert
    props.setProperty("converters", "boolean, datetime");
    props.setProperty(
        "boolean.type", "io.debezium.connector.mysql.converters.TinyIntOneToBooleanConverter");
    props.setProperty("datetime.type", "com.linkall.source.mysql.MySqlDateTimeConverter");

    return props;
  }

  public Set<String> getExcludedTables(String[] excludeTables) {
    Set<String> exclude =
        new HashSet<>(Arrays.asList("information_schema", "mysql", "performance_schema", "sys"));
    if (excludeTables != null) {
      for (String name : excludeTables) {
        exclude.add(name);
      }
    }
    return exclude;
  }

  public String tableFormat(Stream<String> table) {
    return table
        .map(stream -> config.getDatabase() + "." + stream)
        .collect(Collectors.joining(","));
  }
}

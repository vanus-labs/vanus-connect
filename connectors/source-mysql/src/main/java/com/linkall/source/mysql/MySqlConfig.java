package com.linkall.source.mysql;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.linkall.cdk.database.debezium.DebeziumConfig;
import com.linkall.cdk.database.debezium.KvStoreOffsetBackingStore;
import io.debezium.connector.mysql.MySqlConnector;
import io.debezium.connector.mysql.converters.TinyIntOneToBooleanConverter;
import io.debezium.storage.file.history.FileSchemaHistory;
import org.apache.kafka.connect.json.JsonConverter;
import org.apache.kafka.connect.json.JsonConverterConfig;
import org.apache.kafka.connect.storage.Converter;

import java.nio.charset.StandardCharsets;
import java.util.Arrays;
import java.util.HashMap;
import java.util.Map;
import java.util.Properties;
import java.util.stream.Collectors;

public class MySqlConfig extends DebeziumConfig {

  @JsonProperty("name")
  private String name;

  @JsonProperty("db")
  private DbConfig dbConfig;

  @JsonProperty("binlog_offset")
  private BinlogOffset binlogOffset;

  @JsonProperty("db_history_file")
  private String dbHistoryFile;

  @JsonProperty("include_databases")
  private String[] includeDatabases;

  @JsonProperty("exclude_databases")
  private String[] excludeDatabases;

  @JsonProperty("include_tables")
  private String[] includeTables;

  @JsonProperty("exclude_tables")
  private String[] excludeTables;

  @Override
  public Class<?> secretClass() {
    return DbConfig.class;
  }

  @Override
  protected Properties getDebeziumProperties() {
    final Properties props = new Properties();

    // debezium engine configuration
    props.setProperty("connector.class", MySqlConnector.class.getCanonicalName());
    // snapshot config
    props.setProperty("snapshot.mode", "initial");
    // DO NOT include schema change, e.g. DDL
    props.setProperty("include.schema.changes", "false");
    // disable tombstones
    props.setProperty("tombstones.on.delete", "false");

    if (binlogOffset != null) {
      Converter valueConverter = new JsonConverter();
      Map<String, Object> valueConfigs = new HashMap<>();
      valueConfigs.put(JsonConverterConfig.SCHEMAS_ENABLE_CONFIG, false);
      valueConverter.configure(valueConfigs, false);
      byte[] offsetValue = valueConverter.fromConnectData(name, null, binlogOffset);
      props.setProperty(
          KvStoreOffsetBackingStore.OFFSET_CONFIG_VALUE,
          new String(offsetValue, StandardCharsets.UTF_8));
    }

    props.setProperty("offset.flush.interval.ms", "1000");

    // history
    props.setProperty("schema.history.internal", FileSchemaHistory.class.getCanonicalName());
    props.setProperty("schema.history.internal.file.filename", dbHistoryFile);

    // https://debezium.io/documentation/reference/configuration/avro.html
    props.setProperty("key.converter.schemas.enable", "false");
    props.setProperty("value.converter.schemas.enable", "false");

    // debezium names
    props.setProperty("name", name);
    props.setProperty("topic.prefix", name);
    props.setProperty("database.server.id", String.valueOf(System.currentTimeMillis()));

    // db connection configuration
    props.setProperty("database.hostname", dbConfig.getHost());
    props.setProperty("database.port", String.valueOf(dbConfig.getPort()));
    props.setProperty("database.user", dbConfig.getUsername());
    props.setProperty("database.password", dbConfig.getPassword());

    // https://debezium.io/documentation/reference/2.0/connectors/mysql.html#mysql-property-binary-handling-mode
    props.setProperty("binary.handling.mode", "base64");

    if (includeDatabases != null
        && includeDatabases.length > 0
        && excludeDatabases != null
        && excludeDatabases.length > 0) {
      throw new IllegalArgumentException(
          "the include_databases and exclude_databases can't be set together");
    }
    // database selection
    if (includeDatabases != null && includeDatabases.length > 0) {
      props.setProperty(
          "database.include.list",
          Arrays.stream(includeDatabases).collect(Collectors.joining(",")));
    } else if (excludeDatabases != null && excludeDatabases.length > 0) {
      props.setProperty(
          "database.exclude.list",
          Arrays.stream(excludeDatabases).collect(Collectors.joining(",")));
    }

    if (includeTables != null
        && includeTables.length > 0
        && excludeTables != null
        && excludeTables.length > 0) {
      throw new IllegalArgumentException(
          "the include_tables and exclude_tables can't be set together");
    }
    if (includeTables != null && includeTables.length > 0) {
      props.setProperty(
          "table.include.list", Arrays.stream(includeTables).collect(Collectors.joining(",")));
    } else if (excludeTables != null && excludeTables.length > 0) {
      props.setProperty(
          "table.exclude.list", Arrays.stream(excludeTables).collect(Collectors.joining(",")));
    }
    props.setProperty("converters", "boolean, datetime");
    props.setProperty("boolean.type", TinyIntOneToBooleanConverter.class.getCanonicalName());
    props.setProperty("datetime.type", MySqlDateTimeConverter.class.getCanonicalName());
    return props;
  }
}

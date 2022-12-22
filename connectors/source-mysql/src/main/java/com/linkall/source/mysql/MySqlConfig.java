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
import java.util.*;

public class MySqlConfig extends DebeziumConfig {
  private static final Set<String> systemTable = new HashSet(Arrays.asList("information_schema", "mysql", "performance_schema", "sys"));

  @JsonProperty("db_config")
  private DbConfig dbConfig;

  @JsonProperty("binlog_offset")
  private BinlogOffset binlogOffset;

  @JsonProperty("db_history_file")
  private String dbHistoryFile;

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
      byte[] offsetValue =
          valueConverter.fromConnectData(dbConfig.getDatabase(), null, binlogOffset);
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
    props.setProperty("name", dbConfig.getDatabase());
    props.setProperty("topic.prefix", dbConfig.getDatabase());
    props.setProperty("database.server.id", String.valueOf(System.currentTimeMillis()));

    // db connection configuration
    props.setProperty("database.hostname", dbConfig.getHost());
    props.setProperty("database.port", String.valueOf(dbConfig.getPort()));
    props.setProperty("database.dbname", dbConfig.getDatabase());
    if (dbConfig.getUsername() != null && dbConfig.getUsername() != "")
      props.setProperty("database.user", dbConfig.getUsername());
    if (dbConfig.getPassword() != null && dbConfig.getPassword() != "")
      props.setProperty("database.password", dbConfig.getPassword());

    // https://debezium.io/documentation/reference/1.9/connectors/mysql.html#mysql-property-binary-handling-mode
    props.setProperty("binary.handling.mode", "base64");

    // table selection
    props.setProperty("database.include.list", dbConfig.getDatabase());
    if (includeTables != null && includeTables.length > 0) {
      props.setProperty(
          "table.include.list", tableFormat(dbConfig.getDatabase(), Arrays.stream(includeTables)));
    } else {
      Set<String> exclude = new HashSet<>(systemTable);
      if (excludeTables !=null && excludeTables.length>0){
        for (String table: excludeTables){
          exclude.add(table);
        }
      }
      props.setProperty(
          "table.exclude.list", tableFormat(dbConfig.getDatabase(),exclude.stream()));
    }
    props.setProperty("converters", "boolean, datetime");
    props.setProperty("boolean.type", TinyIntOneToBooleanConverter.class.getCanonicalName());
    props.setProperty("datetime.type", MySqlDateTimeConverter.class.getCanonicalName());
    return props;
  }

  public DbConfig getDbConfig() {
    return dbConfig;
  }

  public void setDbConfig(DbConfig dbConfig) {
    this.dbConfig = dbConfig;
  }

  public BinlogOffset getBinlogOffset() {
    return binlogOffset;
  }

  public void setBinlogOffset(BinlogOffset binlogOffset) {
    this.binlogOffset = binlogOffset;
  }

  public String getDbHistoryFile() {
    return dbHistoryFile;
  }

  public void setDbHistoryFile(String dbHistoryFile) {
    this.dbHistoryFile = dbHistoryFile;
  }

  public String[] getIncludeTables() {
    return includeTables;
  }

  public void setIncludeTables(String[] includeTables) {
    this.includeTables = includeTables;
  }

  public String[] getExcludeTables() {
    return excludeTables;
  }

  public void setExcludeTables(String[] excludeTables) {
    this.excludeTables = excludeTables;
  }
}

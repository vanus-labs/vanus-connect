package com.linkall.source.postgresql;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.linkall.cdk.database.debezium.DebeziumConfig;
import com.linkall.cdk.database.debezium.KvStoreOffsetBackingStore;
import io.debezium.connector.postgresql.PostgresConnector;
import org.apache.kafka.connect.json.JsonConverter;
import org.apache.kafka.connect.json.JsonConverterConfig;
import org.apache.kafka.connect.storage.Converter;

import java.nio.charset.StandardCharsets;
import java.util.Arrays;
import java.util.HashMap;
import java.util.Map;
import java.util.Properties;
import java.util.stream.Collectors;

public class PostgreSQLConfig extends DebeziumConfig {

    @JsonProperty("name")
    private String name;

    @JsonProperty("db")
    private DbConfig dbConfig;

    @JsonProperty("offset")
    private PostgreSQLOffset offset;

    @JsonProperty("plugin_name")
    private String pluginName;

    @JsonProperty("slot_name")
    private String slotName;

    @JsonProperty("publication_name")
    private String publicationName;

    @JsonProperty("schema_include")
    private String[] schemaInclude;

    @JsonProperty("schema_exclude")
    private String[] schemaExclude;

    @JsonProperty("table_include")
    private String[] tableInclude;

    @JsonProperty("table_exclude")
    private String[] tableExclude;

    @JsonProperty("column_include")
    private String[] columnInclude;

    @JsonProperty("column_exclude")
    private String[] columnExclude;

    @Override
    public Class<?> secretClass() {
        return DbConfig.class;
    }

    @Override
    protected Properties getDebeziumProperties() {
        final Properties props = new Properties();

        // debezium engine configuration
        props.setProperty("connector.class", PostgresConnector.class.getCanonicalName());
        // snapshot config
        props.setProperty("snapshot.mode", "initial");
        // DO NOT include schema change, e.g. DDL
        props.setProperty("include.schema.changes", "false");
        // disable tombstones
        props.setProperty("tombstones.on.delete", "false");

        if (offset!=null) {
            Converter valueConverter = new JsonConverter();
            Map<String, Object> valueConfigs = new HashMap<>();
            valueConfigs.put(JsonConverterConfig.SCHEMAS_ENABLE_CONFIG, false);
            valueConverter.configure(valueConfigs, false);
            byte[] offsetValue = valueConverter.fromConnectData(name, null, offset);
            props.setProperty(
                    KvStoreOffsetBackingStore.OFFSET_CONFIG_VALUE,
                    new String(offsetValue, StandardCharsets.UTF_8));
        }

        // https://debezium.io/documentation/reference/configuration/avro.html
        props.setProperty("key.converter.schemas.enable", "false");
        props.setProperty("value.converter.schemas.enable", "false");

        // debezium names
        props.setProperty("name", name);
        props.setProperty("topic.prefix", name);

        if (pluginName==null) {
            pluginName = "pgoutput";
        }
        props.setProperty("plugin.name", pluginName);
        if (slotName!=null)
            props.setProperty("slot.name", slotName);
        if (publicationName!=null)
            props.setProperty("publication.name", publicationName);

        // db connection configuration
        props.setProperty("database.hostname", dbConfig.getHost());
        props.setProperty("database.port", String.valueOf(dbConfig.getPort()));
        props.setProperty("database.user", dbConfig.getUsername());
        props.setProperty("database.password", dbConfig.getPassword());
        props.setProperty("database.dbname", dbConfig.getDatabase());

        // https://debezium.io/documentation/reference/stable/connectors/postgresql.html#postgresql-property-binary-handling-mode
        props.setProperty("binary.handling.mode", "base64");

        if (schemaInclude!=null
                && schemaInclude.length > 0
                && schemaExclude!=null
                && schemaExclude.length > 0) {
            throw new IllegalArgumentException(
                    "the schema_include and schema_exclude can't be set together");
        }
        // database selection
        if (schemaInclude!=null && schemaInclude.length > 0) {
            props.setProperty(
                    "schema.include.list", Arrays.stream(schemaInclude).collect(Collectors.joining(",")));
        } else if (schemaExclude!=null && schemaExclude.length > 0) {
            props.setProperty(
                    "schema.exclude.list", Arrays.stream(schemaExclude).collect(Collectors.joining(",")));
        }

        if (tableInclude!=null
                && tableInclude.length > 0
                && tableExclude!=null
                && tableExclude.length > 0) {
            throw new IllegalArgumentException(
                    "the table_include and table_exclude can't be set together");
        }
        if (tableInclude!=null && tableInclude.length > 0) {
            props.setProperty(
                    "table.include.list", Arrays.stream(tableInclude).collect(Collectors.joining(",")));
        } else if (tableExclude!=null && tableExclude.length > 0) {
            props.setProperty(
                    "table.exclude.list", Arrays.stream(tableExclude).collect(Collectors.joining(",")));
        }
        if (columnInclude!=null
                && columnInclude.length > 0
                && columnExclude!=null
                && columnExclude.length > 0) {
            throw new IllegalArgumentException(
                    "the column_include and column_exclude can't be set together");
        }
        if (columnInclude!=null && columnInclude.length > 0) {
            props.setProperty(
                    "column.include.list", Arrays.stream(columnInclude).collect(Collectors.joining(",")));
        } else if (columnExclude!=null && columnExclude.length > 0) {
            props.setProperty(
                    "column.exclude.list", Arrays.stream(columnExclude).collect(Collectors.joining(",")));
        }
        return props;
    }
}

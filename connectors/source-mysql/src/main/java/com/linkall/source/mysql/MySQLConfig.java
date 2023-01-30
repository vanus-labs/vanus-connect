package com.linkall.source.mysql;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.linkall.cdk.database.debezium.DebeziumConfig;
import io.debezium.connector.mysql.MySqlConnector;
import io.debezium.storage.file.history.FileSchemaHistory;

import java.util.Arrays;
import java.util.Properties;
import java.util.stream.Collectors;

public class MySQLConfig extends DebeziumConfig {

    @JsonProperty("name")
    private String name;

    @JsonProperty("db")
    private DbConfig dbConfig;

    @JsonProperty("binlog_offset")
    private BinlogOffset binlogOffset;

    @JsonProperty("db_history_file")
    private String dbHistoryFile;

    @JsonProperty("database_include")
    private String[] databaseInclude;

    @JsonProperty("database_exclude")
    private String[] databaseExclude;

    @JsonProperty("table_include")
    private String[] tableInclude;

    @JsonProperty("table_exclude")
    private String[] tableExclude;

    @Override
    public Class<?> secretClass() {
        return DbConfig.class;
    }

    @Override
    protected Object getOffset() {
        return binlogOffset;
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

        if (databaseInclude!=null && databaseInclude.length > 0
                && databaseExclude!=null && databaseExclude.length > 0) {
            throw new IllegalArgumentException(
                    "the database_include and database_exclude can't be set together");
        }
        // database selection
        if (databaseInclude!=null && databaseInclude.length > 0) {
            props.setProperty(
                    "database.include.list", Arrays.stream(databaseInclude).collect(Collectors.joining(",")));
        } else if (databaseExclude!=null && databaseExclude.length > 0) {
            props.setProperty(
                    "database.exclude.list", Arrays.stream(databaseExclude).collect(Collectors.joining(",")));
        }

        if (tableInclude!=null && tableInclude.length > 0
                && tableExclude!=null && tableExclude.length > 0) {
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
//    props.setProperty("converters", "boolean, datetime");
//    props.setProperty("boolean.type", TinyIntOneToBooleanConverter.class.getCanonicalName());
//    props.setProperty("datetime.type", MySQLDateTimeConverter.class.getCanonicalName());
        return props;
    }
}

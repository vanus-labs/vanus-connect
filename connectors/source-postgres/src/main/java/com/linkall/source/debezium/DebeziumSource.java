package com.linkall.source.debezium;

import com.linkall.vance.common.config.ConfigLoader;
import com.linkall.vance.core.Adapter1;
import com.linkall.vance.core.Source;
import io.debezium.embedded.Connect;
import io.debezium.engine.ChangeEvent;
import io.debezium.engine.DebeziumEngine;
import io.debezium.engine.spi.OffsetCommitPolicy;
import io.vertx.core.json.JsonObject;
import org.apache.kafka.connect.json.JsonConverter;
import org.apache.kafka.connect.json.JsonConverterConfig;
import org.apache.kafka.connect.source.SourceRecord;
import org.apache.kafka.connect.storage.Converter;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.IOException;
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
    protected JsonObject config;

    private DbConfig dbConfig;


    public DebeziumSource() {
        config = ConfigLoader.getUserConfig();
        dbConfig = new DbConfig(config.getString("host"),
                config.getString("port"),
                config.getString("username"),
                config.getString("password"),
                config.getString("db_name"));
        consumer = new DebeziumChangeConsumer((Adapter1<SourceRecord>) getAdapter());
    }

    public abstract String getConnectorClass();

    public abstract Map<String, Object> getConfigOffset();

    public abstract Properties getDebeziumProperties();

    public abstract Set<String> getSystemExcludedTables();

    public void start() throws Exception {
        engine =
                DebeziumEngine.create(Connect.class)
                        .using(getProperties())
                        .using(OffsetCommitPolicy.always())
                        .notifying(consumer)
                        .using(
                                (success, message, error) -> {
                                    LOGGER.info(
                                            "Debezium engine shutdown,success: {},message: {}", success, message, error);
                                })
                        .build();
        executor = Executors.newSingleThreadExecutor();
        executor.execute(engine);
        Runtime.getRuntime()
                .addShutdownHook(
                        new Thread(
                                () -> {
                                    try {
                                        engine.close();
                                    } catch (IOException e) {
                                        LOGGER.error("engine close error", e);
                                    }
                                }));
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
        if (config.containsKey("store_offset_key")) {
            props.setProperty(
                    KvStoreOffsetBackingStore.OFFSET_STORAGE_KV_STORE_KEY_CONFIG, config.getString("store_offset_key"));
        }
        Map<String, Object> configOffset = getConfigOffset();
        if (configOffset!=null && configOffset.size() > 0) {
            Converter valueConverter = new JsonConverter();
            Map<String, Object> valueConfigs = new HashMap<>();
            valueConfigs.put(JsonConverterConfig.SCHEMAS_ENABLE_CONFIG, false);
            valueConverter.configure(valueConfigs, false);
            byte[] offsetValue = valueConverter.fromConnectData(dbConfig.getDatabase(), null, configOffset);
            props.setProperty(
                    KvStoreOffsetBackingStore.OFFSET_CONFIG_VALUE,
                    new String(offsetValue, StandardCharsets.UTF_8));
        }

        props.setProperty("offset.flush.interval.ms", "1000");

        // https://debezium.io/documentation/reference/configuration/avro.html
        props.setProperty("key.converter.schemas.enable", "false");
        props.setProperty("value.converter.schemas.enable", "false");

        // debezium names
        props.setProperty("name", dbConfig.getDatabase());
        props.setProperty("database.server.name", dbConfig.getDatabase());

        // db connection configuration
        props.setProperty("database.hostname", dbConfig.getHost());
        props.setProperty("database.port", dbConfig.getPort());
        props.setProperty("database.user", dbConfig.getUsername());
        props.setProperty("database.dbname", dbConfig.getDatabase());
        props.setProperty("database.password", dbConfig.getPassword());

        props.putAll(getDebeziumProperties());
        if (config.containsKey("debezium")) {
            // other debezium properties
            props.putAll(getDebeziumProperties(config.getJsonObject("debezium")));
        }
        return props;
    }

    public Properties getDebeziumProperties(JsonObject debezium) {
        final Properties debeziumProperties = new Properties();

        debezium.stream().forEach(k -> {
            debeziumProperties.put(k.getKey(), debezium.getString(k.getKey()));
        });

        return debeziumProperties;
    }

    public Set<String> getExcludedTables(Set<String> excludeTables) {
        Set<String> exclude = new HashSet<>(getSystemExcludedTables());
        exclude.addAll(excludeTables);
        return exclude;
    }

    public String tableFormat(String name, Stream<String> table) {
        return table
                .map(stream -> name + "." + stream)
                .collect(Collectors.joining(","));
    }
}

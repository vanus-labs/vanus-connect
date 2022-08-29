// Copyright 2022 Linkall Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package com.linkall.connector.source.mongodb.debezium;

import com.linkall.vance.core.Adapter2;
import io.debezium.engine.ChangeEvent;
import io.debezium.engine.DebeziumEngine;
import io.debezium.engine.format.Json;
import io.debezium.engine.spi.OffsetCommitPolicy;
import org.apache.kafka.connect.json.JsonConverter;
import org.apache.kafka.connect.json.JsonConverterConfig;
import org.apache.kafka.connect.storage.Converter;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.IOException;
import java.nio.charset.StandardCharsets;
import java.util.HashMap;
import java.util.Map;
import java.util.Properties;
import java.util.concurrent.Executor;
import java.util.concurrent.Executors;

public abstract class DebeziumSource implements com.linkall.vance.core.Source {
    private static final Logger LOGGER = LoggerFactory.getLogger(DebeziumSource.class);
    private final DebeziumEngine.ChangeConsumer<ChangeEvent<String, String>> consumer;
    private DebeziumEngine<ChangeEvent<String, String>> engine;
    private Executor executor;

    public DebeziumSource() {
        consumer = new RecordConsumer((Adapter2<String, String>) getAdapter());
    }

    public abstract String getConnectorClass();

    public abstract String getDatabase();

    public abstract String getStoreOffsetKey();

    public abstract Map<String, Object> getConfigOffset();

    public abstract Properties getDebeziumProperties() throws IOException;

    public void start() throws IOException {
        engine =
                DebeziumEngine.create(Json.class)
                        .using(getProperties())
                        .using(OffsetCommitPolicy.always())
                        .notifying(consumer)
                        .using(
                                (success, message, error) -> {
                                    LOGGER.info(
                                            "Debezium engine shutdown,success: {}, message: {},error:{}",
                                            success,
                                            message,
                                            error);
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

    private Properties getProperties() throws IOException {
        final Properties props = new Properties();

        // debezium engine configuration
        props.setProperty("connector.class", getConnectorClass());

        // https://debezium.io/documentation/reference/1.9/connectors/mysql.html#mysql-property-binary-handling-mode
        props.setProperty("binary.handling.mode", "base64");

        // snapshot config
        props.setProperty("snapshot.mode", "initial");
        // DO NOT include schema change, e.g. DDL
        props.setProperty("include.schema.changes", "false");
        // disable tombstones
        props.setProperty("tombstones.on.delete", "false");
        props.setProperty("converter.schemas.enable", "false"); // don't include schema in message

        // https://debezium.io/documentation/reference/stable/integrations/cloudevents.html
        props.setProperty("converter", "io.debezium.converters.CloudEventsConverter");
        props.setProperty("converter.serializer.type", "json");
        props.setProperty("converter.data.serializer.type", "json");

        // history
        props.setProperty("database.history", "io.debezium.relational.history.FileDatabaseHistory");
        props.setProperty("database.history.file.filename", "/tmp/mongodb/history.data");

        // offset
        props.setProperty("offset.flush.interval.ms", "1000");
        props.setProperty("offset.storage", OffsetStore.class.getCanonicalName());
        if (this.getStoreOffsetKey() != null && !this.getStoreOffsetKey().isEmpty()) {
            props.setProperty(
                    OffsetStore.OFFSET_STORAGE_KV_STORE_KEY_CONFIG, this.getStoreOffsetKey());
        }
        Map<String, Object> configOffset = getConfigOffset();
        if (configOffset != null && configOffset.size() > 0) {
            Converter valueConverter = new JsonConverter();
            Map<String, Object> valueConfigs = new HashMap<>();
            valueConfigs.put(JsonConverterConfig.SCHEMAS_ENABLE_CONFIG, false);
            valueConverter.configure(valueConfigs, false);
            byte[] offsetValue = valueConverter.fromConnectData(this.getDatabase(), null, configOffset);
            props.setProperty(
                    OffsetStore.OFFSET_CONFIG_VALUE,
                    new String(offsetValue, StandardCharsets.UTF_8));
        }

        props.putAll(getDebeziumProperties());
        return props;
    }
}

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

package com.linkall.connector.mongodb;

import com.alibaba.fastjson.JSON;
import com.alibaba.fastjson.JSONObject;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.google.protobuf.Struct;
import com.google.protobuf.Value;
import com.google.protobuf.util.JsonFormat;
import com.linkall.connector.proto.base.Base;
import com.linkall.connector.proto.database.Database;
import com.linkall.connector.proto.database.Mongodb;
import com.linkall.vance.core.Adapter2;
import io.cloudevents.CloudEvent;
import io.cloudevents.core.builder.CloudEventBuilder;
import io.cloudevents.jackson.JsonCloudEventData;
import org.apache.commons.text.StringEscapeUtils;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.net.URI;
import java.nio.charset.StandardCharsets;
import java.time.Instant;
import java.time.OffsetDateTime;
import java.time.ZoneOffset;
import java.util.ArrayList;
import java.util.List;
import java.util.Map;

public class MongoDBAdapter implements Adapter2 {
    private static final Logger LOGGER = LoggerFactory.getLogger(MongoDBAdapter.class);
    private static final String EXTENSION_NAME_PREFIX = "vancemongodb";
    private static final List<String> keyFilter = new ArrayList<>();

    static {
        keyFilter.add("db");
        keyFilter.add("collection");
        keyFilter.add("connector");
        keyFilter.add("rs");
        keyFilter.add("ts_ms");
    }

    public static Mongodb.Event parse(String key, String value) {
        Mongodb.Event.Builder builder = Mongodb.Event.newBuilder();
        builder.getMetadataBuilder().setRecognized(true);
        builder.getRawBuilder().setKey(key);
        builder.getRawBuilder().setValue(value);
        try {
            JSONObject obj = JSON.parseObject(value);
            Base.Metadata.Builder mdBuilder = builder.getMetadataBuilder();
            Struct.Builder sb = Struct.newBuilder();
            JsonFormat.parser().merge(obj.getJSONObject("source").toJSONString(), sb);
            mdBuilder.setExtension(sb.build());
            JSONObject id = JSON.parseObject(key);
            if (id.containsKey("id")) {
                mdBuilder.setId(id.getJSONObject("id").get("$oid").toString());
            } else {
                mdBuilder.setRecognized(false);
                return builder.build();
            }
            switch (obj.getString("op")) {
                case "u":
                    builder.setOp(Database.Operation.UPDATE);
                    JSONObject ud = obj.getJSONObject("updateDescription");

                    JsonFormat.parser().merge(
                            StringEscapeUtils.unescapeJava(
                                    ud.getOrDefault("updatedFields", "{}").toString()
                            ),
                            builder.getUpdateBuilder().getUpdateDescriptionBuilder().getUpdatedFieldsBuilder());
                    JsonFormat.parser().merge(
                            StringEscapeUtils.unescapeJava(
                                    ud.getOrDefault("truncatedArrays", "[]").toString()
                            ),
                            builder.getUpdateBuilder().getUpdateDescriptionBuilder().getTruncatedArraysBuilder());
                    JsonFormat.parser().merge(
                            StringEscapeUtils.unescapeJava(
                                    ud.getOrDefault("removedFields", "[]").toString()
                            ),
                            builder.getUpdateBuilder().getUpdateDescriptionBuilder().getRemovedFieldsBuilder());
                case "c":
                    // TODO don't use this for update
                    if (builder.getOp() == Database.Operation.UNKNOWN) {
                        builder.setOp(Database.Operation.INSERT);
                    }
                    JsonFormat.parser().merge(obj.getString("after"),
                            builder.getInsertBuilder().getDocumentBuilder());
                    break;
                case "d":
                    builder.setOp(Database.Operation.DELETE);
                    break;
                default:
                    builder.setOp(Database.Operation.UNKNOWN);
                    mdBuilder.setRecognized(false);
                    return builder.build();
            }
        } catch (Exception e) {
            LOGGER.warn("parse event data failed: {}", e.getMessage());
            builder.getMetadataBuilder().setRecognized(false);
        }
        return builder.build();
    }

    public static CloudEvent proto2CloudEvent(Mongodb.Event event) {
        CloudEventBuilder builder = CloudEventBuilder.v1();

        String ID = "unknown";
        String sourcePrefix = "unknown.unknown";
        String type = "unknown.unknown";
        builder.withId(ID).withType(type).withSource(URI.create(sourcePrefix + "." + type));
        String data = "{\"raw\":{\"key\":\"" + event.getRaw().getKey() + "\",\"value\":\"" + event.getRaw().getValue() + "\"}}\n";
        try {
            builder.withDataContentType("application/json");
            Base.Metadata md = event.getMetadata();

            if (md.getRecognized()) {
                ID = event.getMetadata().getId();
                sourcePrefix = md.getExtension().getFieldsMap().get("connector").getStringValue() + "."
                        + md.getExtension().getFieldsMap().get("rs").getStringValue();

                type = md.getExtension().getFieldsMap().get("db").getStringValue() + "." +
                        md.getExtension().getFieldsMap().get("collection").getStringValue();
                builder.withId(ID).withType(type).withSource(URI.create(sourcePrefix + "." + type));
            }

            Value time = md.getExtension().getFieldsMap().get("ts_ms");
            if (time != null) {
                builder.withTime(
                        OffsetDateTime.ofInstant(
                                Instant.ofEpochMilli(
                                        ((Double) time.getNumberValue()).longValue()
                                ),
                                ZoneOffset.UTC)
                );
            }
            Mongodb.Event.Builder b = event.toBuilder();
            b.clearRaw();
            event.toBuilder().clearRaw();
            ObjectMapper mapper = new ObjectMapper();
            builder.withData(JsonCloudEventData.wrap(
                    mapper.readTree(
                            JsonFormat.printer().omittingInsignificantWhitespace().print(b.build())
                    )
            ));
            for (Map.Entry<String, Value> entry : md.getExtension().getFieldsMap().entrySet()) {
                if (!MongoDBAdapter.keyFilter.contains(entry.getKey()) && entry.getValue() != null) {
                    builder.withExtension(EXTENSION_NAME_PREFIX + entry.getKey(), entry.getValue().getStringValue());
                }
            }

            builder.withExtension(EXTENSION_NAME_PREFIX + "operation", event.getOp().toString());
        } catch (Exception e) {
            builder.withData(data.getBytes(StandardCharsets.UTF_8));
            e.printStackTrace();
        }
        builder.withExtension(EXTENSION_NAME_PREFIX + "recognized", event.getMetadata().getRecognized());
        return builder.build();
    }

    @Override
    public CloudEvent adapt(Object key, Object val) {
        return MongoDBAdapter.proto2CloudEvent(MongoDBAdapter.parse((String) key, (String) val));
    }
}


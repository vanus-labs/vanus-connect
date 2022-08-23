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
import com.alibaba.fastjson.JSONValidator;
import io.cloudevents.CloudEvent;
import io.cloudevents.core.builder.CloudEventBuilder;
import io.cloudevents.core.data.BytesCloudEventData;
import org.apache.commons.lang3.StringUtils;
import org.apache.commons.text.StringEscapeUtils;
import org.apache.logging.log4j.util.Strings;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.net.URI;
import java.time.Instant;
import java.time.OffsetDateTime;
import java.time.ZoneOffset;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

class MongoChangeEvent {
    private static final Logger LOGGER = LoggerFactory.getLogger(MongoChangeEvent.class);
    private static final String EXTENSION_NAME_PREFIX = "vancemongodb";
    private static final List<String> keyFilter = new ArrayList<>();

    static {
        keyFilter.add("db");
        keyFilter.add("collection");
        keyFilter.add("connector");
        keyFilter.add("rs");
        keyFilter.add("ts_ms");
    }

    private final HashMap<String, Object> fullFields = new HashMap<>();
    private final HashMap<String, Object> updatedFields = new HashMap<>();
    private final HashMap<String, Object> deletedFields = new HashMap<>();
    private Map<String, Object> metadata = new HashMap<>();
    private String objectID;
    private String type;
    private boolean isValidate = true;
    private String rawKey;
    private String rawValue;

    public static MongoChangeEvent parse(String key, String value) {
        MongoChangeEvent event = new MongoChangeEvent();
        event.rawKey = key;
        event.rawValue = value;
        try {
            JSONObject obj = JSON.parseObject(value);
            event.metadata = obj.getJSONObject("source").getInnerMap();
            JSONObject id = JSON.parseObject(key);
            if (id.containsKey("id")) {
                event.objectID = id.getJSONObject("id").get("$oid").toString();
            } else {
                event.isValidate = false;
                return event;
            }
            switch (obj.getString("op")) {
                case "c":
                    event.type = "insert";
                    processFullFields(obj.getString("after"), event);
                    break;
                case "u":
                    event.type = "update";
                    processFullFields(obj.getString("after"), event);
                    processUpdateFields(obj.getString("updateDescription"), event);
                    break;
                case "d":
                    event.type = "delete";
                    break;
                default:
                    event.isValidate = false;
                    return event;
            }
        } catch (Exception e) {
            LOGGER.warn("parse event data failed: {}", e.getMessage());
            event.isValidate = false;
            if (Strings.isBlank(event.type)) {
                event.type = "unknown";
            }
        }
        return event;
    }

    private static void processFullFields(String data, MongoChangeEvent event) {
        String body = StringEscapeUtils.unescapeJava(data);
        event.fullFields.put("_id", event.objectID);
        JSONValidator validator = JSONValidator.from(body);
        if (!validator.validate()) {
            event.isValidate = false;
            return;
        }
        JSONObject obj = JSON.parseObject(body);
        for (Map.Entry<String, Object> entry : obj.entrySet()) {
            String key = entry.getKey();
            String val = entry.getValue().toString();
            if (!key.equals("_id")) {
                event.fullFields.put(key, val);
            }
        }
    }

    private static void processUpdateFields(String data, MongoChangeEvent event) {
        JSONObject obj = JSON.parseObject(data);
        String updated = StringEscapeUtils.unescapeJava(obj.getString("updatedFields"));
        if (!StringUtils.isBlank(updated)) {
            JSONValidator validator = JSONValidator.from(updated);
            if (!validator.validate()) {
                event.isValidate = false;
                return;
            }
            event.updatedFields.put("updated", JSON.parseObject(updated, Map.class));
        }

        String removed = StringEscapeUtils.unescapeJava(obj.getString("removedFields"));
        if (!StringUtils.isBlank(removed)) {
            JSONValidator validator = JSONValidator.from(removed);
            if (!validator.validate()) {
                event.isValidate = false;
                return;
            }
            event.updatedFields.put("removed", JSON.parseObject(removed, Map.class));
        }

        String truncated = StringEscapeUtils.unescapeJava(obj.getString("truncated"));
        if (!StringUtils.isBlank(truncated)) {
            JSONValidator validator = JSONValidator.from(truncated);
            if (!validator.validate()) {
                event.isValidate = false;
                return;
            }
            event.updatedFields.put("truncated", JSON.parseObject(truncated, Map.class));
        }
    }

    public CloudEvent getCloudEvent() {
        CloudEventBuilder builder = CloudEventBuilder.v1();
        builder.withDataContentType("application/json")
                .withId(!Strings.isBlank(this.objectID) ? this.objectID : "unknown");
        String type = this.metadata.getOrDefault("db", "unknown") + "." +
                this.metadata.getOrDefault("collection", "unknown");
        String sourcePrefix = this.metadata.getOrDefault("connector", "unknown") + "."
                + this.metadata.getOrDefault("rs", "unknown");
        builder.withType(type).withSource(URI.create(sourcePrefix + "." + type));
        if (this.metadata.get("ts_ms") != null) {
            builder.withTime(OffsetDateTime.ofInstant(Instant.ofEpochMilli((Long) this.metadata.get("ts_ms")), ZoneOffset.UTC));
        }

        Map<String, Object> data = new HashMap<>();
        if (this.isValidate) {
            if (this.fullFields.size() > 0) {
                data.put("full", this.fullFields);
            }
            if (this.updatedFields.size() > 0) {
                data.put("changed", this.updatedFields);
            }
            data.put("id", this.objectID);
        } else {
            data.put("rawKey", this.rawKey);
            data.put("rawValue", this.rawValue);
        }
        builder.withData(BytesCloudEventData.wrap(JSON.toJSONBytes(data)));
        for (Map.Entry<String, Object> entry : this.metadata.entrySet()) {
            if (!MongoChangeEvent.keyFilter.contains(entry.getKey())) {
                if (entry.getValue() == null) {
                    continue;
                }
                builder.withExtension(EXTENSION_NAME_PREFIX + entry.getKey(), entry.getValue().toString());
            }
        }

        builder.withExtension(EXTENSION_NAME_PREFIX + "operation", this.type);
        builder.withExtension(EXTENSION_NAME_PREFIX + "recognized", this.isValidate);
        return builder.build();
    }

    public String getType() {
        return type;
    }

    public String getObjectID() {
        return objectID;
    }

    public boolean isValidate() {
        return isValidate;
    }

    public HashMap<String, Object> getFullFields() {
        return fullFields;
    }

    public HashMap<String, Object> getUpdatedFields() {
        return updatedFields;
    }

    public HashMap<String, Object> getDeletedFields() {
        return deletedFields;
    }

    public String getRawKey() {
        return rawKey;
    }

    public String getRawValue() {
        return rawValue;
    }
}


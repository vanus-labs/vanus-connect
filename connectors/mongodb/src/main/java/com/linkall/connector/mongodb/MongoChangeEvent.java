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
import com.fasterxml.jackson.databind.JsonNode;
import io.cloudevents.CloudEvent;
import io.cloudevents.CloudEventData;
import io.cloudevents.core.builder.CloudEventBuilder;
import io.cloudevents.core.data.BytesCloudEventData;
import io.cloudevents.jackson.JsonCloudEventData;
import org.apache.commons.text.StringEscapeUtils;

import java.net.URI;
import java.time.Instant;
import java.time.OffsetDateTime;
import java.time.ZoneOffset;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

class MongoChangeEvent {
    private static final String EXTENSION_NAME_PREFIX="vancemongodb";
    private static List<String> keyFilter = new ArrayList<>();
    private String objectID;
    private OpType type;
    private HashMap<String, Object> payload = new HashMap<>();
    private Map<String, Object> metadata;
    private boolean isValidate = true;

    static {
        keyFilter.add("db");
        keyFilter.add("collection");
        keyFilter.add("connector");
        keyFilter.add("rs");
        keyFilter.add("ts_ms");
    }

    public static MongoChangeEvent parse(String data) {
        MongoChangeEvent event = new MongoChangeEvent();
        JSONObject obj = JSON.parseObject(StringEscapeUtils.unescapeJava(data));

        JSONObject payload = obj.getJSONObject("payload");
        event.metadata = payload.getJSONObject("source").getInnerMap();
        switch (payload.getString("op")) {
            case "c":
                event.type = OpType.INSERT;
                processInsertEvent(payload, event);
                break;
            case "u":
                event.type = OpType.UPDATE;
                processUpdateEvent(payload, event);
                break;
            case "d":
                event.type = OpType.DELETE;
                processDeleteEvent(payload, event);
                break;
            default:
                event.isValidate = false;
                return event;
        }
        return event;
    }

    private static Map<String, Object> processInsertEvent(JSONObject data, MongoChangeEvent event) {
        HashMap<String, Object> m = new HashMap<>();
        String body =StringEscapeUtils.unescapeJava(data.getString("after"));

        JSONValidator validator = JSONValidator.from(body);
        if(!validator.validate()) {
            event.isValidate=false;
            return m;
        }
        JSONObject obj = data.getJSONObject("after");
        for(Map.Entry<String, Object> entry : obj.entrySet()) {
            String key = entry.getKey();
            String val = entry.getValue().toString();
            if (key.equals( "_id")) {
                JSONObject id = JSON.parseObject(val);
                if( id.containsKey("$numberLong") ){
                    event.objectID = id.getString("$numberLong");
                    event.payload.put("_id", Long.parseLong(id.getString("$numberLong")));
                }
                continue;
            } else  {
                event.payload.put(key, val);
            }
        }
        return m;
    }

    private static void processUpdateEvent(JSONObject data, MongoChangeEvent event) {

    }

    private static void processDeleteEvent(JSONObject data, MongoChangeEvent event) {

    }

    public CloudEvent getCloudEvent() {
        CloudEventBuilder builder = CloudEventBuilder.v1();
        String type = this.metadata.get("db")+"."+this.metadata.get("collection");
        String sourcePrefix = this.metadata.get("connector")+"."+this.metadata.get("rs");
        builder.withDataContentType("application/json")
                .withId(this.objectID)
                .withType(type)
                .withSource(URI.create(sourcePrefix+"."+type))
                .withTime(OffsetDateTime.ofInstant(Instant.ofEpochMilli((Long)this.metadata.get("ts_ms")), ZoneOffset.UTC))
                .withData(BytesCloudEventData.wrap(JSON.toJSONBytes(this.payload)));

        for (Map.Entry<String, Object> entry : this.metadata.entrySet()) {
            if(!MongoChangeEvent.keyFilter.contains(entry.getKey())) {
                builder.withExtension(EXTENSION_NAME_PREFIX+entry.getKey(), entry.getValue().toString());
            }
        }
        return builder.build();
    }

    public OpType getType() {
        return type;
    }

    public String getObjectID() {
        return objectID;
    }

    public boolean isValidate() {
        return isValidate;
    }
}


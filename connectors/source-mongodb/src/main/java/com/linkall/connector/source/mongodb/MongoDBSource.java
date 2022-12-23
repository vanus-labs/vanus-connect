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

package com.linkall.connector.source.mongodb;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.linkall.cdk.config.Config;
import io.cloudevents.CloudEventData;
import io.cloudevents.jackson.JsonCloudEventData;
import org.apache.commons.text.StringEscapeUtils;

import java.io.IOException;
import java.nio.charset.StandardCharsets;
import java.util.HashMap;
import java.util.Map;

public class MongoDBSource extends com.linkall.cdk.database.debezium.DebeziumSource {
    private static final String CONNECTOR_NAME = "Source MongoDB";
    private final ObjectMapper mapper = new ObjectMapper();

    @Override
    public Class<? extends Config> configClass() {
        return MongoDBConfig.class;
    }

    @Override
    public String name() {
        return CONNECTOR_NAME;
    }

    @Override
    protected CloudEventData convertData(Object data) throws IOException {
        Map<String, Object> m = (Map) data;
        Map<String, Object> result = new HashMap<>();
        for (Map.Entry<String, Object> entry : m.entrySet()) {
            if (entry.getValue() == null) {
                continue;
            }
            switch (entry.getKey()) {
                case "before":
                    break;
                case "after":
                    String json = StringEscapeUtils.unescapeJson(entry.getValue().toString());
                    Map<String, Object> value = this.mapper.readValue(json.getBytes(StandardCharsets.UTF_8), Map.class);
                    if (value.get("_id") != null) {
                        value.put("_id", ((Map) value.get("_id")).get("$oid"));
                    }
                    result.put(entry.getKey(), value);
                    break;
                case "patch":
                    break;
                case "filter":
                    break;
                case "updateDescription":
                    result.put(entry.getKey(), processUpdate(entry.getValue()));
                    break;
            }

        }
        return JsonCloudEventData.wrap(mapper.valueToTree(result));
    }

    // TODO more tests
    private Object processUpdate(Object obj) throws IOException {
        Map<String, Object> m = (Map) obj;
        Map<String, Object> result = new HashMap<>();
        for (Map.Entry<String, Object> entry : m.entrySet()) {
            if (entry.getValue() == null) {
                continue;
            }
            String json = StringEscapeUtils.unescapeJson(entry.getValue().toString());
            Map<String, Object> value = this.mapper.readValue(json.getBytes(StandardCharsets.UTF_8), Map.class);
            result.put(entry.getKey(), value);
        }
        return result;
    }
}



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
import io.cloudevents.CloudEvent;
import org.junit.jupiter.api.Test;

import java.util.Map;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertTrue;

public class MongoDBChangeEventTest {

    @Test
    public void TestCreateEvent(){
        String createEvent = "{\"schema\":{\"type\":\"struct\",\"fields\":[{\"type\":\"string\",\"optional\":true," +
                "\"name\":\"io.debezium.data.Json\",\"version\":1,\"field\":\"after\"},{\"type\":\"string\"," +
                "\"optional\":true,\"name\":\"io.debezium.data.Json\",\"version\":1,\"field\":\"patch\"}," +
                "{\"type\":\"struct\",\"fields\":[{\"type\":\"string\",\"optional\":false,\"field\":\"version\"}," +
                "{\"type\":\"string\",\"optional\":false,\"field\":\"connector\"},{\"type\":\"string\"," +
                "\"optional\":false,\"field\":\"name\"},{\"type\":\"int64\",\"optional\":false,\"field\":\"ts_ms\"}," +
                "{\"type\":\"boolean\",\"optional\":true,\"default\":false,\"field\":\"snapshot\"}," +
                "{\"type\":\"string\",\"optional\":false,\"field\":\"db\"},{\"type\":\"string\",\"optional\":false," +
                "\"field\":\"rs\"},{\"type\":\"string\",\"optional\":false,\"field\":\"collection\"},{\"type\":\"int32\"," +
                "\"optional\":false,\"field\":\"ord\"},{\"type\":\"int64\",\"optional\":true,\"field\":\"h\"}]," +
                "\"optional\":false,\"name\":\"io.debezium.connector.mongo.Source\",\"field\":\"source\"}," +
                "{\"type\":\"string\",\"optional\":true,\"field\":\"op\"},{\"type\":\"int64\",\"optional\":true," +
                "\"field\":\"ts_ms\"}],\"optional\":false,\"name\":\"dbserver1.inventory.customers.Envelope\"}," +
                "\"payload\":{\"after\":\"{\\\\\\\"_id\\\\\\\" : {\\\\\\\"$numberLong\\\\\\\" : \\\\\\\"1004\\\\\\\"}," +
                "\\\\\\\"first_name\\\\\\\" : \\\\\\\"Anne\\\\\\\",\\\\\\\"last_name\\\\\\\" : \\\\\\\"Kretchmar\\\\\\\"," +
                "\\\\\\\"email\\\\\\\" : \\\\\\\"annek@noanswer.org\\\\\\\"}\",\"patch\":null," +
                "\"source\":{\"version\":\"1.9.5.Final\",\"connector\":\"mongodb\",\"name\":\"fulfillment\"," +
                "\"ts_ms\":1558965508000,\"snapshot\":false,\"db\":\"inventory\",\"rs\":\"rs0\"," +
                "\"collection\":\"customers\",\"ord\":31,\"h\":1546547425148722000},\"op\":\"c\"," +
                "\"ts_ms\":1558965515240}}\n";
        MongoChangeEvent event = MongoChangeEvent.parse(createEvent);
        assertEquals(OpType.INSERT, event.getType());
        assertEquals("1004", event.getObjectID());
        assertTrue(event.isValidate());
        CloudEvent ce = event.getCloudEvent();
        assertEquals(5, ce.getExtensionNames().size());
        assertEquals("1004", ce.getId());
        assertEquals("mongodb.rs0.inventory.customers", ce.getSource().toString());
        assertEquals("inventory.customers", ce.getType());
        assertEquals("application/json", ce.getDataContentType());
        Map<String, Object> data = JSON.parseObject(ce.getData().toBytes(), Map.class);
        assertEquals(1004, data.get("_id"));
        assertEquals("Anne", data.get("first_name"));
        assertEquals("Kretchmar", data.get("last_name"));
        assertEquals("annek@noanswer.org", data.get("email"));


    }
}

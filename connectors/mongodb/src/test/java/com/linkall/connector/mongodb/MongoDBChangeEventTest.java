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
    public void TestCreateEvent() {
        String id = "{\"id\":\"{\\\"$oid\\\": \\\"62ff236a99b4cfeac7ed54c6\\\"}\"}";
        String createEvent = "{\"after\":\"{\\\"_id\\\": {\\\"$oid\\\": \\\"62ff236a99b4cfeac7ed54c6\\\"}," +
                "\\\"a\\\": \\\"a\\\"}\",\"patch\":null,\"filter\":null,\"updateDescription\":null," +
                "\"source\":{\"version\":\"1.9.4.Final\",\"connector\":\"mongodb\",\"name\":\"test\"," +
                "\"ts_ms\":1660887914000,\"snapshot\":\"false\",\"db\":\"test\",\"sequence\":null," +
                "\"rs\":\"replicaset-01\",\"collection\":\"source\",\"ord\":1,\"h\":null,\"tord\":null," +
                "\"stxnid\":null,\"lsid\":null,\"txnNumber\":null},\"op\":\"c\",\"ts_ms\":1661223842688," +
                "\"transaction\":null}";
        MongoChangeEvent event = MongoChangeEvent.parse(id, createEvent);
        assertEquals("insert", event.getType());
        assertEquals("62ff236a99b4cfeac7ed54c6", event.getObjectID());
        assertTrue(event.isValidate());
        CloudEvent ce = event.getCloudEvent();
        assertEquals("62ff236a99b4cfeac7ed54c6", ce.getId());
        assertEquals("mongodb.replicaset-01.test.source", ce.getSource().toString());
        assertEquals("test.source", ce.getType());
        assertEquals("application/json", ce.getDataContentType());
        Map<String, Object> data = JSON.parseObject(ce.getData().toBytes(), Map.class);
        Map<String, Object> full = (Map<String, Object>) data.get("full");
        assertEquals("62ff236a99b4cfeac7ed54c6", full.get("_id"));
        assertEquals("a", full.get("a"));
        assertEquals(6, ce.getExtensionNames().size());
    }

    @Test
    public void TestUpdateEvent() {
        String id = "{\"id\":\"{\\\"$oid\\\": \\\"63044b3fccaea8fcf8a159ef\\\"}\"}";
        String updateEvent = "{\"after\":\"{\\\"_id\\\": {\\\"$oid\\\": \\\"63044b3fccaea8fcf8a159ef\\\"}," +
                "\\\"b\\\": \\\"1213\\\"}\",\"patch\":null,\"filter\":null," +
                "\"updateDescription\":{\"removedFields\":null,\"updatedFields\":\"{\\\"b\\\": \\\"1213\\\"}\"," +
                "\"truncatedArrays\":null},\"source\":{\"version\":\"1.9.4.Final\",\"connector\":\"mongodb\"," +
                "\"name\":\"test\",\"ts_ms\":1661225902000,\"snapshot\":\"false\",\"db\":\"test\",\"sequence\":null," +
                "\"rs\":\"replicaset-01\",\"collection\":\"source\",\"ord\":1,\"h\":null,\"tord\":null," +
                "\"stxnid\":null,\"lsid\":null,\"txnNumber\":null},\"op\":\"u\",\"ts_ms\":1661225902776," +
                "\"transaction\":null}";
        MongoChangeEvent event = MongoChangeEvent.parse(id, updateEvent);
        assertEquals("update", event.getType());
        assertEquals("63044b3fccaea8fcf8a159ef", event.getObjectID());
        assertEquals(2, event.getFullFields().size());
        assertEquals(1, event.getUpdatedFields().size());
        assertEquals(0, event.getDeletedFields().size());
        assertTrue(event.isValidate());
        CloudEvent ce = event.getCloudEvent();
        assertEquals("63044b3fccaea8fcf8a159ef", ce.getId());
        assertEquals("mongodb.replicaset-01.test.source", ce.getSource().toString());
        assertEquals("test.source", ce.getType());
        assertEquals("application/json", ce.getDataContentType());
        Map<String, Object> data = JSON.parseObject(ce.getData().toBytes(), Map.class);
        Map<String, Object> full = (Map<String, Object>) data.get("full");
        assertEquals(2, full.size());
        assertEquals("63044b3fccaea8fcf8a159ef", full.get("_id"));
        assertEquals("1213", full.get("b"));
        Map<String, Object> changed = (Map<String, Object>) data.get("changed");
        assertEquals(1, changed.size());
        Map<String, Object> updated = (Map<String, Object>) changed.get("updated");
        assertEquals(1, updated.size());
        assertEquals("1213", updated.get("b"));
        assertEquals(6, ce.getExtensionNames().size());
    }

    @Test
    void TestDeletedEvent() {
        String id = "{\"id\":\"{\\\"$oid\\\": \\\"63044b3fccaea8fcf8a159ef\\\"}\"}";
        String deleted = "{\"after\":null,\"patch\":null,\"filter\":null,\"updateDescription\":null," +
                "\"source\":{\"version\":\"1.9.4.Final\",\"connector\":\"mongodb\",\"name\":\"test\"," +
                "\"ts_ms\":1661232012000,\"snapshot\":\"false\",\"db\":\"test\",\"sequence\":null," +
                "\"rs\":\"replicaset-01\",\"collection\":\"source\",\"ord\":1,\"h\":null,\"tord\":null," +
                "\"stxnid\":null,\"lsid\":null,\"txnNumber\":null},\"op\":\"d\",\"ts_ms\":1661232012563," +
                "\"transaction\":null}";
        MongoChangeEvent event = MongoChangeEvent.parse(id, deleted);
        assertEquals("delete", event.getType());
        assertEquals("63044b3fccaea8fcf8a159ef", event.getObjectID());
        assertEquals(0, event.getFullFields().size());
        assertEquals(0, event.getUpdatedFields().size());
        assertEquals(0, event.getDeletedFields().size());
        assertTrue(event.isValidate());
        CloudEvent ce = event.getCloudEvent();
        assertEquals("63044b3fccaea8fcf8a159ef", ce.getId());
        assertEquals("mongodb.replicaset-01.test.source", ce.getSource().toString());
        assertEquals("test.source", ce.getType());
        assertEquals("application/json", ce.getDataContentType());
        Map<String, Object> data = JSON.parseObject(ce.getData().toBytes(), Map.class);
        assertEquals(0, data.size());
        assertEquals(6, ce.getExtensionNames().size());
    }

    @Test
    void TestUnrecognizedEvent() {
        String id = "1{\"id\":\"{\\\"$oid\\\": \\\"63044b3fccaea8fcf8a159ef\\\"}\"}";
        String unknown = "1{\"after\":null,\"patch\":null,\"filter\":null,\"updateDescription\":null," +
                "\"source\":{\"version\":\"1.9.4.Final\",\"connector\":\"mongodb\",\"name\":\"test\"," +
                "\"ts_ms\":1661232012000,\"snapshot\":\"false\",\"db\":\"test\",\"sequence\":null," +
                "\"rs\":\"replicaset-01\",\"collection\":\"source\",\"ord\":1,\"h\":null,\"tord\":null," +
                "\"stxnid\":null,\"lsid\":null,\"txnNumber\":null},\"op\":\"d\",\"ts_ms\":1661232012563," +
                "\"transaction\":null}";
        MongoChangeEvent event = MongoChangeEvent.parse(id, unknown);
        assertEquals("unknown", event.getType());
        assertEquals(id, event.getRawKey());
        assertEquals(unknown, event.getRawValue());
        CloudEvent ce = event.getCloudEvent();
        assertEquals("unknown", ce.getId());
        assertEquals("unknown.unknown.unknown.unknown", ce.getSource().toString());
        assertEquals("unknown.unknown", ce.getType());
        assertEquals("application/json", ce.getDataContentType());
        Map<String, Object> data = JSON.parseObject(ce.getData().toBytes(), Map.class);
        assertEquals(2, data.size());
        assertEquals(id, data.get("rawKey"));
        assertEquals(unknown, data.get("rawValue"));
        assertEquals(2, ce.getExtensionNames().size());
        assertEquals(false, ce.getExtension("vancemongodbrecognized"));
        assertEquals("unknown", ce.getExtension("vancemongodboperation"));
    }
}

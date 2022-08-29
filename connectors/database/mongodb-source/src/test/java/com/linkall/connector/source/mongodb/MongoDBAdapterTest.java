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
import io.cloudevents.CloudEvent;
import org.junit.jupiter.api.Test;

import java.io.IOException;

import static org.junit.jupiter.api.Assertions.assertEquals;

public class MongoDBAdapterTest {

    @Test
    public void TestCreateEvent() throws IOException {
        String id = "{\"id\":\"{\\\"$oid\\\": \\\"62ff236a99b4cfeac7ed54c6\\\"}\"}";
        String createEvent = "{\"after\":\"{\\\"_id\\\": {\\\"$oid\\\": \\\"62ff236a99b4cfeac7ed54c6\\\"}," +
                "\\\"a\\\": \\\"a\\\"}\",\"patch\":null,\"filter\":null,\"updateDescription\":null," +
                "\"source\":{\"version\":\"1.9.4.Final\",\"connector\":\"mongodb\",\"name\":\"test\"," +
                "\"ts_ms\":1660887914000,\"snapshot\":\"false\",\"db\":\"test\",\"sequence\":null," +
                "\"rs\":\"replicaset-01\",\"collection\":\"source\",\"ord\":1,\"h\":null,\"tord\":null," +
                "\"stxnid\":null,\"lsid\":null,\"txnNumber\":null},\"op\":\"c\",\"ts_ms\":1661223842688," +
                "\"transaction\":null}";
        CloudEvent ce = MongoDBAdapter.proto2CloudEvent(MongoDBAdapter.parse(id, createEvent));
        assertEquals("62ff236a99b4cfeac7ed54c6", ce.getId());
        assertEquals("mongodb.replicaset-01.test.source", ce.getSource().toString());
        assertEquals("test.source", ce.getType());
        assertEquals("application/json", ce.getDataContentType());
        JSONObject obj = JSON.parseObject(ce.getData().toBytes(), JSONObject.class);
        JSONObject document = obj.getJSONObject("insert").getJSONObject("document");
        assertEquals("62ff236a99b4cfeac7ed54c6", document.getJSONObject("_id").get("$oid"));
        assertEquals("a", document.get("a"));
        assertEquals(6, ce.getExtensionNames().size());
        assertEquals("INSERT", ce.getExtension("vancemongodboperation"));
    }

    @Test
    public void TestUpdateEvent() throws IOException {
        String id = "{\"id\":\"{\\\"$oid\\\": \\\"63044b3fccaea8fcf8a159ef\\\"}\"}";
        String updateEvent = "{\"after\":\"{\\\"_id\\\": {\\\"$oid\\\": \\\"63044b3fccaea8fcf8a159ef\\\"}," +
                "\\\"b\\\": \\\"1213\\\"}\",\"patch\":null,\"filter\":null," +
                "\"updateDescription\":{\"removedFields\":null,\"updatedFields\":\"{\\\"b\\\": \\\"1213\\\"}\"," +
                "\"truncatedArrays\":null},\"source\":{\"version\":\"1.9.4.Final\",\"connector\":\"mongodb\"," +
                "\"name\":\"test\",\"ts_ms\":1661225902000,\"snapshot\":\"false\",\"db\":\"test\",\"sequence\":null," +
                "\"rs\":\"replicaset-01\",\"collection\":\"source\",\"ord\":1,\"h\":null,\"tord\":null," +
                "\"stxnid\":null,\"lsid\":null,\"txnNumber\":null},\"op\":\"u\",\"ts_ms\":1661225902776," +
                "\"transaction\":null}";

        CloudEvent ce = MongoDBAdapter.proto2CloudEvent(MongoDBAdapter.parse(id, updateEvent));
        assertEquals("63044b3fccaea8fcf8a159ef", ce.getId());
        assertEquals("mongodb.replicaset-01.test.source", ce.getSource().toString());
        assertEquals("test.source", ce.getType());
        assertEquals("application/json", ce.getDataContentType());
        JSONObject obj = JSON.parseObject(ce.getData().toBytes(), JSONObject.class);
        JSONObject document = obj.getJSONObject("insert").getJSONObject("document");
        assertEquals(2, document.size());
        assertEquals("63044b3fccaea8fcf8a159ef", document.getJSONObject("_id").get("$oid"));
        assertEquals("1213", document.get("b"));
        JSONObject update = obj.getJSONObject("update").getJSONObject("updateDescription");
        assertEquals(3, update.size());
        JSONObject updatedFields = update.getJSONObject("updatedFields");
        assertEquals(1, updatedFields.size());
        assertEquals("1213", updatedFields.get("b"));
        assertEquals(6, ce.getExtensionNames().size());
        assertEquals("UPDATE", ce.getExtension("vancemongodboperation"));
    }

    @Test
    void TestDeletedEvent() throws IOException {
        String id = "{\"id\":\"{\\\"$oid\\\": \\\"63044b3fccaea8fcf8a159ef\\\"}\"}";
        String deleted = "{\"after\":null,\"patch\":null,\"filter\":null,\"updateDescription\":null," +
                "\"source\":{\"version\":\"1.9.4.Final\",\"connector\":\"mongodb\",\"name\":\"test\"," +
                "\"ts_ms\":1661232012000,\"snapshot\":\"false\",\"db\":\"test\",\"sequence\":null," +
                "\"rs\":\"replicaset-01\",\"collection\":\"source\",\"ord\":1,\"h\":null,\"tord\":null," +
                "\"stxnid\":null,\"lsid\":null,\"txnNumber\":null},\"op\":\"d\",\"ts_ms\":1661232012563," +
                "\"transaction\":null}";
        CloudEvent ce = MongoDBAdapter.proto2CloudEvent(MongoDBAdapter.parse(id, deleted));
        assertEquals("63044b3fccaea8fcf8a159ef", ce.getId());
        assertEquals("mongodb.replicaset-01.test.source", ce.getSource().toString());
        assertEquals("test.source", ce.getType());
        assertEquals("application/json", ce.getDataContentType());
        JSONObject obj = JSON.parseObject(ce.getData().toBytes(), JSONObject.class);
        assertEquals(3, obj.size());
        assertEquals(6, ce.getExtensionNames().size());
        assertEquals("DELETE", ce.getExtension("vancemongodboperation"));
    }

    @Test
    void TestUnrecognizedEvent() throws IOException {
        String id = "1{\"id\":\"{\\\"$oid\\\": \\\"63044b3fccaea8fcf8a159ef\\\"}\"}";
        String unknown = "1{\"after\":null,\"patch\":null,\"filter\":null,\"updateDescription\":null," +
                "\"source\":{\"version\":\"1.9.4.Final\",\"connector\":\"mongodb\",\"name\":\"test\"," +
                "\"ts_ms\":1661232012000,\"snapshot\":\"false\",\"db\":\"test\",\"sequence\":null," +
                "\"rs\":\"replicaset-01\",\"collection\":\"source\",\"ord\":1,\"h\":null,\"tord\":null," +
                "\"stxnid\":null,\"lsid\":null,\"txnNumber\":null},\"op\":\"d\",\"ts_ms\":1661232012563," +
                "\"transaction\":null}";
        CloudEvent ce = MongoDBAdapter.proto2CloudEvent(MongoDBAdapter.parse(id, unknown));
        assertEquals("unknown", ce.getId());
        assertEquals("unknown.unknown.unknown.unknown", ce.getSource().toString());
        assertEquals("unknown.unknown", ce.getType());
        assertEquals("application/json", ce.getDataContentType());
        JSONObject obj = JSON.parseObject(ce.getData().toBytes(), JSONObject.class);
        assertEquals(2, obj.size());
        assertEquals(id, obj.getJSONObject("raw").get("key"));
        assertEquals(unknown, obj.getJSONObject("raw").get("value"));
        assertEquals(2, ce.getExtensionNames().size());
        assertEquals(false, ce.getExtension("vancemongodbrecognized"));
        assertEquals("UNKNOWN", ce.getExtension("vancemongodboperation"));
    }
}

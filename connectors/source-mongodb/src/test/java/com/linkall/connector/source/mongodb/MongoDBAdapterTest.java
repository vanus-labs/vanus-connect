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
import junit.framework.Assert;
import org.junit.jupiter.api.Test;

import java.io.IOException;
import java.nio.charset.StandardCharsets;
import java.util.Map;

public class MongoDBAdapterTest {

    @Test
    public void TestConvertData() throws IOException {
        MongoDBSource source = new MongoDBSource();
        ObjectMapper mapper = new ObjectMapper();
        String data = "{\"id\":\"name:test;ts:1671652565000;ord:1\",\"source\":\"/debezium/mongodb/test\"," +
                "\"specversion\":\"1.0\",\"type\":\"io.debezium.mongodb.datachangeevent\",\"time\":\"2022-12-21T19:56:05.000Z\"," +
                "\"datacontenttype\":\"application/json\",\"iodebeziumop\":\"u\",\"iodebeziumversion\":\"2.0.1.Final\"," +
                "\"iodebeziumconnector\":\"mongodb\",\"iodebeziumname\":\"test\",\"iodebeziumtsms\":\"1671652565000\"," +
                "\"iodebeziumsnapshot\":\"false\",\"iodebeziumdb\":\"test\",\"iodebeziumsequence\":null," +
                "\"iodebeziumrs\":\"replicaset-01\",\"iodebeziumcollection\":\"a\",\"iodebeziumord\":1," +
                "\"iodebeziumlsid\":null,\"iodebeziumtxnNumber\":null,\"iodebeziumtxid\":null," +
                "\"iodebeziumtxtotalorder\":null,\"iodebeziumtxdatacollectionorder\":null,\"data\":{\"before\":null," +
                "\"after\":\"{\\\\\\\"_id\\\\\\\": {\\\\\\\"$oid\\\\\\\": \\\\\\\"63a364ad8835b568e786e262\\\\\\\"}," +
                "\\\\\\\"a\\\\\\\": 1234,\\\\\\\"b\\\\\\\": \\\\\\\"3\\\\\\\"}\",\"patch\":null,\"filter\":null," +
                "\"updateDescription\":{\"removedFields\":null,\"updatedFields\":\"{\\\\\\\"a\\\\\\\": 1234}\"," +
                "\"truncatedArrays\":null}}}";
        Map<String, Object> value = mapper.readValue(data.getBytes(StandardCharsets.UTF_8), Map.class);
        Assert.assertEquals(source.convertData(value.get("data")).toString(),"JsonCloudEventData{node={" +
                "\"updateDescription\":{\"updatedFields\":{\"a\":1234}},\"after\":{\"_id\":\"63a364ad8835b568e786e262\"," +
                "\"a\":1234,\"b\":\"3\"}}}" );
    }

}

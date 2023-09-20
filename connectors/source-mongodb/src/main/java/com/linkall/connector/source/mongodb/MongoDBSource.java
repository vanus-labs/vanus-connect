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
import com.linkall.cdk.database.debezium.DebeziumSource;
import io.cloudevents.CloudEventData;
import io.cloudevents.core.builder.CloudEventBuilder;
import io.cloudevents.jackson.JsonCloudEventData;
import io.debezium.connector.AbstractSourceInfo;
import io.debezium.data.Envelope;
import org.apache.commons.text.StringEscapeUtils;
import org.apache.kafka.connect.data.Schema;
import org.apache.kafka.connect.data.Struct;

import java.io.IOException;
import java.nio.charset.StandardCharsets;
import java.util.*;

public class MongoDBSource extends DebeziumSource {
    private static final String CONNECTOR_NAME = "Source MongoDB";
    protected static Set<String> extensionSourceName = new HashSet<>(Arrays.asList(
            AbstractSourceInfo.DATABASE_NAME_KEY,
            AbstractSourceInfo.COLLECTION_NAME_KEY
    ));
    @Override
    public Class<? extends Config> configClass() {
        return MongoDBConfig.class;
    }

    @Override
    public String name() {
        return CONNECTOR_NAME;
    }


    @Override
    protected byte[] eventData(Struct struct) {
        String fieldName = Envelope.FieldName.AFTER;
        Object dataValue = struct.get(fieldName);
        if (dataValue==null) {
            fieldName = Envelope.FieldName.BEFORE;
            dataValue = struct.get(fieldName);
        }
        if (dataValue instanceof String){
            return dataValue.toString().getBytes(StandardCharsets.UTF_8);
        }
        Schema dataSchema = struct.schema().field(fieldName).schema();
        return jsonDataConverter.fromConnectData("debezium", dataSchema, dataValue);
    }

    @Override
    protected void eventExtension(CloudEventBuilder builder, Struct struct) {
        Struct source = struct.getStruct(Envelope.FieldName.SOURCE);
        for (String name : extensionSourceName) {
            builder.withExtension(extensionName(name), source.getString(name));
        }
    }
}



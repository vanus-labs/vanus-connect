package com.linkall.source.postgresql;


import com.fasterxml.jackson.annotation.JsonInclude;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.linkall.cdk.config.Config;
import com.linkall.cdk.database.debezium.DebeziumSource;
import io.cloudevents.CloudEventData;
import io.cloudevents.jackson.JsonCloudEventData;

import java.io.IOException;
import java.util.Map;

public class PostgreSQLSource extends DebeziumSource {

    private ObjectMapper objectMapper = new ObjectMapper();

    public PostgreSQLSource() {
        objectMapper.setSerializationInclusion(JsonInclude.Include.NON_NULL);
    }

    @Override
    protected CloudEventData convertData(Object data) throws IOException {
        Map<String, Object> m = (Map) data;
        Object result = m.get("after");
        if (result==null) {
            // op:d
            result = m.get("before");
        }
        return JsonCloudEventData.wrap(objectMapper.valueToTree(result));
    }

    @Override
    public Class<? extends Config> configClass() {
        return PostgreSQLConfig.class;
    }

    @Override
    public String name() {
        return null;
    }
}

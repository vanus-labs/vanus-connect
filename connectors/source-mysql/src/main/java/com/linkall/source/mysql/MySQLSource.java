package com.linkall.source.mysql;

import com.linkall.cdk.config.Config;
import com.linkall.cdk.database.debezium.DebeziumSource;
import io.cloudevents.core.builder.CloudEventBuilder;
import io.debezium.connector.AbstractSourceInfo;
import io.debezium.data.Envelope;
import org.apache.kafka.connect.data.Struct;

import java.util.Arrays;
import java.util.HashSet;
import java.util.Set;

public class MySQLSource extends DebeziumSource {

    protected static Set<String> extensionSourceName = new HashSet<>(Arrays.asList(
            AbstractSourceInfo.DATABASE_NAME_KEY,
            AbstractSourceInfo.TABLE_NAME_KEY
    ));

    @Override
    protected void eventExtension(CloudEventBuilder builder, Struct struct) {
        Struct source = struct.getStruct(Envelope.FieldName.SOURCE);
        for (String name : extensionSourceName) {
            builder.withExtension(extensionName(name), source.getString(name));
        }
    }


    @Override
    public Class<? extends Config> configClass() {
        return MySQLConfig.class;
    }

    @Override
    public String name() {
        return "Source MySQL";
    }
}

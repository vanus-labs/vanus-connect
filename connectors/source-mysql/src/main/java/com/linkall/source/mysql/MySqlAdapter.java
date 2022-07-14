package com.linkall.source.mysql;

import com.linkall.vance.core.Adapter1;
import io.cloudevents.CloudEvent;
import io.cloudevents.core.builder.CloudEventBuilder;
import io.debezium.connector.AbstractSourceInfo;
import io.debezium.data.Envelope;
import org.apache.kafka.connect.data.Field;
import org.apache.kafka.connect.data.Schema;
import org.apache.kafka.connect.data.Struct;
import org.apache.kafka.connect.json.JsonConverter;
import org.apache.kafka.connect.json.JsonConverterConfig;
import org.apache.kafka.connect.source.SourceRecord;

import java.net.URI;
import java.time.Instant;
import java.time.OffsetDateTime;
import java.time.ZoneOffset;
import java.util.HashMap;
import java.util.Map;
import java.util.function.Function;

public class MySqlAdapter implements Adapter1<SourceRecord> {
  private static final String EXTENSION_NAME_PREFIX = "vancedebezium";
  private static final String TX_ATTRIBUTE_PREFIX = "tx";

  private final JsonConverter jsonDataConverter = new JsonConverter();

  public MySqlAdapter() {
    Map<String, String> ceJsonConfig = new HashMap<>();
    ceJsonConfig.put(JsonConverterConfig.SCHEMAS_ENABLE_CONFIG, "false");
    jsonDataConverter.configure(ceJsonConfig, false);
  }

  @Override
  public CloudEvent adapt(SourceRecord record) {
    Struct struct = (Struct) record.value();
    Struct source = struct.getStruct(Envelope.FieldName.SOURCE);
    CloudEventBuilder eventBuilder = CloudEventBuilder.v1();
    eventBuilder
        .withId(ceId(source))
        .withSource(source())
        .withType(ceType(source))
        .withTime(ceTime(struct));
    eventBuilder.withData("application/json", data(struct));
    extension(eventBuilder, struct);
    return eventBuilder.build();
  }

  private URI source() {
    return URI.create("vance.debezium.mysql");
  }

  private String ceId(Struct source) {
    return source.get(AbstractSourceInfo.DATABASE_NAME_KEY)
        + "."
        + source.get(AbstractSourceInfo.TABLE_NAME_KEY)
        + ":"
        + source.get("file")
        + ":"
        + source.get("pos");
  }

  private String ceType(Struct source) {
    return "debezium.mysql."
        + source.get(AbstractSourceInfo.DATABASE_NAME_KEY)
        + "."
        + source.get(AbstractSourceInfo.TABLE_NAME_KEY);
  }

  private OffsetDateTime ceTime(Struct struct) {
    return OffsetDateTime.ofInstant(
        Instant.ofEpochMilli(struct.getInt64(Envelope.FieldName.TIMESTAMP)), ZoneOffset.UTC);
  }

  private byte[] data(Struct struct) {
    String fieldName = Envelope.FieldName.AFTER;
    Object dataValue = struct.get(fieldName);
    if (dataValue == null) {
      fieldName = Envelope.FieldName.BEFORE;
      dataValue = struct.get(fieldName);
    }
    Schema dataSchema = struct.schema().field(fieldName).schema();
    byte[] data = jsonDataConverter.fromConnectData("debezium", dataSchema, dataValue);
    return data;
  }

  private void extension(CloudEventBuilder eventBuilder, Struct struct) {
    // source
    Struct source = struct.getStruct(Envelope.FieldName.SOURCE);
    setExtension(eventBuilder, source, MySqlAdapter::adjustExtensionName);
    // op
    eventBuilder.withExtension(
        adjustExtensionName(Envelope.FieldName.OPERATION),
        (String) struct.get(Envelope.FieldName.OPERATION));
    // transaction
    Struct transaction = struct.getStruct(Envelope.FieldName.TRANSACTION);
    if (transaction != null) {
      setExtension(eventBuilder, transaction, MySqlAdapter::txExtensionName);
    }
  }

  private void setExtension(
      CloudEventBuilder eventBuilder, Struct struct, Function<String, String> nameFunc) {
    for (Field field : struct.schema().fields()) {
      Object value = struct.get(field);
      if (value == null) {
        continue;
      }
      switch (field.schema().type()) {
        case STRING:
          eventBuilder.withExtension(nameFunc.apply(field.name()), (String) value);
          break;
        case BOOLEAN:
          eventBuilder.withExtension(nameFunc.apply(field.name()), (Boolean) value);
          break;
        case INT64:
        case INT32:
        case INT16:
        case INT8:
        case FLOAT64:
        case FLOAT32:
          eventBuilder.withExtension(nameFunc.apply(field.name()), (Number) value);
          break;
      }
    }
  }

  private static String txExtensionName(String name) {
    return adjustExtensionName(TX_ATTRIBUTE_PREFIX + name);
  }

  private static String adjustExtensionName(String original) {
    StringBuilder sb = new StringBuilder(EXTENSION_NAME_PREFIX);

    char c;
    for (int i = 0; i != original.length(); ++i) {
      c = original.charAt(i);
      if (isExtensionNameValidChar(c)) {
        sb.append(c);
      }
    }
    return sb.toString();
  }

  private static boolean isExtensionNameValidChar(char c) {
    return (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9');
  }
}

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
import java.util.Arrays;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.function.Function;

public class MySqlAdapter implements Adapter1<SourceRecord> {
  private static final String EXTENSION_NAME_PREFIX = "vancedebezium";
  private final JsonConverter jsonDataConverter = new JsonConverter();
  private final List<String> dataFields =
      Arrays.asList(Envelope.FieldName.AFTER, Envelope.FieldName.BEFORE);
  private final CloudEventBuilder eventBuilder;

  public MySqlAdapter() {
    Map<String, String> ceJsonConfig = new HashMap<>();
    ceJsonConfig.put(JsonConverterConfig.SCHEMAS_ENABLE_CONFIG, "false");
    jsonDataConverter.configure(ceJsonConfig, false);
    eventBuilder = CloudEventBuilder.v1().withSource(URI.create("vance.debezium.mysql"));
  }

  static String adjustExtensionName(String original) {
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

  @Override
  public CloudEvent adapt(SourceRecord record) {
    Struct struct = (Struct) record.value();
    Struct source = struct.getStruct(Envelope.FieldName.SOURCE);
    eventBuilder
        .withId(ceId(source))
        .withType(ceType(source))
        .withExtension(
            adjustExtensionName(Envelope.FieldName.OPERATION),
            (String) struct.get(Envelope.FieldName.OPERATION))
        .withTime(ceTime(struct));
    setExtension(eventBuilder, source, MySqlAdapter::adjustExtensionName);
    Object dataValue = null;
    Schema dataSchema = null;
    for (String fieldName : dataFields) {
      dataValue = struct.get(fieldName);
      dataSchema = struct.schema().field(Envelope.FieldName.AFTER).schema();
      if (dataValue != null) {
        break;
      }
    }
    byte[] data = jsonDataConverter.fromConnectData("debezium", dataSchema, dataValue);
    eventBuilder.withData("application/json", data);
    return eventBuilder.build();
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

  private void setExtension(
      CloudEventBuilder eventBuilder, Struct struct, Function<String, String> modifyName) {
    for (Field field : struct.schema().fields()) {
      Object value = struct.get(field);
      if (value == null) {
        continue;
      }
      switch (field.schema().type()) {
        case STRING:
          eventBuilder.withExtension(modifyName.apply(field.name()), (String) value);
          break;
        case BOOLEAN:
          eventBuilder.withExtension(modifyName.apply(field.name()), (Boolean) value);
          break;
        case INT64:
        case INT32:
        case INT16:
        case INT8:
        case FLOAT64:
        case FLOAT32:
          eventBuilder.withExtension(modifyName.apply(field.name()), (Number) value);
          break;
      }
    }
  }
}

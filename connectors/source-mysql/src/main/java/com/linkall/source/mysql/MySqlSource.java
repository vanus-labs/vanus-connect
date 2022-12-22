package com.linkall.source.mysql;

import com.fasterxml.jackson.annotation.JsonInclude;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.linkall.cdk.config.Config;
import com.linkall.cdk.database.debezium.DebeziumSource;
import io.cloudevents.CloudEventData;
import io.cloudevents.jackson.JsonCloudEventData;

import java.io.IOException;

public class MySqlSource extends DebeziumSource {

  private ObjectMapper objectMapper = new ObjectMapper();

  public MySqlSource() {
    objectMapper.setSerializationInclusion(JsonInclude.Include.NON_NULL);
  }

  @Override
  protected CloudEventData convertData(Object data) throws IOException {
    return JsonCloudEventData.wrap(objectMapper.valueToTree(data));
  }

  @Override
  public Class<? extends Config> configClass() {
    return MySqlConfig.class;
  }

  @Override
  public String name() {
    return "Source MySQL";
  }
}

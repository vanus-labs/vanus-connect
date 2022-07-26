package com.linkall.source.mysql;

import com.linkall.source.debezium.DebeziumSource;
import com.linkall.vance.core.Adapter;
import com.linkall.vance.core.Source;
import io.debezium.connector.mysql.MySqlConnector;
import io.debezium.connector.mysql.converters.TinyIntOneToBooleanConverter;

import java.util.HashMap;
import java.util.Map;
import java.util.Properties;

public class MySqlSource extends DebeziumSource implements Source {

  private MySqlOffset offset;

  public MySqlSource() {
    offset = new MySqlOffset();
  }

  @Override
  public String getConnectorClass() {
    return MySqlConnector.class.getCanonicalName();
  }

  @Override
  public Map<String, Object> getConfigOffset() {
    Map<String, Object> offsets = new HashMap<>();
    if (offset.getPos() != null) offsets.put("pos", offset.getPos());
    if (offset.getFile() != null && !offset.getFile().isEmpty())
      offsets.put("file", offset.getFile());
    return offsets;
  }

  @Override
  public Properties getDebeziumProperties() {
    final Properties props = new Properties();

    // convert
    props.setProperty("converters", "boolean, datetime");
    props.setProperty("boolean.type", TinyIntOneToBooleanConverter.class.getCanonicalName());
    props.setProperty("datetime.type", MySqlDateTimeConverter.class.getCanonicalName());

    return props;
  }

  @Override
  public Adapter getAdapter() {
    return new MySqlAdapter();
  }
}

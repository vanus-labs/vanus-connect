package com.linkall.source.mysql;

import com.linkall.source.debezium.DebeziumSource;
import com.linkall.vance.core.Adapter;
import com.linkall.vance.core.Source;
import io.debezium.connector.mysql.MySqlConnector;
import io.debezium.connector.mysql.converters.TinyIntOneToBooleanConverter;

import java.util.Properties;

public class MySqlSource extends DebeziumSource implements Source {

  @Override
  public String getConnectorClass() {
    return MySqlConnector.class.getCanonicalName();
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

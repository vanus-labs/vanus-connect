package com.linkall.source.mysql;

import io.debezium.spi.converter.CustomConverter;
import io.debezium.spi.converter.RelationalColumn;
import org.apache.kafka.connect.data.SchemaBuilder;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.sql.Timestamp;
import java.time.*;
import java.time.format.DateTimeFormatter;
import java.util.Date;
import java.util.Properties;
import java.util.function.Consumer;

public class MySqlDateTimeConverter implements CustomConverter<SchemaBuilder, RelationalColumn> {
  private static final Logger LOGGER = LoggerFactory.getLogger(MySqlDateTimeConverter.class);
  private DateTimeFormatter dateFormatter = DateTimeFormatter.ISO_DATE;
  private DateTimeFormatter timeFormatter = DateTimeFormatter.ISO_TIME;
  private DateTimeFormatter timestampFormatter = DateTimeFormatter.ISO_DATE_TIME;

  private ZoneId timestampZoneId = ZoneOffset.systemDefault();

  @Override
  public void configure(Properties props) {
    readProps(props, "format.date", p -> dateFormatter = DateTimeFormatter.ofPattern(p));
    readProps(props, "format.time", p -> timeFormatter = DateTimeFormatter.ofPattern(p));
    readProps(props, "format.timestamp", p -> timestampFormatter = DateTimeFormatter.ofPattern(p));
    readProps(props, "format.timestamp.zone", z -> timestampZoneId = ZoneId.of(z));
  }

  private void readProps(Properties properties, String settingKey, Consumer<String> callback) {
    String settingValue = (String) properties.get(settingKey);
    if (settingValue == null || settingValue.length() == 0) {
      return;
    }
    try {
      callback.accept(settingValue.trim());
    } catch (IllegalArgumentException | DateTimeException e) {
      LOGGER.error("The \"{}\" setting is illegal:{}", settingKey, settingValue);
      throw e;
    }
  }

  @Override
  public void converterFor(RelationalColumn column, ConverterRegistration<SchemaBuilder> registration) {
    String sqlType = column.typeName().toUpperCase();
    Converter converter = null;
    if ("DATE".equals(sqlType)) {
      converter = this::convertDate;
    }
    if ("TIME".equals(sqlType)) {
      converter = this::convertTime;
    }
    if ("DATETIME".equals(sqlType)) {
      converter = this::convertTimestamp;
    }
    if ("TIMESTAMP".equals(sqlType)) {
      converter = this::convertTimestamp;
    }
    if (converter != null) {
      registration.register(SchemaBuilder.string(), converter);
    }
  }

  private String convertDate(Object input) {
    if (input == null){
      return null;
    }
    if (input instanceof LocalDate) {
      return dateFormatter.format((LocalDate) input);
    }
    if (input instanceof Number) {
      LocalDate date = LocalDate.ofEpochDay( ((Number) input).longValue());
      return dateFormatter.format(date);
    }
    return null;
  }

  private String convertTime(Object input) {
    if (input == null){
      return null;
    }
    if (input instanceof Duration) {
      Duration duration = (Duration) input;
      long seconds = duration.getSeconds();
      int nano = duration.getNano();
      LocalTime time = LocalTime.ofSecondOfDay(seconds).withNano(nano);
      return timeFormatter.format(time);
    }
    return null;
  }

  private String convertTimestamp(Object input) {
    if (input == null){
      return null;
    }
    if (input instanceof Timestamp){
      Timestamp timestamp = (Timestamp) input;
      ZonedDateTime zonedDateTime = ZonedDateTime.ofInstant(timestamp.toInstant(),ZoneOffset.UTC);
      LocalDateTime localDateTime = zonedDateTime.withZoneSameInstant(timestampZoneId).toLocalDateTime();
      return timestampFormatter.format(localDateTime);
    }
    if (input instanceof LocalDateTime) {
      return timestampFormatter.format((LocalDateTime) input);
    }
    if (input instanceof ZonedDateTime) {
      ZonedDateTime zonedDateTime = (ZonedDateTime) input;
      LocalDateTime localDateTime = zonedDateTime.withZoneSameInstant(timestampZoneId).toLocalDateTime();
      return timestampFormatter.format(localDateTime);
    }
    return null;
  }

  public static void main(String[] args) {
    Timestamp timestamp= new Timestamp(new Date().getTime());
    for (String zone : ZoneId.getAvailableZoneIds()) {
      ZoneId zoneId  = ZoneId.of(zone);

      String f = DateTimeFormatter.ISO_OFFSET_DATE_TIME.format(ZonedDateTime.ofInstant(timestamp.toInstant(),zoneId));
      System.out.println(zone);
      System.out.println(f);
    }
  }
}

package com.linkall.source.mysql;

import io.debezium.spi.converter.CustomConverter;
import io.debezium.spi.converter.RelationalColumn;
import org.apache.kafka.connect.data.SchemaBuilder;

import java.sql.Timestamp;
import java.time.*;
import java.time.format.DateTimeFormatter;
import java.util.Arrays;
import java.util.Properties;

public class MySqlDateTimeConverter implements CustomConverter<SchemaBuilder, RelationalColumn> {
  private final DateTimeFormatter timeFormatter = DateTimeFormatter.ISO_TIME;
  private final DateTimeFormatter dateFormatter = DateTimeFormatter.ISO_DATE;
  private final DateTimeFormatter dateTimeFormatter = DateTimeFormatter.ISO_DATE_TIME;
  private ZoneId timestampZoneId = ZoneId.systemDefault();
  private final String[] DATE_TYPES = {"DATE", "DATETIME", "TIMESTAMP"};

  @Override
  public void configure(Properties props) {
    String zoneString = (String) props.get("format.timestamp.zone");
    if (zoneString == null || zoneString.length() == 0) {
      return;
    }
    timestampZoneId = ZoneId.of(zoneString);
  }

  /**
   * | mysql | debezium | binlog | |-----------|-----------|---------------| | time | Duration |
   * Duration | | date | LocalDate | LocalDate | | datetime | Timestamp | LocalDateTime | |
   * timestamp | Timestamp | ZonedDateTime |
   */
  @Override
  public void converterFor(
      RelationalColumn column, ConverterRegistration<SchemaBuilder> registration) {
    String sqlType = column.typeName().toUpperCase();
    if ("TIME".equals(sqlType)) {
      registerTime(column, registration);
    } else if (Arrays.stream(DATE_TYPES).anyMatch(s -> s.equals(sqlType))) {
      registerDate(column, registration);
    }
  }

  private void registerTime(
      final RelationalColumn field, final ConverterRegistration<SchemaBuilder> registration) {
    registration.register(
        SchemaBuilder.string(),
        x -> {
          if (x == null) {
            return null;
          }
          if (x instanceof Duration) {
            Duration duration = (Duration) x;
            long seconds = duration.getSeconds();
            int nano = duration.getNano();
            LocalTime time = LocalTime.ofSecondOfDay(seconds).withNano(nano);
            return timeFormatter.format(time);
          }
          return null;
        });
  }

  private void registerDate(
      final RelationalColumn field, final ConverterRegistration<SchemaBuilder> registration) {
    registration.register(
        SchemaBuilder.string(),
        x -> {
          if (x == null) {
            return null;
          }
          if (x instanceof LocalDate) {
            return dateFormatter.format((LocalDate) x);
          } else if (x instanceof LocalDateTime) {
            return dateTimeFormatter.format((LocalDateTime) x);
          } else if (x instanceof Timestamp) {
            return dateTimeFormatter.format(((Timestamp) x).toLocalDateTime());
          } else if (x instanceof ZonedDateTime) {
            return dateTimeFormatter.format(
                ((ZonedDateTime) x).withZoneSameInstant(timestampZoneId).toLocalDateTime());
          } else if (x instanceof Number) {
            return dateTimeFormatter.format(
                new Timestamp(((Number) x).longValue()).toLocalDateTime());
          }
          return x;
        });
  }
}

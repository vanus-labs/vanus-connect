package com.linkall.sink.mysql;

import com.linkall.vance.common.env.EnvUtil;
import com.linkall.vance.core.Sink;
import com.linkall.vance.core.http.HttpServer;
import io.vertx.core.json.JsonObject;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.sql.SQLException;
import java.util.ArrayList;
import java.util.List;

public class MySqlSink implements Sink {
  private static final Logger LOGGER = LoggerFactory.getLogger(MySqlSink.class);
  private SqlWriter sqlWriter;
  private final MySqlConfig config;

  public MySqlSink() {
    config =
        new MySqlConfig(
            EnvUtil.getEnvOrConfig("host"),
            EnvUtil.getEnvOrConfig("port"),
            EnvUtil.getEnvOrConfig("username"),
            EnvUtil.getEnvOrConfig("password"),
            EnvUtil.getEnvOrConfig("database"),
            EnvUtil.getEnvOrConfig("table_name"),
            EnvUtil.getEnvOrConfig("insert_mode"),
            EnvUtil.getEnvOrConfig("commit_interval"));
  }

  @Override
  public void start() throws Exception {
    HttpServer server = HttpServer.createHttpServer();
    server.ceHandler(
        event -> {
          LOGGER.info("receive a new event: {}", event);
          if (!"application/json".equals(event.getDataContentType())) {
            LOGGER.info(
                "only process contentType application/json, now contentType: {}",
                event.getDataContentType());
            return;
          }
          JsonObject data = new JsonObject(new String(event.getData().toBytes()));
          if (sqlWriter == null) {
            try {
              this.sqlWriter = getSqlWriter(data);
            } catch (SQLException e) {
              LOGGER.error("get sql writer fail", e);
            }
          }
          try {
            sqlWriter.add(data);
          } catch (SQLException e) {
            LOGGER.error("write data has error", e);
          }
        });
    server.listen();
  }

  public SqlWriter getSqlWriter(JsonObject data) throws SQLException {
    List<String> columnNames = new ArrayList<>(data.fieldNames());
    TableMetadata meta = new TableMetadata(config.getTableName(), columnNames);
    SqlWriter sqlWriter = new SqlWriter(this.config, meta);
    sqlWriter.init();
    return sqlWriter;
  }
}

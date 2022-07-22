package com.linkall.sink.mysql;

import com.linkall.sink.mysql.connection.ConnectionProvider;
import com.linkall.sink.mysql.connection.HikariConnectionProvider;
import com.linkall.sink.mysql.dialect.Dialect;
import com.linkall.sink.mysql.dialect.MySqlDialect;
import com.linkall.sink.mysql.executor.InsertExecutor;
import com.linkall.sink.mysql.executor.SqlExecutor;
import com.linkall.sink.mysql.executor.UpsertExecutor;
import io.vertx.core.json.JsonObject;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.sql.SQLException;
import java.util.concurrent.Executors;
import java.util.concurrent.ScheduledExecutorService;
import java.util.concurrent.TimeUnit;

public class SqlWriter {
  private static final Logger LOGGER = LoggerFactory.getLogger(SqlWriter.class);
  private final SqlExecutor<JsonObject> sqlExecutor;
  private final TableMetadata metadata;
  private int batchSize;
  private final int commitSize;
  private final ScheduledExecutorService executorService;
  public static final int MAX_BATCH_SIZE = 2000;
  private final ConnectionProvider connectionProvider;
  private final MySqlConfig sqlConfig;

  public SqlWriter(MySqlConfig sqlConfig, TableMetadata metadata) {
    this.sqlConfig = sqlConfig;
    this.metadata = metadata;
    this.connectionProvider = new HikariConnectionProvider(sqlConfig);
    this.executorService = Executors.newSingleThreadScheduledExecutor();
    this.sqlExecutor = getSqlExecutor();
    this.commitSize = MAX_BATCH_SIZE;
  }

  private SqlExecutor getSqlExecutor() {
    Dialect dialect = new MySqlDialect();
    switch (sqlConfig.getInsertMode()) {
      case INSERT:
        return new InsertExecutor<>(metadata, dialect);
      case UPSERT:
        return new UpsertExecutor<>(metadata, dialect);
      default:
        throw new RuntimeException("Invalid insert mode");
    }
  }

  public void init() throws SQLException {
    sqlExecutor.prepareStatement(connectionProvider.getConnection());
    executorService.scheduleAtFixedRate(
        () -> {
          try {
            flush();
          } catch (SQLException e) {
            LOGGER.warn("sql writer flush has error", e);
          }
        },
        10*1000,
        sqlConfig.getCommitInterval(),
        TimeUnit.MILLISECONDS);
  }

  public synchronized void add(JsonObject data) throws SQLException {
    sqlExecutor.addToBatch(data);
    batchSize++;
    if (batchSize >= commitSize) {
      flush();
    }
  }

  public synchronized void flush() throws SQLException {
    sqlExecutor.executeBatch();
    batchSize = 0;
  }

  public void close() throws Exception{
    executorService.shutdown();
    flush();
    connectionProvider.close();
  }
}

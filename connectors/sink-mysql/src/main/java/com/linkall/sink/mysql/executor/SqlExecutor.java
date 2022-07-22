package com.linkall.sink.mysql.executor;

import java.sql.Connection;
import java.sql.SQLException;

public interface SqlExecutor<T> {

  void prepareStatement(Connection connection) throws SQLException;

  void addToBatch(T data);

  /** Submits a batch of commands to the database for execution. */
  void executeBatch() throws SQLException;

  /** Close JDBC related statements. */
  void close() throws SQLException;
}

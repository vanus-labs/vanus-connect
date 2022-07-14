package com.linkall.sink.mysql.binder;

import java.sql.SQLException;

/** A function to bind the cloud event data into a prepared statement. */
@FunctionalInterface
public interface StatementBinder<T> {
  void bindData(T data) throws SQLException;
}

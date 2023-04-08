package com.linkall.sink.mysql.connection;

import java.sql.Connection;
import java.sql.SQLException;

public interface ConnectionProvider {
    Connection getConnection() throws SQLException;

    boolean isConnectionValid() throws SQLException;

    void close();
}

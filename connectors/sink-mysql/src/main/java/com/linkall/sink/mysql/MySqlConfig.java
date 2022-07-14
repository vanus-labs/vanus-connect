package com.linkall.sink.mysql;

public class MySqlConfig {
  public enum InsertMode {
    INSERT,
    UPSERT,
  }

  public static InsertMode getInsertMode(String insertMode) {
    for (InsertMode mode : InsertMode.values()) {
      if (mode.toString().equalsIgnoreCase(insertMode)) {
        return mode;
      }
    }
    return InsertMode.INSERT;
  }

  private String host;
  private String port;
  private String username;
  private String password;
  private String database;
  private String tableName;
  private InsertMode insertMode;

  public MySqlConfig(
      String host,
      String port,
      String username,
      String password,
      String database,
      String tableName,
      String insertMode) {
    this.host = host;
    this.port = port;
    this.username = username;
    this.password = password;
    this.database = database;
    this.tableName = tableName;
    this.insertMode = getInsertMode(insertMode);
  }

  public String getHost() {
    return host;
  }

  public void setHost(String host) {
    this.host = host;
  }

  public String getPort() {
    return port;
  }

  public void setPort(String port) {
    this.port = port;
  }

  public String getUsername() {
    return username;
  }

  public void setUsername(String username) {
    this.username = username;
  }

  public String getPassword() {
    return password;
  }

  public void setPassword(String password) {
    this.password = password;
  }

  public String getDatabase() {
    return database;
  }

  public void setDatabase(String database) {
    this.database = database;
  }

  public String getTableName() {
    return tableName;
  }

  public void setTableName(String tableName) {
    this.tableName = tableName;
  }

  public InsertMode getInsertMode() {
    return insertMode;
  }

  public void setInsertMode(InsertMode insertMode) {
    this.insertMode = insertMode;
  }
}

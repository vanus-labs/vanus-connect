package com.linkall.source.mysql;

public class MySqlConfig {
  private String host;
  private String port;
  private String username;
  private String password;
  private String database;
  private String[] includeTables;
  private String[] excludeTables;

  public MySqlConfig(
      String host,
      String port,
      String username,
      String password,
      String database,
      String includeTable,
      String excludeTable) {
    this.host = host;
    this.port = port;
    this.username = username;
    this.password = password;
    this.database = database;
    if (includeTable != null) {
      includeTables = includeTable.split(",");
    }
    if (excludeTable != null) {
      excludeTables = excludeTable.split(",");
    }
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

  public String[] getIncludeTables() {
    return includeTables;
  }

  public void setIncludeTables(String[] includeTables) {
    this.includeTables = includeTables;
  }

  public String[] getExcludeTables() {
    return excludeTables;
  }

  public void setExcludeTables(String[] excludeTables) {
    this.excludeTables = excludeTables;
  }
}

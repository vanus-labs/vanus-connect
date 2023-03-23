/*
 * Copyright (c) 2023 Airbyte, Inc., all rights reserved.
 */

package com.linkall.sink.snowflake;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.sql.Connection;
import java.sql.DriverManager;
import java.sql.SQLException;
import java.time.Duration;
import java.util.Properties;

/**
 * SnowflakeDatabase contains helpers to create connections to and run queries on Snowflake.
 */
public class SnowflakeDatabase {

    private static final Logger LOGGER = LoggerFactory.getLogger(SnowflakeDatabase.class);

    private static final Duration NETWORK_TIMEOUT = Duration.ofMinutes(1);
    private static final Duration QUERY_TIMEOUT = Duration.ofHours(3);
    private static final String DRIVER_CLASS_NAME = "net.snowflake.client.jdbc.SnowflakeDriver";

    private Connection connection;
    private Properties properties;
    private String url;
    private int retry = 3;

    static {
        try {
            Class.forName(DRIVER_CLASS_NAME);
        } catch (ClassNotFoundException e) {
            e.printStackTrace();
            System.exit(1);
        }
    }

    public SnowflakeDatabase(DbConfig dbConfig) {
        properties = createProperties(dbConfig);
        url = String.format("jdbc:snowflake://%s", dbConfig.getHost());
    }

    public boolean execute(String sql) throws SQLException {
        LOGGER.debug("execute sql:{}", sql);
        for (int i = 0; i < retry; i++) {
            try {
                return getConnection().createStatement().execute(sql);
            } catch (SQLException e) {
                LOGGER.error("sql execute error, retry times = {}", i, e);
                if (i >= retry) {
                    throw e;
                }
                if (isConnectionValid()) {
                    LOGGER.info("connection valid false, will close connection");
                    closeConnection();
                }
            }
        }
        return false;
    }

    public void closeConnection() {
        if (connection!=null) {
            try {
                connection.close();
            } catch (SQLException e) {
                LOGGER.warn("connection close failed", e);
            } finally {
                connection = null;
            }
        }
    }

    public boolean isConnectionValid() {
        try {
            return connection!=null && connection.isValid(5);
        } catch (SQLException e) {
            LOGGER.error("connection valid error", e);
            return false;
        }
    }

    public synchronized Connection getConnection() throws SQLException {
        if (connection!=null) {
            return connection;
        }
        connection = createConnection();
        if (connection==null) {
            throw new SQLException("new connection is null");
        }
        return connection;
    }

    public Connection createConnection() throws SQLException {
        return DriverManager.getConnection(url, properties);
    }

    private Properties createProperties(DbConfig dbConfig) {
        Properties prop = new Properties();
        prop.put("user", dbConfig.getUsername());
        prop.put("password", dbConfig.getPassword());
        prop.put("db", dbConfig.getDatabase());
        prop.put("schema", dbConfig.getSchema());
        prop.put("warehouse", dbConfig.getWarehouse());
        prop.put("role", dbConfig.getRole());
        prop.put("networkTimeout", Math.toIntExact(NETWORK_TIMEOUT.toSeconds()));
        prop.put("queryTimeout", Math.toIntExact(QUERY_TIMEOUT.toSeconds()));
        prop.put("application", "vanus");
        if (dbConfig.getProperties()!=null) {
            prop.putAll(dbConfig.getProperties());
        }
        return prop;
    }

}

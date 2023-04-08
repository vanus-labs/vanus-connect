package com.linkall.sink.mysql.connection;

import com.linkall.sink.mysql.DbConfig;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.sql.Connection;
import java.sql.DriverManager;
import java.sql.SQLException;
import java.util.Properties;

public class SimpleConnectionProvider implements ConnectionProvider {
    private static final String DRIVER_CLASS_NAME = "com.mysql.cj.jdbc.Driver";
    private static final Logger LOGGER = LoggerFactory.getLogger(SimpleConnectionProvider.class);
    private Connection connection;
    private Properties properties;
    private String url;

    static {
        try {
            Class.forName(DRIVER_CLASS_NAME);
        } catch (ClassNotFoundException e) {
            LOGGER.error("driver class {} not found", DRIVER_CLASS_NAME, e);
            System.exit(1);
        }
    }

    public SimpleConnectionProvider(DbConfig dbConfig) {
        properties = createProperties(dbConfig);
        url = String.format("jdbc:mysql://%s:%s/%s", dbConfig.getHost(), dbConfig.getPort(), dbConfig.getDatabase());
    }

    @Override
    public synchronized boolean isConnectionValid() throws SQLException {
        return connection!=null && connection.isValid(5);
    }

    @Override
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

    @Override
    public void close() {
        if (connection==null) {
            return;
        }
        try {
            connection.close();
        } catch (SQLException e) {
            LOGGER.warn("connection close failed", e);
        } finally {
            connection = null;
        }
    }

    public Connection createConnection() throws SQLException {
        return DriverManager.getConnection(url, properties);
    }

    private Properties createProperties(DbConfig dbConfig) {
        Properties prop = new Properties();
        prop.put("user", dbConfig.getUsername());
        prop.put("password", dbConfig.getPassword());
        prop.putAll(initializeDefaultJdbcProperties());
        if (dbConfig.getProperties()!=null) {
            prop.putAll(dbConfig.getProperties());
        }
        return prop;
    }

    private static Properties initializeDefaultJdbcProperties() {
        Properties defaultJdbcProperties = new Properties();
        defaultJdbcProperties.setProperty("characterEncoding", "UTF-8");
        defaultJdbcProperties.setProperty("characterSetResults", "UTF-8");
        return defaultJdbcProperties;
    }
}

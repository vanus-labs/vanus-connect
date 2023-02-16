/*
 * Copyright (c) 2023 Airbyte, Inc., all rights reserved.
 */

package com.linkall.sink.snowflake;

import com.zaxxer.hikari.HikariConfig;
import com.zaxxer.hikari.HikariDataSource;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import javax.sql.DataSource;
import java.sql.Connection;
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

    private DataSource dataSource;

    public SnowflakeDatabase(DbConfig dbConfig) {
        dataSource = createDataSource(dbConfig);
    }

    public boolean execute(String sql) throws SQLException {
        LOGGER.info("execute sql:{}", sql);
        return dataSource.getConnection().createStatement().execute(sql);
    }

    public Connection getConnection() throws Exception {
        return dataSource.getConnection();
    }

    private DataSource createDataSource(DbConfig dbConfig) {
        final HikariConfig config = new HikariConfig();
        config.setPoolName("connection-pool-" + dbConfig.getHost());
        config.setUsername(dbConfig.getUsername());
        config.setPassword(dbConfig.getPassword());
        config.setDriverClassName(DRIVER_CLASS_NAME);
        config.setJdbcUrl(String.format("jdbc:snowflake://%s", dbConfig.getHost()));
        config.setDataSourceProperties(createProperties(dbConfig));
        return new HikariDataSource(config);
    }


    private Properties createProperties(DbConfig dbConfig) {
        Properties prop = new Properties();
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
